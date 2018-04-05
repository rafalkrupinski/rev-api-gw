package handlers

import (
	"log"
	"net/http"
)

type directHandler struct {
	roundTripper http.RoundTripper
}

func (p *directHandler) HandleRequest(req *http.Request) (*http.Request, *http.Response, error) {
	log.Printf("url:%v", req.URL)
	resp, err := p.roundTripper.RoundTrip(req)
	return req, resp, err
}
