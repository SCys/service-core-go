package core

import (
	"context"
	"net/http"
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
)

// GetResponse execute Get method request with timeout & headers
func GetResponse(_ *context.Context, url string, timeout time.Duration, headers http.Header) (*http.Response, error) {
	client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

	return client.Get(url, headers)
}
