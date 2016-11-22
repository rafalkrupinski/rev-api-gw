package httplog

import (
	"bytes"
	"compress/gzip"
	"fmt"
	ht "github.com/rafalkrupinski/revapigw/http"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type LoggingRoundTripper struct {
	Super http.RoundTripper
}

func (rt *LoggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	dumpRequest(r)

	res, err := rt.Super.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	if res != nil {
		dumpResponse(res)
	}

	return res, err
}

func dumpRequest(r *http.Request) error {
	fmt.Printf("[%v] %v %v %v\n", r.URL.Scheme, r.Method, r.URL.RequestURI(), r.Proto)
	fmt.Printf("Host: %v\n", r.Host)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = r.Body.Close()
	if err != nil {
		return err
	}

	origBody, err := dump(body, r.Header)

	if err != nil {
		return err
	}

	r.Body = origBody
	return nil
}

func dumpResponse(r *http.Response) error {
	fmt.Printf("%v\n", r.Status)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = r.Body.Close()
	if err != nil {
		return err
	}

	origBody, err := dump(body, r.Header)
	if err != nil {
		return err
	}

	r.Body = origBody
	return nil
}

func dump(body []byte, h http.Header) (origBody io.ReadCloser, _ error) {
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
