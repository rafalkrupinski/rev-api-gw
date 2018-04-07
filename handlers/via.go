package handlers

import (
	"github.com/rafalkrupinski/rev-api-gw/morego/morehttp"
	"net/http"
	"os"
)

var via string

func init() {
	var host, _ = os.Hostname()
	via = "HTTP/1.1 RevApiGW 1 " + host
}

func ViaIn(req *http.Request) (*http.Request, *http.Response, error) {
	req.Header.Add(morehttp.VIA, via)
	return req, nil, nil
}

func ViaOut(_ *http.Request, resp *http.Response) *http.Response {
	resp.Header.Add(morehttp.VIA, via)
	return nil
}
