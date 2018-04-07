package handlers

import (
	"bytes"
	"compress/gzip"
	"github.com/rafalkrupinski/rev-api-gw/morego/morehttp"
	"github.com/rafalkrupinski/rev-api-gw/morego/moreio"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func configureVerbosity(c *HandlerChain, verbose bool) {
	if !verbose {
		return
	}
	c.AddRequestHandlerFunc(DumpRequest)
	c.AddResponseHandlerFunc(DumpResponse)
}

func DumpRequest(r *http.Request) (*http.Request, *http.Response, error) {
	log.Printf("[%v] %v %v %v\n", r.URL.Scheme, r.Method, r.URL.RequestURI(), r.Proto)
	log.Printf("Host: %v\n", r.Host)

	log.Printf("url:%v", r.URL.String())

	body, err := dump(r.Body, r.Header)
	r.Body = body
	return nil, nil, err
}

func DumpResponse(_ *http.Request, r *http.Response) *http.Response {
	log.Printf("%v\n", r.Status)

	body, err := dump(r.Body, r.Header)
	if err != nil {
		log.Println(err)
	}
	r.Body = body
	return nil
}

func dump(body io.ReadCloser, h http.Header) (io.ReadCloser, error) {
	if body == nil {
		return nil, nil
	}

	buf, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	err = body.Close()
	if err != nil {
		return nil, err
	}

	origBody, err := doDump(buf, h)

	if err != nil {
		return nil, err
	} else {
		return origBody, nil
	}
}

func doDump(body []byte, h http.Header) (origBody io.ReadCloser, _ error) {
	origBody = moreio.BytesReadCloser(body)

	if h.Get(morehttp.CONTENT_ENC) == "gzip" {
		reader, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		body, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}
	}

	for k, v := range h {
		for _, h := range v {
			log.Printf("%v: %v\n", k, h)
		}
	}
	log.Println("Body follows:")
	log.Println(string(body))
	log.Println(":Body ended")
	return
}
