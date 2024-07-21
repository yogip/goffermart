package accrual

import (
	"fmt"
	"time"
)

type ErrorTooManyRequests struct {
	RetryAfter time.Duration
}

func NewErrorTooManyRequests(retryAfer int64) *ErrorTooManyRequests {
	return &ErrorTooManyRequests{
		RetryAfter: time.Duration(retryAfer) * time.Second,
	}
}

func (e *ErrorTooManyRequests) Error() string {
	return fmt.Sprintf("too many requests, retry after: %d", e.RetryAfter)
}
