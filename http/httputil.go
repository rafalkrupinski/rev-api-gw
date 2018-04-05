package http

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	CONTENT_ENC  = "Content-Encoding"
	VIA          = "Via"
	SCHEME_HTTP  = "http"
	PORT_HTTP    = 80
	SCHEME_HTTPS = "https"
	PORT_HTTPS   = 443
)

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
