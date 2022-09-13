package service

import (
	"context"
	"time"
)

// StoreCmdsService describes the service.
type StoreCmdsService interface {
	Store(ctx context.Context, timestamp_exec time.Time, cmd string, success bool, exit_code int, stdout string, stderr string) (err error)
	GetFromTo(ctx context.Context, from time.Time, to time.Time) (res []*CmdExecutedEntry, err error)
}

type basicStoreCmdsService struct {
	r Repository
}

func (b *basicStoreCmdsService) Store(ctx context.Context, timestamp time.Time, cmd string, success bool, exit_code int, stdout string, stderr string) (err error) {
	b.r.CreateCmdExec(ctx, &CmdExecutedEntry{
		Cmd:           cmd,
		TimestampExec: timestamp,
		Success:       success,
		ExitCode:      exit_code,
		Stdout:        stdout,
		Stderr:        stderr,
	})

	return err
}

func (b *basicStoreCmdsService) GetFromTo(ctx context.Context, from time.Time, to time.Time) (res []*CmdExecutedEntry, err error) {
	res, err = b.r.GetCmdExecFromTo(ctx, from, to)
	return
}

// NewBasicStoreCmdsService returns a naive, stateless implementation of StoreCmdsService.
func NewBasicStoreCmdsService(repo Repository) StoreCmdsService {
	return &basicStoreCmdsService{
		r: repo,
	}
}

// New returns a StoreCmdsService with all of the expected middleware wired in.
func New(repository Repository, middleware []Middleware) StoreCmdsService {
	var svc StoreCmdsService = NewBasicStoreCmdsService(repository)
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
