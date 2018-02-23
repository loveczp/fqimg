package lib

import (
	"testing"
	"net/http/httptest"
)

func TestGetCommands(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo?c=abc&d=aef&c=efg&a=a&e=e&11=2&0&9", nil)

	t.Log(req.URL.Query())
	t.Log(req.URL.RawQuery)
}

func TestParseQueryString(t *testing.T)  {
	req := httptest.NewRequest("GET", "http://example.com/foo?c=abc&d=aef&c=efg&a=a&e=e&11=2&0&9", nil)
	t.Log(parseQueryString(req.URL.RawQuery))
}