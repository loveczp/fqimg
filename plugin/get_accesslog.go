package plugin

import (
"net/http"
)

func Plugin_get_accesslog(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
