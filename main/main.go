package main

import (
	"fmt"
	"net/http"
	"gopkg.in/ini.v1"
	"strconv"
	"strings"
	"os"
)

func main() {
	http.HandleFunc("/upload", uploadBinHandler)
	http.HandleFunc("/uploadm", uploadMultiHandler)
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
var defaultAction string
var headers map[string]string
var logDir string
var logfile string


func init() {
	cfg, _ := ini.Load("./config.conf")
	storage_type = cfg.Section("").Key("storage_type").MustString("file")
	port = cfg.Section("").Key("port").MustInt(12345)
	favicoPath = cfg.Section("").Key("faviconIcoPath").MustString("./favicon.ico")
	defaultAction= cfg.Section("").Key("defaultAction").MustString("")
	headerString:= cfg.Section("").Key("headers").MustString("")
	logDir = cfg.Section("").Key("logDir").MustString("")

	//header
	if len(headerString)>2{
		headers=make(map[string]string)
		harray := strings.Split(headerString,";")
		for i := 0; i < len(harray); i++ {
			itemarray:=strings.Split(harray[i],":")
			if len(itemarray)==2 {
				headers[itemarray[0]]=itemarray[1]
			}
		}
	}



	if len(logDir)>2 {
		os.MkdirAll(logDir, 0777);

		if _, err := os.Stat(logDir + "/access.log"); os.IsNotExist(err) {
			_, _ = os.Create(logDir + "/access.log");
			logfile = logDir + "/access.log";
		}

		if (storage_type == "file") {
			store = initFile(cfg);
		}
	}
}



