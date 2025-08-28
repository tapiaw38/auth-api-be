package workers

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
	"github.com/tapiaw38/auth-api-be/internal/adapters/queue"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations/notification"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
)

type (
	Worker interface {
		Start(ctx context.Context) error
		Stop() error
	}

	Consumer interface {
		StartConsumer(topic queue.Topic, handler ConsumerHandler) error
	}

	ConsumerHandler func(any) error

	ConsumerManager struct {
		conn      *amqp.Connection
		consumers map[queue.Topic]*consumer
		mutex     sync.Mutex
	}

	consumer struct {
		topic   queue.Topic
		ch      *amqp.Channel
		handler ConsumerHandler
	}
)

func NewConsumerManager(conn *amqp.Connection) *ConsumerManager {
	return &ConsumerManager{
		conn:      conn,
		consumers: make(map[queue.Topic]*consumer),
	}
}

func (cm *ConsumerManager) GetConsumer(topic queue.Topic, handler ConsumerHandler) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if _, ok := cm.consumers[topic]; ok {
		return nil
	}

	ch, err := cm.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	cons := &consumer{
		topic:   topic,
		ch:      ch,
		handler: handler,
	}

	cm.consumers[topic] = cons
	return nil
}

func (cm *ConsumerManager) StartConsumer(topic queue.Topic, handler ConsumerHandler) error {
	log.Printf("Starting consumer for topic: %s", topic)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := cm.GetConsumer(topic, handler)
	if err != nil {
		return fmt.Errorf("failed to create consumer for topic %s: %w", topic, err)
	}

	log.Printf("Consumer for topic %s created successfully", topic)

	if err := cm.Consume(ctx, topic); err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	log.Printf("Consumer for topic %s is running", topic)

	<-ctx.Done()
	log.Println("Shutting down consumer gracefully...")

	return nil
}

func (cm *ConsumerManager) Consume(ctx context.Context, topic queue.Topic) error {
	cm.mutex.Lock()
	cons, ok := cm.consumers[topic]
	cm.mutex.Unlock()

	if !ok {
		return fmt.Errorf("consumer for topic %s not found", topic)
	}

	defer cons.ch.Close()

	q, err := cons.ch.QueueDeclare(
		string(cons.topic), // queue name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return fmt.Errorf("queue declare failed: %w", err)
	}

	msgs, err := cons.ch.Consume(
		string(cons.topic), // queue name
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
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

				if cons.topic == queue.TopicSendEmail {
					var sendEmailBody notification.SendEmailInput
					if err := json.Unmarshal(d.Body, &sendEmailBody); err != nil {
						log.Printf("Error unmarshaling message: %v", err)
						return
					}

					if err := cons.handler(sendEmailBody); err != nil {
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

func RegisterWorkers(ctx context.Context, mq *queue.RabbitMQ, contextFactory appcontext.Factory) error {
	log.Println("Starting to register workers...")

	consumerManager := NewConsumerManager(mq.GetConnection())

	app := contextFactory()
	integrations := app.Integrations

	emailWorker := NewEmailWorker(consumerManager, integrations)
	if err := emailWorker.Start(ctx); err != nil {
		return fmt.Errorf("failed to register email worker: %w", err)
	}

	log.Println("All workers registered successfully")

	return nil
}
