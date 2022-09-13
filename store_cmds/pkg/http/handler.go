package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	endpoint "github.com/gigi214/services_example/store_cmds/pkg/endpoint"
	http1 "github.com/go-kit/kit/transport/http"
)

// makeStoreHandler creates the handler logic
func makeStoreHandler(m *http.ServeMux, endpoints endpoint.Endpoints, options []http1.ServerOption) {
	m.Handle("/store", http1.NewServer(endpoints.StoreEndpoint, decodeStoreRequest, encodeStoreResponse, options...))
}

// decodeStoreRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeStoreRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := endpoint.StoreRequest{ExitCode: -999}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.ExitCode == -999 {
		return nil, err
	}
	return req, err
}

// encodeStoreResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeStoreResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeGetFromToHandler creates the handler logic
func makeGetFromToHandler(m *http.ServeMux, endpoints endpoint.Endpoints, options []http1.ServerOption) {
	m.Handle("/get-from-to", http1.NewServer(endpoints.GetFromToEndpoint, decodeGetFromToRequest, encodeGetFromToResponse, options...))
}

// decodeGetFromToRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetFromToRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := endpoint.GetFromToRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeGetFromToResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeGetFromToResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}
func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}
func ErrorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

// This is used to set the http status, see an example here :
// https://github.com/go-kit/kit/blob/master/examples/addsvc/pkg/addtransport/http.go#L133
func err2code(err error) int {
	return http.StatusInternalServerError
}

type errorWrapper struct {
	Error string `json:"error"`
}
