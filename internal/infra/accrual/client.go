package accrual

import (
	"context"
	"fmt"
	"goffermart/internal/core/config"
	"goffermart/internal/logger"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type StopLimiter struct {
	lock      *sync.RWMutex
	stopUnill time.Time
}

func NewStopLimiter() *StopLimiter {
	return &StopLimiter{
		lock:      &sync.RWMutex{},
		stopUnill: time.Now(),
	}
}

func (s *StopLimiter) CanWork() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return time.Now().After(s.stopUnill)
}

func (s *StopLimiter) StopUntil(interval time.Duration) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.stopUnill = time.Now().Add(interval)
}

type AccrualClient struct {
	baseURL string
	config  *config.Config
	client  *resty.Client
}

func NewAccrualClient(config *config.Config) *AccrualClient {
	client := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second)

	return &AccrualClient{
		baseURL: config.Accrual.Address,
		config:  config,
		client:  client,
	}
}

func (c *AccrualClient) GetOrderAccrual(ctx context.Context, orderID string) (*Accrual, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, orderID)
	output := &Accrual{}

	resp, err := c.client.R().
		SetHeader("content-type", "application/json").
		SetContext(ctx).
		SetResult(output).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("getting order accrual request error: %w", err)
	}
	if resp.StatusCode() == http.StatusNoContent {
		return nil, nil
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		raw := resp.Header().Get("Retry-After")
		ra, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			logger.Log.Error(
				"Could not parse Retry-After for GetOrderAccrual",
				zap.String("OrderId", orderID),
				zap.String("RawRetryAfter", raw),
			)
			ra = 60
		}
		return nil, NewErrorTooManyRequests(ra)
	}

	return output, nil
}
