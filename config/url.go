package config

import "net/url"

type YamlUrl struct {
	*url.URL
}

func (u *YamlUrl) MarshalYAML() (interface{}, error) {
	return u.String(), nil
}

func (u *YamlUrl) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw string
	err := unmarshal(&raw)
	if err != nil {
		return err
	}
	parse, e := u.Parse(raw)
	if e != nil {
		return err
	}

	u.URL = parse
	return nil
}

func NewYamlUrl(raw string) (*YamlUrl, error) {
	parse, e := url.Parse(raw)
	if e != nil {
		return nil, e
	}
	return &YamlUrl{parse}, e
}
