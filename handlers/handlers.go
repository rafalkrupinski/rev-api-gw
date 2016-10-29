package handlers

import (
	"github.com/elazarl/goproxy"
	"net/http"
	"os"
)

func LogIn(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	ctx.Logf("%+v", req)
	return req, nil
}

func LogOut(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	ctx.Logf("%+v", resp)
	return resp
}

const https = "https"

func UpgradeToHttps(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	scheme := req.URL.Scheme
	if https != scheme {
		ctx.Logf("scheme: %v -> %v", scheme, https)
		req.URL.Scheme = https
	} else {
		ctx.Logf("scheme: %v", scheme)
	}
	return req, nil
}

const kVia = "Via"

var host, _ = os.Hostname()

var via = "EtsyGW 1 " + host

func ViaIn(req *http.Request, _ *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	req.Header.Add(kVia, via)
	return req, nil
}

func ViaOut(resp *http.Response, _ *goproxy.ProxyCtx) *http.Response {
	resp.Header.Add(kVia, via)
	return resp
}
