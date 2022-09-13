package http

import (
	"bytes"
	"context"
	"encoding/json"
	endpoint1 "github.com/gigi214/services_example/store_cmds/pkg/endpoint"
	http2 "github.com/gigi214/services_example/store_cmds/pkg/http"
	service "github.com/gigi214/services_example/store_cmds/pkg/service"
	endpoint "github.com/go-kit/kit/endpoint"
	http "github.com/go-kit/kit/transport/http"
	"io/ioutil"
	http1 "net/http"
	"net/url"
	"strings"
)

// New returns an AddService backed by an HTTP server living at the remote
// instance. We expect instance to come from a service discovery system, so
// likely of the form "host:port".
func New(instance string, options map[string][]http.ClientOption) (service.StoreCmdsService, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}
	var storeEndpoint endpoint.Endpoint
	{
		storeEndpoint = http.NewClient("POST", copyURL(u, "/store"), encodeHTTPGenericRequest, decodeStoreResponse, options["Store"]...).Endpoint()
	}

	var getFromToEndpoint endpoint.Endpoint
	{
		getFromToEndpoint = http.NewClient("POST", copyURL(u, "/get-from-to"), encodeHTTPGenericRequest, decodeGetFromToResponse, options["GetFromTo"]...).Endpoint()
	}

	return endpoint1.Endpoints{
		GetFromToEndpoint: getFromToEndpoint,
		StoreEndpoint:     storeEndpoint,
	}, nil
}

// EncodeHTTPGenericRequest is a transport/http.EncodeRequestFunc that
// SON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPGenericRequest(_ context.Context, r *http1.Request, request interface{}) error {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// decodeStoreResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded concat response from the HTTP response body. If the response
// as a non-200 status code, we will interpret that as an error and attempt to
//  decode the specific error message from the response body.
func decodeStoreResponse(_ context.Context, r *http1.Response) (interface{}, error) {
	if r.StatusCode != http1.StatusOK {
		return nil, http2.ErrorDecoder(r)
	}
	var resp endpoint1.StoreResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// decodeGetFromToResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded concat response from the HTTP response body. If the response
// as a non-200 status code, we will interpret that as an error and attempt to
//  decode the specific error message from the response body.
func decodeGetFromToResponse(_ context.Context, r *http1.Response) (interface{}, error) {
	if r.StatusCode != http1.StatusOK {
		return nil, http2.ErrorDecoder(r)
	}
	var resp endpoint1.GetFromToResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
func copyURL(base *url.URL, path string) (next *url.URL) {
	n := *base
	n.Path = path
	next = &n
	return
}
