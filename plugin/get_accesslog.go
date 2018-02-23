package plugin

import (
"net/http"
"github.com/loveczp/fqimg/lib"
)

func Plugin_get_accesslog(h http.HandlerFunc, conf lib.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
