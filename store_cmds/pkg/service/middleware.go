package service

import (
	"context"
	"time"

	log "github.com/go-kit/log"
)

type Middleware func(StoreCmdsService) StoreCmdsService

type loggingMiddleware struct {
	logger log.Logger
	next   StoreCmdsService
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next StoreCmdsService) StoreCmdsService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) Store(ctx context.Context, timestamp_exec time.Time, cmd string, success bool, exit_code int, stdout string, stderr string) (err error) {
	defer func() {
		l.logger.Log("method", "Store", "timestamp_exec", timestamp_exec, "cmd", cmd, "success", success, "exit_code", exit_code, "stdout", stdout, "stderr", stderr, "err", err)
	}()
	return l.next.Store(ctx, timestamp_exec, cmd, success, exit_code, stdout, stderr)
}

func (l loggingMiddleware) GetFromTo(ctx context.Context, from time.Time, to time.Time) (res []*CmdExecutedEntry, err error) {
	defer func() {
		l.logger.Log("method", "GetFromTo", "from", from, "to", to, "res", res, "err", err)
	}()
	return l.next.GetFromTo(ctx, from, to)
}
