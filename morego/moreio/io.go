package moreio

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
)

func StrReadCloser(str string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(str))
}

func BytesReadCloser(buf []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(buf))
}
