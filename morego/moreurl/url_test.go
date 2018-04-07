package moreurl

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"net/url"
	"testing"
)

const myUrl = "http://example.com/path?arg=val"

func TestUrl_UnmarshalYAML(t *testing.T) {
	//given
	expected, err := url.Parse(myUrl)
	assert.NoError(t, err)
	u := &URL{}

	// when
	err = yaml.UnmarshalStrict([]byte(myUrl), u)

	//then
	assert.NoError(t, err)
	assert.Equal(t, &URL{expected}, u)
}

func TestURL_MarshalYAML(t *testing.T) {
	// given
	p, e := url.Parse(myUrl)
	assert.NoError(t, e)
	u := &URL{p}

	// when
	str, e := u.MarshalYAML()

	// then
	assert.NoError(t, e)
	assert.Equal(t, myUrl, str.(string))
}

func TestURL_MustParseUrl_PanicsOnInvalidUrl(t *testing.T) {
	assert.Panics(t, func() {
		MustParseNetURL("*:")
	})
}

func TestURL_String(t *testing.T) {
	parsed, err := url.Parse(myUrl)
	assert.NoError(t, err)
	u := &URL{URL: parsed}

	assert.Equal(t, myUrl, u.String())
}
