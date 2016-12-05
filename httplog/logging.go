package httplog

import (
	"bytes"
	"compress/gzip"
	"fmt"
	ht "github.com/rafalkrupinski/rev-api-gw/http"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type LoggingRoundTripper struct {
	Super http.RoundTripper
}

func (rt *LoggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	DumpRequest(r)

	res, err := rt.Super.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	if res != nil {
		DumpResponse(res)
	}

	return res, err
}

func DumpRequest(r *http.Request) error {
	fmt.Printf("[%v] %v %v %v\n", r.URL.Scheme, r.Method, r.URL.RequestURI(), r.Proto)
	fmt.Printf("Host: %v\n", r.Host)

	body, err := dump(r.Body, r.Header)
	r.Body = body
	return err
}

func DumpResponse(r *http.Response) error {
	fmt.Printf("%v\n", r.Status)

	body, err := dump(r.Body, r.Header)
	r.Body = body
	return err
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
	origBody = ReadCloser{bytes.NewReader(body)}

	if h.Get(ht.CONTENT_ENC) == "gzip" {
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
			fmt.Printf("%v: %v\n", k, h)
		}
	}
	fmt.Println()
	os.Stdout.Write(body)
	fmt.Println()

	return
}

type ReadCloser struct {
	io.Reader
}

func (ReadCloser) Close() error {
	return nil
}
