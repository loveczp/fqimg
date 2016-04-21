package main

import (
	"fmt"
	"net/http"
	"gopkg.in/ini.v1"
	"strconv"
	"strings"
	"os"
	"flag"
	"log"
)

func main() {
	http.HandleFunc("/upload", uploadBinHandler)
	http.HandleFunc("/uploadm", uploadMultiHandler)
	http.HandleFunc("/favicon.ico", handleFav)
	http.HandleFunc("/", getHandler)
	http.HandleFunc("/hello", helloHandle)
	err := http.ListenAndServe(":" + strconv.Itoa(port), nil)
	if err != nil {
		log.Panic("server start failed :", err)
		return
	} else {
		log.Printf("server start at port: " + strconv.Itoa(port))
	}
}

var (
	store storage
	storage_type string
	port int
	favicon_path string
	default_action string
	headers map[string]string
	log_dir string
	sformat = "%-20s%-20s\n"
	configPath string
	accessLog *log.Logger
)

func init() {
	flag.StringVar(&configPath, "c", "./config.conf", "config file path")
	flag.Parse()
	fmt.Printf(sformat, "configPath:", configPath)
	cfg, err := ini.Load(configPath)
	if err != nil {
		fmt.Printf("config file does not exsit, default config will be loaded\n")
		cfg = ini.Empty()
	}
	storage_type = cfg.Section("").Key("storage_type").MustString("file")
	fmt.Printf(sformat, "storage_type:", storage_type)

	port = cfg.Section("").Key("port").MustInt(12345)
	fmt.Printf(sformat, "port:", strconv.Itoa(port))

	favicon_path = cfg.Section("").Key("favicon_path").MustString("./favicon.ico")
	fmt.Printf(sformat, "favicon_path:", favicon_path)
	default_action = cfg.Section("").Key("default_action").MustString("")
	fmt.Printf(sformat, "default_action:", default_action)
	headerString := cfg.Section("").Key("headers").MustString("Cache-Control:max-age=9999999")
	fmt.Printf(sformat, "headers:", headerString)
	log_dir = cfg.Section("").Key("log_dir").MustString("/var/go_image_server")
	fmt.Printf(sformat, "log_dir:", log_dir)

	//header
	if len(headerString) > 2 {
		headers = make(map[string]string)
		harray := strings.Split(headerString, ";")
		for i := 0; i < len(harray); i++ {
			itemarray := strings.Split(harray[i], ":")
			if len(itemarray) == 2 {
				headers[itemarray[0]] = itemarray[1]
			}
		}
	}

	if len(log_dir) > 2 {
		os.MkdirAll(log_dir, 0777);
		logfile :=log_dir + "/access.log"
		if _, err := os.Stat(logfile); os.IsNotExist(err) {
			_, _ = os.Create(logfile);
		}

		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file", "output:", err)
		}

		accessLog = log.New(file,"access: ",		log.Ldate|log.Ltime)

		if (storage_type == "file") {
			store = initFile(cfg);
		}
	}
}



