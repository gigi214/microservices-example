package endpoint

import (
	service "bash_exec/pkg/service"
	"context"

	endpoint "github.com/go-kit/kit/endpoint"
)

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// ExecCmdRequest collects the request parameters for the ExecCmd method.
type ExecCmdRequest struct {
	Cmd string `json:"cmd"`
}

// ExecCmdResponse collects the response parameters for the ExecCmd method.
type ExecCmdResponse struct {
	StdOut   string `json:"std_out"`
	StdErr   string `json:"std_err"`
	ExitCode int    `json:"exit_code"`
	Err      error  `json:"err"`
}

// MakeExecCmdEndpoint returns an endpoint that invokes ExecCmd on the service.
func MakeExecCmdEndpoint(s service.BashExecService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ExecCmdRequest)
		stdOut, stdErr, exitCode, err := s.ExecCmd(ctx, req.Cmd)
		return ExecCmdResponse{
			Err:      err,
			ExitCode: exitCode,
			StdErr:   stdErr,
			StdOut:   stdOut,
		}, nil
	}
}

// Failed implements Failer.
func (r ExecCmdResponse) Failed() error {
	return r.Err
}

// ExecCmd implements Service. Primarily useful in a client.
func (e Endpoints) ExecCmd(ctx context.Context, cmd string) (stdOut string, stdErr string, exitCode int, err error) {
	request := ExecCmdRequest{Cmd: cmd}
	response, err := e.ExecCmdEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(ExecCmdResponse).StdOut, response.(ExecCmdResponse).StdErr, response.(ExecCmdResponse).ExitCode, response.(ExecCmdResponse).Err
}
