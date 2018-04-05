package util

import "net/url"

func MustParseUrl(rawurl string) *url.URL {
	parsed, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return parsed
}
