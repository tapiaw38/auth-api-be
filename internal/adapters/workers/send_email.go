package workers

import (
	"context"
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

	go func() {
		if err := w.consumerManager.StartConsumer(
			queue.TopicSendEmail,
			func(data any) error {
				input, ok := data.(notification.SendEmailInput)
				if !ok {
					return fmt.Errorf("invalid data type, expected notification.SendEmailInput")
				}
				log.Printf("Processing email for: %s", input.To)
				return w.integrations.Notification.SendEmail(input)
			},
		); err != nil {
			log.Fatalf("Failed to start email consumer: %v", err)
		}
	}()

	log.Println("Email worker started successfully")

	return nil
}

func (w *EmailWorker) Stop() error {
	if w.cancel != nil {
		w.cancel()
	}
	log.Println("Email worker stopped")
	return nil
}
