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
	//upload
	up := lib.UploadHandler(lib.Storage_instance)
	up = plugin.Plugin_throttle_ip(up)
	up = plugin.Plugin_throttle_total(up)
	up = plugin.Plugin_upload_cors(up)
	up = plugin.Plugin_upload_iplimit(up)
	route.PathPrefix("/put").HandlerFunc(up)

	//get
	get := lib.GetHandler(lib.Storage_instance)
	get = plugin.Plugin_get_headers(get)
	get = plugin.Plugin_get_filecache(get)
	route.PathPrefix("/get").HandlerFunc(get)
	//log.Fatal(http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(lib.Conf.HttpPort), route))
}
