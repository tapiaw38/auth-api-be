package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/streadway/amqp"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations"
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

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return pub.ch.Publish(
		"", string(topic), false, false,
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

	q, err := pub.ch.QueueDeclare(
		string(pub.topic), false, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("queue declare failed: %w", err)
	}

	msgs, err := pub.ch.Consume(
		q.Name, "", false, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("consume failed: %w", err)
	}

	go func() {
		for {
			select {
			case d, ok := <-msgs:
				if !ok {
					return
				}
				func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Printf("Recovered in message handler: %v\n", r)
						}
					}()
					if err := pub.handler(d.Body); err != nil {
						fmt.Printf("Error handling message: %v\n", err)
					} else {
						_ = d.Ack(false)
					}
				}()
			case <-ctx.Done():
				fmt.Println("Consumer context cancelled")
				return
			}
		}
	}()

	return nil
}

func (r *RabbitMQ) GetPublishers(integrations *integrations.Integrations) error {
	if err := r.GetPublisher(
		TopicSendEmail,
		func(data any) error {
			input, ok := data.(notification.SendEmailInput)
			if !ok {
				return fmt.Errorf("invalid data type, expected notification.SendEmailInput")
			}
			return integrations.Notification.SendEmail(input)
		},
	); err != nil {
		r.Close()
		return err
	}

	return nil
}
