package errors

import (
	"context"
	"log"
)

func (r *applicationError) Log(ctx context.Context) {
	logFields := map[string]interface{}{
		"internal_code":    r.InternalCode(),
		"status_code":      r.StatusCode(),
		"message":          r.Message(),
		"original_message": r.OriginalMessage(),
	}

	// Add extra fields if any
	for k, v := range r.extraFields {
		logFields[k] = v
	}

	// For now using standard log, can be replaced with structured logger like zap or logrus
	log.Printf("[ERROR] %+v", logFields)
}
