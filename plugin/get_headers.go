package plugin

import (
"net/http"
"github.com/loveczp/fqimg/lib"
)

func Plugin_get_headers(h http.HandlerFunc, conf lib.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for key:= range( conf.Headers){
			w.Header().Add(key,conf.Headers[key])
		}
		h.ServeHTTP(w,r)
	}
}
