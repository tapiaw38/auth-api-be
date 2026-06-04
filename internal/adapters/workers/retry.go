package workers

import "errors"

type nonRetryableError struct {
	err error
}

func (e nonRetryableError) Error() string {
	return e.err.Error()
}

func (e nonRetryableError) Unwrap() error {
	return e.err
}

func Permanent(err error) error {
	if err == nil {
		return nil
	}

	return nonRetryableError{err: err}
}

func shouldRequeue(err error) bool {
	var nonRetryable nonRetryableError
	return err != nil && !errors.As(err, &nonRetryable)
}
