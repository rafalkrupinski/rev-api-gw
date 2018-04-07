package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type AppConfig struct {
	Verbose    bool
	Address    string
	ConfigFile string

	Configuration *EndpointConfig
}

func ParseFlags() (*AppConfig, bool, error) {
	config := &AppConfig{}

	flag.BoolVar(&config.Verbose, "v", false, "should every proxy request be logged to stdout")
	flag.StringVar(&config.Address, "addr", ":8080", "proxy listen address")
	usage := flag.Bool("h", false, "print help")
	flag.StringVar(&config.ConfigFile, "config", "application.yaml", "path to configuration file")
	flag.Parse()

	endpointConfig, err := config.endpointConfig()
	config.Configuration = endpointConfig
	return config, *usage, err
}

func (c *AppConfig) endpointConfig() (*EndpointConfig, error) {
	buf, err := ioutil.ReadFile(c.ConfigFile)
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

func (c *AppConfig) Dump() {
	str, err := yaml.Marshal(c)
	if err != nil {
		log.Panic(err)
	}
	log.Print(string(str))
}
