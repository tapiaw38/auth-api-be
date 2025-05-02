package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/streadway/amqp"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations/notification"
	"github.com/tapiaw38/auth-api-be/internal/platform/config"
)

type Topic string

const (
	TopicSendEmail Topic = "send_email"
)

type (
	Publisher interface {
		Publish(topic Topic, data any) error
	}

	ConsumerHandler func(any) error

	RabbitMQ struct {
		conn       *amqp.Connection
		publishers map[Topic]*publisher
		mutex      sync.Mutex
	}

	publisher struct {
		topic   Topic
		ch      *amqp.Channel
		handler ConsumerHandler
	}
)

func NewRabbitMQ(cfg *config.ConfigurationService) (*RabbitMQ, error) {
	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return &RabbitMQ{
		conn:       conn,
		publishers: make(map[Topic]*publisher),
	}, nil
}

func (r *RabbitMQ) Close() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, pub := range r.publishers {
		if pub.ch != nil {
			_ = pub.ch.Close()
		}
	}
	return r.conn.Close()
}

func (r *RabbitMQ) GetPublisher(topic Topic, handler ConsumerHandler) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.publishers[topic]; ok {
		return nil
	}

	ch, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	pub := &publisher{
		topic:   topic,
		ch:      ch,
		handler: handler,
	}

	r.publishers[topic] = pub
	return nil
}

func (r *RabbitMQ) Publish(topic Topic, data interface{}) error {
	r.mutex.Lock()
	pub, ok := r.publishers[topic]
	r.mutex.Unlock()

	if !ok {
		return fmt.Errorf("publisher for topic %s not found", topic)
	}

	q, err := pub.ch.QueueDeclare(
		string(pub.topic), // Asegúrate de que este nombre coincida con el tema
		false,             // Durable
		false,             // Auto-delete
		false,             // Exclusive
		false,             // No-wait
		nil,               // Arguments
	)
	if err != nil {
		return fmt.Errorf("queue declare failed: %w", err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return pub.ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
}

func (r *RabbitMQ) Consume(ctx context.Context, topic Topic) error {
	r.mutex.Lock()
	pub, ok := r.publishers[topic]
	r.mutex.Unlock()

	if !ok {
		return fmt.Errorf("publisher for topic %s not found", topic)
	}

	defer pub.ch.Close()

	q, err := pub.ch.QueueDeclare(
		string(pub.topic), // queue name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return fmt.Errorf("queue declare failed: %w", err)
	}

	msgs, err := pub.ch.Consume(
		string(pub.topic), // queue name
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		return fmt.Errorf("consume failed: %w", err)
	}

	log.Printf("Waiting for messages on queue: %s", q.Name)

	for {
		select {
		case d, ok := <-msgs:
			if !ok {
				log.Println("Message channel closed")
				return nil
			}
			log.Printf("Received message: %s", string(d.Body))
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered in message handler: %v", r)
					}
				}()

				if pub.topic == TopicSendEmail {
					var sendEmailBody notification.SendEmailInput
					if err := json.Unmarshal(d.Body, &sendEmailBody); err != nil {
						log.Printf("Error unmarshaling message: %v", err)
						return
					}

					if err := pub.handler(sendEmailBody); err != nil {
						log.Printf("Error handling message: %v", err)
					}
				}
			}()
		case <-ctx.Done():
			log.Println("Consumer context cancelled")
			return nil
		}
	}
}

func (r *RabbitMQ) StartConsumer(topic Topic, handler ConsumerHandler) error {
	log.Printf("Starting consumer for topic: %s", topic)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := r.GetPublisher(
		topic,
		func(data any) error {
			return handler(data)
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create publisher for topic %s: %w", topic, err)
	}

	log.Printf("Publisher for topic %s created successfully", topic)

	if err := r.Consume(ctx, topic); err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	log.Printf("Consumer for topic %s is running", topic)

	<-ctx.Done()
	log.Println("Shutting down consumer gracefully...")

	return nil
}
