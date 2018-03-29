package plugin

import (
	"net/http"
)

func Plugin_jwtauth(h http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
	}
}

func Plugin_jwtauth_config_check() {

}
