package moreurl

import "net/url"

// Successfully parse net.URL or panic
func MustParseNetURL(rawurl string) *url.URL {
	parsed, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return parsed
}

// Successfully parse URL or panic
func MustParseURL(rawurl string) *URL {
	return &URL{MustParseNetURL(rawurl)}
}

// YAML marshalable replacement for url.URL
// this could use url.URL as underlying type but the whole code gets more messy as we don't inherit its methods
type URL struct{ *url.URL }

func (u *URL) MarshalYAML() (interface{}, error) {
	return u.URL.String(), nil
}

func (u *URL) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw string
	err := unmarshal(&raw)
	if err != nil {
		return err
	}
	u.URL, err = url.Parse(raw)
	if err != nil {
		return err
	}
	return nil
}
