package workers

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
	"github.com/tapiaw38/auth-api-be/internal/adapters/queue"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
)

type (
	Worker interface {
		Start(ctx context.Context) error
		Stop() error
	}

	ConsumerHandler func([]byte) error

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

func (cm *ConsumerManager) Consume(ctx context.Context, topic queue.Topic) error {
	cm.mutex.Lock()
	cons, ok := cm.consumers[topic]
	cm.mutex.Unlock()

	if !ok {
		return fmt.Errorf("consumer for topic %s not found", topic)
	}

	defer cons.ch.Close()

	q, err := cons.ch.QueueDeclare(
		string(cons.topic),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue declare failed: %w", err)
	}

	if err := cons.ch.Qos(1, 0, false); err != nil {
		return fmt.Errorf("qos setup failed: %w", err)
	}

	msgs, err := cons.ch.Consume(
		string(cons.topic),
		"",
		false,
		false,
		false,
		false,
		nil,
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
				ackMessage := false
				requeueMessage := false

				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered in message handler: %v", r)
						requeueMessage = true
					}

					if ackMessage {
						if err := d.Ack(false); err != nil {
							log.Printf("Failed to ack message: %v", err)
						}
						return
					}

					if err := d.Nack(false, requeueMessage); err != nil {
						log.Printf("Failed to nack message: %v", err)
					}
				}()

				if err := cons.handler(d.Body); err != nil {
					log.Printf("Error handling message: %v", err)
					requeueMessage = shouldRequeue(err)
					return
				}

				ackMessage = true
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
