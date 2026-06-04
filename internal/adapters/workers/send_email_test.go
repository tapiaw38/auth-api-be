package workers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations/notification"
)

type notificationStub struct {
	sendEmail func(notification.SendEmailInput) error
}

func (s notificationStub) SendEmail(input notification.SendEmailInput) error {
	return s.sendEmail(input)
}

func TestEmailWorkerHandler(t *testing.T) {
	type fields struct {
		notification notification.Integration
	}

	expectedErr := errors.New("smtp unavailable")

	tests := map[string]struct {
		body            []byte
		prepare         func(f *fields)
		expectedErr     error
		expectedRequeue bool
	}{
		"invalid JSON is marked as non retryable": {
			body: []byte(`invalid-json`),
			prepare: func(f *fields) {
				f.notification = notificationStub{
					sendEmail: func(notification.SendEmailInput) error { return nil },
				}
			},
			expectedErr:     errors.New("failed to unmarshal email input: invalid character 'i' looking for beginning of value"),
			expectedRequeue: false,
		},
		"integration failure remains retryable": {
			body: []byte(`{"to":"user@example.com","subject":"hello","template_name":"welcome","variables":{"name":"John"}}`),
			prepare: func(f *fields) {
				f.notification = notificationStub{
					sendEmail: func(input notification.SendEmailInput) error {
						assert.Equal(t, "user@example.com", input.To)
						return expectedErr
					},
				}
			},
			expectedErr:     expectedErr,
			expectedRequeue: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f := fields{}
			if tc.prepare != nil {
				tc.prepare(&f)
			}

			worker := &EmailWorker{
				integrations: &integrations.Integrations{
					Notification: f.notification,
				},
			}

			err := worker.handler()(tc.body)

			assert.Error(t, err)
			assert.EqualError(t, err, tc.expectedErr.Error())
			assert.Equal(t, tc.expectedRequeue, shouldRequeue(err))
		})
	}
}

func TestShouldRequeue(t *testing.T) {
	tests := map[string]struct {
		err             error
		expectedRequeue bool
	}{
		"nil error does not requeue": {
			err:             nil,
			expectedRequeue: false,
		},
		"transient error requeues": {
			err:             errors.New("transient error"),
			expectedRequeue: true,
		},
		"permanent error does not requeue": {
			err:             Permanent(errors.New("bad payload")),
			expectedRequeue: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedRequeue, shouldRequeue(tc.err))
		})
	}
}
