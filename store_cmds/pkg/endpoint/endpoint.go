package endpoint

import (
	"context"
	"errors"
	"time"

	service "github.com/gigi214/services_example/store_cmds/pkg/service"
	endpoint "github.com/go-kit/kit/endpoint"
)

// StoreRequest collects the request parameters for the Store method.
type StoreRequest struct {
	TimestampExec time.Time `json:"timestamp_exec"`
	Cmd           string    `json:"cmd"`
	Success       bool      `json:"success"`
	ExitCode      int       `json:"exit_code"`
	Stdout        string    `json:"stdout"`
	Stderr        string    `json:"stderr"`
}

// StoreResponse collects the response parameters for the Store method.
type StoreResponse struct {
	Err error `json:"err"`
}

var ErrInvalidInput = errors.New("invalid inputs")

// MakeStoreEndpoint returns an endpoint that invokes Store on the service.
func MakeStoreEndpoint(s service.StoreCmdsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if request == nil {
			return nil, ErrInvalidInput
		}

		req := request.(StoreRequest)
		err := s.Store(ctx, req.TimestampExec, req.Cmd, req.Success, req.ExitCode, req.Stdout, req.Stderr)
		return StoreResponse{Err: err}, nil
	}
}

// Failed implements Failer.
func (r StoreResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// Store implements Service. Primarily useful in a client.
func (e Endpoints) Store(ctx context.Context, timestamp_exec time.Time, cmd string, success bool, exit_code int, stdout string, stderr string) (err error) {
	request := StoreRequest{
		Cmd:           cmd,
		ExitCode:      exit_code,
		Stderr:        stderr,
		Stdout:        stdout,
		Success:       success,
		TimestampExec: timestamp_exec,
	}
	response, err := e.StoreEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(StoreResponse).Err
}

// GetFromToRequest collects the request parameters for the GetFromTo method.
type GetFromToRequest struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// GetFromToResponse collects the response parameters for the GetFromTo method.
type GetFromToResponse struct {
	Res []*service.CmdExecutedEntry `json:"res"`
	Err error                       `json:"err"`
}

// MakeGetFromToEndpoint returns an endpoint that invokes GetFromTo on the service.
func MakeGetFromToEndpoint(s service.StoreCmdsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetFromToRequest)
		res, err := s.GetFromTo(ctx, req.From, req.To)
		return GetFromToResponse{
			Err: err,
			Res: res,
		}, nil
	}
}

// Failed implements Failer.
func (r GetFromToResponse) Failed() error {
	return r.Err
}

// GetFromTo implements Service. Primarily useful in a client.
func (e Endpoints) GetFromTo(ctx context.Context, from time.Time, to time.Time) (res []*service.CmdExecutedEntry, err error) {
	request := GetFromToRequest{
		From: from,
		To:   to,
	}
	response, err := e.GetFromToEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetFromToResponse).Res, response.(GetFromToResponse).Err
}
