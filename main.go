package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/loveczp/fqimg/lib"
	"github.com/loveczp/fqimg/plugin"
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

	if lib.Conf.HttpPort > 0 && lib.Conf.HttpsPort > 0 {
		go func() {
			log.Fatal(http.ListenAndServe(":"+strconv.Itoa(lib.Conf.HttpPort), route))
		}()
		err_https := http.ListenAndServeTLS(":"+strconv.Itoa(lib.Conf.HttpsPort), lib.Conf.HttpsPublicKeyPath, lib.Conf.HttpsPrivateKeyPath, route)
		if err_https != nil {
			log.Fatal("Web server (HTTPS): ", err_https)
		}
	} else {
		if lib.Conf.HttpPort > 0 {
			log.Fatal(http.ListenAndServe(":"+strconv.Itoa(lib.Conf.HttpPort), route))
		}

		if lib.Conf.HttpsPort > 0 {
			log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(lib.Conf.HttpsPort), lib.Conf.HttpsPublicKeyPath, lib.Conf.HttpsPrivateKeyPath, route))
		}
	}

}
