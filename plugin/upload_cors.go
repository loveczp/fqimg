package plugin

import (
	"fqimg/lib"
	"io"
	"net/http"
)

func Plugin_upload_cors(h http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if lib.Conf.CorsAllow {
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
				io.WriteString(writer, "")
				return
			}
		}
		h.ServeHTTP(writer, request)
	}
}
