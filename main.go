package main

import (
	"net/http"
	"strconv"
	"log"
	"github.com/loveczp/fqimg/lib"
	"github.com/gorilla/mux"
	"github.com/loveczp/fqimg/plugin"
	"os"
)

func main() {
	lib.InitConfig()
	log.SetOutput(os.Stdout)
	route := mux.NewRouter()
	upload := plugin.Plugin_upload_cors(lib.UploadHandler(lib.Storage_instance),lib.Conf)
	route.PathPrefix("/put").HandlerFunc(upload)



	get := plugin.Plugin_get_filecache(lib.GetHandler(lib.Storage_instance), lib.Conf)
	get = plugin.Plugin_get_headers(get, lib.Conf)
	route.PathPrefix("/get").HandlerFunc(get)
	route.HandleFunc("/hello", lib.HelloHandler())
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(lib.Conf.Port), route))
}