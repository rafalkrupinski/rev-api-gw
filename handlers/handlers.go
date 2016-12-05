package handlers

import (
	"github.com/elazarl/goproxy"
	ht "github.com/rafalkrupinski/rev-api-gw/http"
	"net/http"
	"os"
)

func UpgradeToHttps(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	scheme := req.URL.Scheme
	if ht.SCHEME_HTTPS != scheme {
		ctx.Logf("scheme: %v -> %v", scheme, ht.SCHEME_HTTPS)
		req.URL.Scheme = ht.SCHEME_HTTPS
	} else {
		ctx.Logf("scheme: %v", scheme)
	}
	return req, nil
}

var via *ht.ViaAdder

func init() {
	var host, _ = os.Hostname()
	via = ht.NewViaAdder("RevApiGW 1 " + host)
}

func ViaIn(req *http.Request, _ *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	via.Alter(&req.Header)
	return req, nil
}

func ViaOut(resp *http.Response, _ *goproxy.ProxyCtx) *http.Response {
	via.Alter(&resp.Header)
	return resp
}

func Pass(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	resp, err := ctx.RoundTrip(req)
	if err != nil {
		ctx.Logf("Error", err)
	}

	return req, resp
}

func RespondWith(contentType string, status int, body string) goproxy.FuncReqHandler {
	return func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return r, goproxy.NewResponse(r, contentType, status, body)
	}
}

func CleanupHandler(req *http.Request, _ *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	ht.CleanupRequest(req)
	return req, nil
}
