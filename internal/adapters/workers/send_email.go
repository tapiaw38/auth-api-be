package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tapiaw38/auth-api-be/internal/adapters/queue"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations/notification"
)

type EmailWorker struct {
	consumerManager *ConsumerManager
	integrations    *integrations.Integrations
	ctx             context.Context
	cancel          context.CancelFunc
}

func NewEmailWorker(consumerManager *ConsumerManager, integrations *integrations.Integrations) *EmailWorker {
	return &EmailWorker{
		consumerManager: consumerManager,
		integrations:    integrations,
	}
}

func (w *EmailWorker) Start(ctx context.Context) error {
	w.ctx, w.cancel = context.WithCancel(ctx)

	if err := w.consumerManager.GetConsumer(queue.TopicSendEmail, w.handler()); err != nil {
		return fmt.Errorf("failed to initialize email consumer: %w", err)
	}

	go func() {
		if err := w.consumerManager.Consume(w.ctx, queue.TopicSendEmail); err != nil {
			log.Printf("Email consumer stopped with error: %v", err)
		}
	}()

	log.Println("Email worker started successfully")

	return nil
}

func (w *EmailWorker) handler() ConsumerHandler {
	return func(body []byte) error {
		var input notification.SendEmailInput
		if err := json.Unmarshal(body, &input); err != nil {
			return Permanent(fmt.Errorf("failed to unmarshal email input: %w", err))
		}
		log.Printf("Processing email for: %s", input.To)
		return w.integrations.Notification.SendEmail(input)
	}
}

func (w *EmailWorker) Stop() error {
	if w.cancel != nil {
		w.cancel()
	}
	log.Println("Email worker stopped")
	return nil
}
