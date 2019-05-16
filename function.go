package wait

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go/wait"
)

var _ wait.Strategy = (*FunctionStrategy)(nil)

type Function func(ctx context.Context, target wait.StrategyTarget) error

type FunctionStrategy struct {
	startupTimeout time.Duration
	f              Function
	retryCount     int
}

func NewFunctionStrategy(f Function, retryCount int) *FunctionStrategy {
	return &FunctionStrategy{
		startupTimeout: 60 * time.Second,
		f:              f,
		retryCount:     retryCount,
	}
}

func (ws *FunctionStrategy) WithStartupTimeout(startupTimeout time.Duration) *FunctionStrategy {
	ws.startupTimeout = startupTimeout
	return ws
}

func (ws *FunctionStrategy) WithRetryCount(retryCount int) *FunctionStrategy {
	ws.retryCount = retryCount
	return ws
}

func (ws *FunctionStrategy) WaitUntilReady(ctx context.Context, target wait.StrategyTarget) error {
	ctx, cancel := context.WithTimeout(ctx, ws.startupTimeout)
	defer cancel()
	retry := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			retry++
			err := ws.f(ctx, target)
			if err == nil {
				return nil
			}
			if retry > ws.retryCount {
				return err
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}
