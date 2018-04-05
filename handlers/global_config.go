package handlers

import "github.com/rafalkrupinski/rev-api-gw/config"

type GlobalConfig struct {
	*config.EndpointConfig
	// informative, for logging
	Config  string
	Verbose bool
}
