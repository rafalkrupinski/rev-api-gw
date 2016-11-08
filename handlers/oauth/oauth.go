package oauth

import (
	"github.com/elazarl/goproxy"
	"net/http"
)

type signer struct {
	client *http.Client
}

func (s *signer) Handle(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	if _, has := req.Header["oauth_signature"]; has {
		return req, nil
	}

	ctx.Logf("OAuth1.0a Sign %v", req.URL)

	resp, err := s.client.Do(req)

	if err != nil {
		resp = goproxy.NewResponse(req, "", 500, "")
		resp.Header = http.Header{}
	}

	return req, resp
}

func New(client *http.Client) *signer {
	return &signer{client}
}
