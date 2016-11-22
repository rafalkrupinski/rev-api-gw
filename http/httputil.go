package http

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	CONTENT_TYPE    = "Content-Type"
	CONTENT_LEN     = "Content-Length"
	CONTENT_ENC     = "Content-Encoding"
	VIA             = "Via"
	DEFAULT_TIMEOUT = time.Second * 30
	SCHEME_HTTP     = "http"
	PORT_HTTP       = 80
	SCHEME_HTTPS    = "https"
	PORT_HTTPS      = 443
)

type HeaderAlterer interface {
	Alter(*http.Header)
}

type ViaAdder struct {
	value string
}

func NewViaAdder(via string) *ViaAdder {
	return &ViaAdder{via}
}

func (v *ViaAdder) Alter(h *http.Header) {
	h.Add(VIA, v.value)
}

type ClientBuilder struct {
	c *http.Client
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{&http.Client{Timeout: DEFAULT_TIMEOUT}}
}

func (b *ClientBuilder) WithTransport(v http.RoundTripper) *ClientBuilder {
	b.c.Transport = v
	return b
}

type CheckRedirect func(req *http.Request, via []*http.Request) error

func (b *ClientBuilder) WithCheckRedirect(v CheckRedirect) *ClientBuilder {
	b.c.CheckRedirect = v
	return b
}

func (b *ClientBuilder) WithJar(v http.CookieJar) *ClientBuilder {
	b.c.Jar = v
	return b
}

func (b *ClientBuilder) Build() *http.Client {
	return b.c
}

func CleanupRequest(req *http.Request) {
	addr := req.URL

	if req.RequestURI != "" {
		addr, _ = url.Parse(req.RequestURI)
		req.RequestURI = ""
	}

	if addr.Host == "" && req.Host != "" {
		addr.Host = req.Host
	}

	if addr.Scheme == SCHEME_HTTP {
		strings.TrimSuffix(addr.Host, ":"+strconv.Itoa(PORT_HTTP))
	} else if addr.Scheme == SCHEME_HTTPS {
		strings.TrimSuffix(addr.Host, ":"+strconv.Itoa(PORT_HTTPS))
	}
}
