package handlers

import (
	ht "github.com/rafalkrupinski/rev-api-gw/http"
	"net/http"
	"os"
)

var via string

func init() {
	var host, _ = os.Hostname()
	via = "HTTP/1.1 RevApiGW 1 " + host
}

func ViaIn(req *http.Request) (*http.Request, *http.Response, error) {
	req.Header.Add(ht.VIA, via)
	return req, nil, nil
}

func ViaOut(_ *http.Request, resp *http.Response) *http.Response {
	resp.Header.Add(ht.VIA, via)
	return nil
}
