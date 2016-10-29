package oauth

import (
	"github.com/drone/go-bitbucket/oauth1"
	"github.com/elazarl/goproxy"
	"net/http"
	"strings"
)

type signer struct {
	token    oauth1.Token
	consumer *oauth1.Consumer
}

func (s *signer) Handle(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	req.URL.Host = strings.TrimSuffix(req.URL.Host, ":443")

	ctx.Logf("OAuth1.0a Sign %v", req.URL)

	err := s.consumer.Sign(req, s.token)
	if err != nil {
		panic(err)
	}

	return req, nil
}

func New(token oauth1.Token, consumer *oauth1.Consumer) *signer {
	return &signer{token, consumer}
}

func NewConsumer(consumerKey, consumerSecret string) *oauth1.Consumer {
	return &oauth1.Consumer{ConsumerKey: consumerKey, ConsumerSecret: consumerSecret}
}
