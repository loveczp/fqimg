package plugin

import (
	"net/http"

	"github.com/loveczp/fqimg/lib"
)

func Plugin_get_headers(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(lib.Conf.Headers) > 0 {
			for key := range lib.Conf.Headers {
				w.Header().Add(key, lib.Conf.Headers[key])
			}
		}
		h.ServeHTTP(w, r)
	}
}
