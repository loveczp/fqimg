package main

import (
	"fmt"
	"net/http"
	"strconv"
	"os"
	"flag"
	"log"
	_ "net/http/pprof"
	"github.com/BurntSushi/toml"
	"image"
	"github.com/disintegration/imaging"
	imslib "go_image_server/imageserverlib"
)

func main() {
	http.HandleFunc("/upload", uploadBinHandler)
	http.HandleFunc("/uploadm", uploadMultiHandler)
	http.HandleFunc("/favicon.ico", handleFav)
	http.HandleFunc("/", getHandler)
	http.HandleFunc("/hello", helloHandle)
	http.HandleFunc("/test", uploadTestHandler)
	http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir("/pictest"))))

	err := http.ListenAndServe(":" + strconv.Itoa(conf.Port), nil)
	if err != nil {
		log.Panic("server start failed :", err)
		return
	} else {
		log.Printf("server start at port: " + strconv.Itoa(conf.Port))
	}
}

var conf Config
var configPath string
var accessLog  *log.Logger
var store storage
var sformat = "%-20s%-20s\n"

type Config struct {
	StorageType            string
	FileDir                string
	WeedMasterUrl          string
	FastdfsConfigFilePath  string
	Port                   int
	FaviconPath            string
	DefaultAction          string
	Headers                map[string]string
	LogDir                 string
	Markers                map[string]string
	UploadAllowed          []string //from  conf
	UploadAllowedInterface []interface{}
	UploadDeny             []string
	UploadDenyInterface    []interface{}
	FileCacheDir           string
	FileCacheSize          int
}

var markHash = make(map[string]image.Image)

func init() {
	flag.StringVar(&configPath, "c", "./config.conf", "config file path")
	flag.Parse()
	fmt.Printf(sformat, "configPath:", configPath)
	if _, err := toml.DecodeFile(configPath, &conf); err != nil {
		log.Panic("config file decode error.\n", err)
	}
	fmt.Printf(sformat, "storage_type:", conf.StorageType)
	fmt.Printf(sformat, "file_dir:", conf.FileDir)
	fmt.Printf(sformat, "weed_master_url:", conf.WeedMasterUrl)
	fmt.Printf(sformat, "fastdfs_config_file_path:", conf.FastdfsConfigFilePath)
	fmt.Printf("%-20s%-20d\n", "port:", conf.Port)
	fmt.Printf(sformat, "favicon_path:", conf.FaviconPath)
	fmt.Printf(sformat, "default_action:", conf.DefaultAction)
	fmt.Printf(sformat, "headers:", conf.Headers)
	fmt.Printf(sformat, "log_dir:", conf.LogDir)
	fmt.Printf(sformat, "Markers:", conf.Markers)
	fmt.Printf(sformat, "UploadAllowed:", conf.UploadAllowed)
	fmt.Printf(sformat, "UploadDeny:", conf.UploadDeny)
	fmt.Printf(sformat, "FileCacheDir:", conf.FileCacheDir)
	fmt.Printf(sformat, "FileCacheSize:", conf.FileCacheSize)

	if (conf.FileCacheSize < 3 || conf.FileCacheSize > 10000) {
		conf.FileCacheSize = 10000
	}

	if (len(conf.FileCacheDir) < 2) {
		conf.FileCacheDir = "/var/go_image_server/"
	}

	if (conf.UploadDeny != nil && conf.UploadAllowed != nil) {
		log.Panic("uploadDeny and uploadAllowed and not be set at same time, please  delete uploadDeny or uploadAllowed ")
	}
	conf.UploadAllowedInterface = parseIp(conf.UploadAllowed)
	conf.UploadDenyInterface = parseIp(conf.UploadDeny)

	fmt.Printf(sformat, "UploadAllowedInterface:", conf.UploadAllowedInterface)
	fmt.Printf(sformat, "UploadDenyInterface:", conf.UploadDenyInterface)

	if (len(conf.Markers) > 0) {
		for k, v := range conf.Markers {
			mreader, error := os.Open(v)
			if error != nil {
				log.Panic("open ", v, "error :", error)
			}
			outImage, error := imaging.Decode(mreader)

			if error != nil {
				log.Panic("decode file ", v, "error :", error)
			}
			markHash[k] = outImage
		}
	}

	if len(conf.LogDir) > 2 {
		os.MkdirAll(conf.LogDir, 0777);
		logfile := conf.LogDir + "/access.log"
		if _, err := os.Stat(logfile); os.IsNotExist(err) {
			_, _ = os.Create(logfile);
		}

		file, err := os.OpenFile(logfile, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file", "output:", err)
		}

		accessLog = log.New(file, "access: ", log.Ldate | log.Ltime)

		if (conf.StorageType == "file") {
			store = initFile(conf);
			fmt.Printf("-------------::::::::::use file as storage ::::::::--------------")
		} else if (conf.StorageType == "weed") {
			store, _ = initWeed(conf)
			fmt.Printf("-------------::::::::::use weed as storage ::::::::--------------")
		} else if (conf.StorageType == "fastdfs") {
			store, _ = initFast(conf)
			fmt.Printf("-------------::::::::::use fastdfs as storage ::::::::--------------")
		}

	}

	if len(conf.FileCacheDir) > 2 {
		os.MkdirAll(conf.FileCacheDir, 0777);
	}
}



