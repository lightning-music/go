package lightning

import (
	"github.com/bmizerany/assert"
	"net/url"
	"testing"
)

func TestGetPath(t *testing.T) {
	u, pe := url.Parse("http://foo.org/bar")
	assert.Equal(t, pe, nil)
	p, ge := GetPath(u)
	assert.Equal(t, ge, nil)
	assert.Equal(t, p, "/bar")

	u, pe = url.Parse("http://foo.org/bar?key=val")
	assert.Equal(t, pe, nil)
	p, ge = GetPath(u)
	assert.Equal(t, ge, nil)
	assert.Equal(t, p, "/bar")

	u, pe = url.Parse("http://foo.org/bar?key=val&key2=val2")
	assert.Equal(t, pe, nil)
	p, ge = GetPath(u)
	assert.Equal(t, ge, nil)
	assert.Equal(t, p, "/bar")
}
