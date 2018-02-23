package plugin

import (
	"net/http"
	"github.com/loveczp/fqimg/lib"
	"io"
)

func Plugin_upload_cors(h http.HandlerFunc, conf lib.Config) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if conf.CorsAllow {

			writer.Header().Add("Access-Control-Allow-Origin", "*")
			writer.Header().Add(
				"Access-Control-Allow-Methods",
				"OPTIONS, HEAD, GET, POST, DELETE",
			)
			writer.Header().Add(
				"Access-Control-Allow-Headers",
				"Content-Type, Content-Range, Content-Disposition",
			)

			if request.Method == http.MethodOptions {
				writer.WriteHeader(http.StatusOK)
				io.WriteString(writer, "");
				return
			}
		}
		h.ServeHTTP(writer, request)
	}
}
