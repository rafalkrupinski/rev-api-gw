package handlers

import (
	"github.com/rafalkrupinski/rev-api-gw/config"
	"github.com/rafalkrupinski/rev-api-gw/morego/moreio"
	"github.com/rafalkrupinski/rev-api-gw/morego/moreurl"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConfigure(t *testing.T) {
	body := "My Mock HTTP Result"

	container := make(MockHandlerContainer)
	roundTripper := &MockRoundTripper{Response: &http.Response{StatusCode: 200, Header: http.Header{}, Body: moreio.StrReadCloser(body)}}
	Configure(container, createEndpointConfig(), roundTripper, false)

	recorder := httptest.NewRecorder()

	request := httptest.NewRequest("GET", "/bin", strings.NewReader(body))
	container["bin"].ServeHTTP(recorder, request)
	log.Println(recorder.Code)
	log.Println(recorder.HeaderMap)
	log.Println("body follows:")
	log.Println(recorder.Body)
	assert.Equal(t, body, recorder.Body.String())
}

func createEndpointConfig() *config.EndpointConfig {
	cfg := &config.EndpointConfig{Endpoints: make(map[string]*config.Endpoint)}
	cfg.Endpoints["bin"] = &config.Endpoint{
		Target: moreurl.MustParseURL("https://httpbin.org/get"),
	}
	return cfg
}

type MockHandlerContainer map[string]http.Handler

func (m MockHandlerContainer) Handle(pattern string, handler http.Handler) {
	m[pattern] = handler
}

type MockRoundTripper struct {
	*http.Request
	*http.Response
	error
}

func (rt *MockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	rt.Request = r

	return rt.Response, rt.error
}
