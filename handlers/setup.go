package handlers

import (
	"github.com/dghubble/oauth1"
	"github.com/rafalkrupinski/rev-api-gw/config"
	"golang.org/x/net/context"
	"log"
	"net/http"
)

// implemented by http.ServeMux
type HandlerRegistry interface {
	Handle(pattern string, handler http.Handler)
}

func Configure(container HandlerRegistry, cfg *config.EndpointConfig, rt http.RoundTripper, verbose bool) {
	for path, endpoint := range cfg.Endpoints {
		configureEndpointChain(endpoint, rt, container, path, verbose)
	}
}

func configureEndpointChain(endpoint *config.Endpoint, rt http.RoundTripper, container HandlerRegistry, path string, verbose bool) {
	log.Printf("%v -> %v; oauth=%v; verbose=%v", path, endpoint.Target, endpoint.Oauth1 != nil, verbose)

	chain := &HandlerChain{}
	// TODO setup correlation id
	chain.AddRequestHandlerFunc(CleanupHandler)
	chain.AddRequestHandlerFunc(ViaIn)
	configureVerbosity(chain, verbose)
	configureEndpoint(path, endpoint, chain, rt)
	chain.AddResponseHandlerFunc(ViaOut)

	container.Handle(path, chain)
}

func configureVerbosity(c *HandlerChain, verbose bool) {
	if !verbose {
		return
	}
	log.Println("Verbose=true")
	c.AddRequestHandlerFunc(DumpRequest)
	c.AddResponseHandlerFunc(DumpResponse)
}

func configureEndpoint(path string, config *config.Endpoint, chain *HandlerChain, rt http.RoundTripper) {
	chain.AddRequestHandler(&requestRewriter{Source: path, Target: config.Target.URL})
	chain.AddRequestHandler(createTransportHandler(config, rt))
}

func createTransportHandler(config *config.Endpoint, rt http.RoundTripper) RequestHandler {
	var roundTripper http.RoundTripper
	if config.Oauth1 != nil {
		roundTripper = createOAuth1Transport(config.Oauth1, rt)
	} else {
		roundTripper = rt
	}
	return &directHandler{roundTripper: roundTripper}
}

func createOAuth1Transport(from *config.Oauth1, rt http.RoundTripper) http.RoundTripper {
	cfg := oauth1.NewConfig(from.ConsumerKey, from.ConsumerSecret)
	token := oauth1.NewToken(from.TokenKey, from.TokenSecret)
	ctx := context.WithValue(context.TODO(), oauth1.HTTPClient, &http.Client{Transport: rt})

	return cfg.Client(ctx, token).Transport
}
