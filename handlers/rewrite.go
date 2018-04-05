package handlers

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

type requestRewriter struct {
	Source string
	Target *url.URL
}

func (r *requestRewriter) HandleRequest(req *http.Request) (*http.Request, *http.Response, error) {
	origUrl := req.URL.String()
	req.Host = ""

	requestUrl, e := r.Target.Parse(strings.TrimPrefix(req.URL.Path, r.Source))
	if e != nil {
		return nil, nil, e
	}
	req.URL.Host = requestUrl.Host
	req.URL.Path = requestUrl.Path
	req.URL.Scheme = requestUrl.Scheme

	log.Printf("%v -> %v", origUrl, req.URL)

	return nil, nil, nil
}
