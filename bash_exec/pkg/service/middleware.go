package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	log "github.com/go-kit/log"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

type Middleware func(BashExecService) BashExecService

type loggingMiddleware struct {
	logger log.Logger
	next   BashExecService
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next BashExecService) BashExecService {
		return &loggingMiddleware{logger, next}
	}
}

func (l loggingMiddleware) ExecCmd(ctx context.Context, cmd string) (stdOut string, stdErr string, exitCode int, err error) {
	defer func() {
		l.logger.Log(
			"method", "ExecCmd",
			"cmd", cmd,
			"stdOut", stdOut,
			"stdErr", stdErr,
			"exitCode", exitCode,
			"err", err,
		)
	}()

	return l.next.ExecCmd(ctx, cmd)
}

type proxyStoreMiddleware struct {
	storeService endpoint.Endpoint
	next         BashExecService
}

// ProxyStoreMiddleware returns a BashExecService Middleware.
// the instances is a string with the StoreService instances address separed by comma, if more than one.
func ProxyStoreMiddleware(ctx context.Context, instances string, logger log.Logger) Middleware {
	if instances == "" {
		logger.Log("call_to", "none")
		return func(next BashExecService) BashExecService { return next }
	}

	// Set some parameters for our client.
	var (
		qps         = 100                    // beyond which we will return an error
		maxAttempts = 3                      // per request, before giving up
		maxTime     = 250 * time.Millisecond // wallclock time, before giving up
	)

	var (
		instanceList = splitInstances(instances)
		endpointer   sd.FixedEndpointer
	)
	logger.Log("call_to", fmt.Sprint(instanceList))
	for _, instance := range instanceList {
		var e endpoint.Endpoint
		e = makeStoreProxy(ctx, instance)
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e)
		endpointer = append(endpointer, e)
	}

	// Now, build a single, retrying, load-balancing endpoint out of all of
	// those individual endpoints.
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)

	// And finally, return the ServiceMiddleware, implemented by proxyStoreMiddleware
	return func(next BashExecService) BashExecService {
		return &proxyStoreMiddleware{storeService: retry, next: next}
	}

}
func (s proxyStoreMiddleware) ExecCmd(ctx context.Context, cmd string) (stdOut string, stdErr string, exitCode int, err error) {
	stdOut, stdErr, exitCode, err = s.next.ExecCmd(ctx, cmd)

	// Call store endpoint for save history of cmds
	_, errDb := s.storeService(ctx, StoreRequest{
		Cmd:      cmd,
		Success:  err == nil,
		ExitCode: exitCode,
		Stdout:   stdOut,
		Stderr:   stdErr,
	})

	fmt.Println("errDB: ", errDb)

	return
}

func makeStoreProxy(ctx context.Context, instance string) endpoint.Endpoint {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		panic(err)
	}
	if u.Path == "" {
		u.Path = "/store"
	}
	return httptransport.NewClient(
		"POST",
		u,
		encodeStoreRequest,
		decodeStoreResponse,
	).Endpoint()
}

func splitInstances(s string) []string {
	a := strings.Split(s, ",")
	for i := range a {
		a[i] = strings.TrimSpace(a[i])
	}
	return a
}

func encodeStoreRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeStoreResponse(context.Context, *http.Response) (response interface{}, err error) {
	return
}

type StoreRequest struct {
	Cmd           string    `json:"cmd"`
	TimestampExec time.Time `json:"timestamp_exec"`
	Success       bool      `json:"success"`
	ExitCode      int       `json:"exit_code"`
	Stdout        string    `json:"stdout,omitempty"`
	Stderr        string    `json:"stderr,omitempty"`
}
