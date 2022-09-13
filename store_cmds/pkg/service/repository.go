package service

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	CreateCmdExec(ctx context.Context, e *CmdExecutedEntry) (err error)
	GetAllCmdExec(ctx context.Context) (res []*CmdExecutedEntry, err error)
	GetCmdExecFromTo(ctx context.Context, from, to time.Time) (res []*CmdExecutedEntry, err error)
}

type CmdExecutedEntry struct {
	Cmd           string    `json:"cmd"`
	TimestampExec time.Time `json:"timestamp_exec"`
	Success       bool      `json:"success"`
	ExitCode      int       `json:"exit_code"`
	Stdout        string    `json:"stdout,omitempty"`
	Stderr        string    `json:"stderr,omitempty"`
}

type repoInMem struct {
	mtx sync.RWMutex
	Db  []*CmdExecutedEntry
}

func NewInMemRepository() (Repository, error) {
	return &repoInMem{
		Db: make([]*CmdExecutedEntry, 0, 4),
	}, nil
}

func (r *repoInMem) CreateCmdExec(ctx context.Context, e *CmdExecutedEntry) (err error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.Db = append(r.Db, e)
	return
}

func (r *repoInMem) GetAllCmdExec(ctx context.Context) (res []*CmdExecutedEntry, err error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	res = r.Db
	return
}

func (r *repoInMem) GetCmdExecFromTo(ctx context.Context, from, to time.Time) (res []*CmdExecutedEntry, err error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	for _, e := range r.Db {
		if from.After(e.TimestampExec) && to.Before(e.TimestampExec) {
			res = append(res, e)
		}
	}

	return
}
