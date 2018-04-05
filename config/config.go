package config

import (
	"errors"
	"github.com/rafalkrupinski/rev-api-gw/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type EndpointConfig struct {
	Endpoints map[string]*Endpoint
}

func (c *EndpointConfig) Sanitize() error {
	if len(c.Endpoints) == 0 {
		return errors.New("no endpoints defined")
	}

	for path, endpoint := range c.Endpoints {
		err := endpoint.Sanitize(path)
		if err != nil {
			return err
		}
	}
	return nil
}

type Endpoint struct {
	Target *util.URL
	Oauth1 *Oauth1
}

func (e *Endpoint) Sanitize(path string) error {
	if e.Target == nil || e.Target.String() == "" {
		return errors.New("missing target")
	}

	if e.Oauth1 != nil {
		err := e.Oauth1.Sanitize()
		if err != nil {
			return errors.New(path + ": " + err.Error())
		}
	}

	return nil
}

type Oauth1 struct {
	ConsumerKey    string `yaml:"consumer_key"`
	ConsumerSecret string `yaml:"consumer_secret"`
	TokenKey       string `yaml:"token_key"`
	TokenSecret    string `yaml:"token_secret"`
}

func (o *Oauth1) Sanitize() error {
	if o.ConsumerKey == "" {
		return errors.New("missing consumer_key")
	}
	if o.ConsumerSecret == "" {
		return errors.New("missing consumer_secret")
	}
	if o.TokenKey == "" {
		return errors.New("missing token_key")
	}
	if o.TokenSecret == "" {
		return errors.New("missing token_secret")
	}
	return nil
}

func ReadEndpointConfig(path string) (*EndpointConfig, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	result := &EndpointConfig{}
	err = yaml.UnmarshalStrict(buf, result)
	if err != nil {
		return nil, err
	}
	return result, result.Sanitize()
}
