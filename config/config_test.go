package config

import (
	"github.com/rafalkrupinski/rev-api-gw/util"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestNilEndpoints(t *testing.T) {
	cfg := &EndpointConfig{}
	err := cfg.Sanitize()
	assert.Error(t, err)
}

func TestEmptyEndpoints(t *testing.T) {
	cfg := &EndpointConfig{Endpoints: make(map[string]*Endpoint)}
	err := cfg.Sanitize()
	assert.Error(t, err)
}

func TestMarshall(t *testing.T) {
	cfg := createEndpointConfig()
	out, _ := yaml.Marshal(cfg)
	assert.Equal(t, exampleValue, string(out))
}

func TestUnmarshal(t *testing.T) {
	config := &EndpointConfig{}
	err := yaml.Unmarshal([]byte(exampleValue), config)
	assert.Nil(t, err)
	assert.Equal(t, createEndpointConfig(), config)
}

func TestUnmarshalEmpty(t *testing.T) {
	config := &EndpointConfig{}
	err := yaml.Unmarshal([]byte(""), config)
	assert.Nil(t, err)
	assert.Error(t, config.Sanitize())
}

var exampleValue = `endpoints:
  etsy:
    target: https://httpbin.org/get
    oauth1:
      consumer_key: ck
      consumer_secret: cs
      token_key: tk
      token_secret: ts
`

func createEndpointConfig() *EndpointConfig {
	cfg := &EndpointConfig{Endpoints: make(map[string]*Endpoint)}
	cfg.Endpoints["etsy"] = &Endpoint{
		Target: util.MustParseURL("https://httpbin.org/get"),
		Oauth1: &Oauth1{
			ConsumerKey:    "ck",
			ConsumerSecret: "cs",
			TokenKey:       "tk",
			TokenSecret:    "ts",
		},
	}
	return cfg
}
