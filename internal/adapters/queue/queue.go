package queue

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/streadway/amqp"
	"github.com/tapiaw38/auth-api-be/internal/platform/config"
)

type Topic string

const (
	TopicSendEmail Topic = "send_email"
)

type Publisher interface {
	Publish(topic Topic, data interface{}) error
}

type (
	RabbitMQ struct {
		conn       *amqp.Connection
		publishers map[Topic]*publisher
		mutex      sync.Mutex
	}

	publisher struct {
		topic Topic
		ch    *amqp.Channel
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

func (r *RabbitMQ) GetConnection() *amqp.Connection {
	return r.conn
}

func (r *RabbitMQ) GetPublisher(topic Topic) error {
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
		topic: topic,
		ch:    ch,
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
