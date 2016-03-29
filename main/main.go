package main

import (
	"fmt"
	"net/http"
	"gopkg.in/ini.v1"
	"strconv"
)

func main() {
	http.HandleFunc("/upload", uploadBinHandler)
	http.HandleFunc("/uploadMulti", uploadMultiHandler)
	http.HandleFunc("/favicon.ico", handleFav)
	http.HandleFunc("/", getHandler)
	http.HandleFunc("/hello", helloHandle)
	err := http.ListenAndServe(":" + strconv.Itoa(port), nil)
	if err != nil {
		fmt.Println("server start failed :", err)
		return
	}else {
		fmt.Println("server start at port: "+ strconv.Itoa(port))
	}
}

var store storage
var storage_type string
var port int
var favicoPath string

func init() {
	cfg, _ := ini.Load("./config.conf")
	storage_type = cfg.Section("").Key("storage_type").MustString("file")
	port = cfg.Section("").Key("port").MustInt(8080)
	favicoPath = cfg.Section("").Key("faviconIcoPath").MustString("./favicon.ico")
	if (storage_type == "file") {
		store = initFile(cfg);
	}
}



