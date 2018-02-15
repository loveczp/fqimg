package imageserverlib

import (
	"log"
	"image"
	"flag"
	"fmt"
	"os"
	"github.com/disintegration/imaging"
	"github.com/BurntSushi/toml"
)

var Conf Config
var configPath string
var accessLog *log.Logger
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
	CorsAllow              bool
	ImageUrlPrefix         string
	FileCacheDir           string
	FileCacheSize          int
}

var markHash = make(map[string]image.Image)

func init() {
	flag.StringVar(&configPath, "c", "./config.conf", "config file path")
	flag.Parse()
	fmt.Printf(sformat, "configPath:", configPath)
	if _, err := toml.DecodeFile(configPath, &Conf); err != nil {
		log.Panic("config file decode error.\n", err)
	}
	fmt.Printf(sformat, "storage_type:", Conf.StorageType)
	fmt.Printf(sformat, "file_dir:", Conf.FileDir)
	fmt.Printf(sformat, "weed_master_url:", Conf.WeedMasterUrl)
	fmt.Printf(sformat, "fastdfs_config_file_path:", Conf.FastdfsConfigFilePath)
	fmt.Printf("%-20s%-20d\n", "port:", Conf.Port)
	fmt.Printf(sformat, "favicon_path:", Conf.FaviconPath)
	fmt.Printf(sformat, "default_action:", Conf.DefaultAction)
	fmt.Printf(sformat, "Headers:", Conf.Headers)
	fmt.Printf(sformat, "log_dir:", Conf.LogDir)
	fmt.Printf(sformat, "Markers:", Conf.Markers)
	fmt.Printf(sformat, "UploadAllowed:", Conf.UploadAllowed)
	fmt.Printf(sformat, "UploadDeny:", Conf.UploadDeny)
	fmt.Printf(sformat, "CorsAllow:", Conf.CorsAllow)
	fmt.Printf(sformat, "imageUrlPrefix:", Conf.ImageUrlPrefix)
	fmt.Printf(sformat, "FileCacheDir:", Conf.FileCacheDir)
	fmt.Printf(sformat, "FileCacheSize:", Conf.FileCacheSize)

	if (Conf.FileCacheSize < 3 || Conf.FileCacheSize > 10000) {
		Conf.FileCacheSize = 10000
	}

	if (len(Conf.FileCacheDir) < 2) {
		Conf.FileCacheDir = "/var/go_image_server/"
	}

	if (Conf.UploadDeny != nil && Conf.UploadAllowed != nil) {
		log.Panic("uploadDeny and uploadAllowed and not be set at same time, please  delete uploadDeny or uploadAllowed ")
	}
	Conf.UploadAllowedInterface = parseIp(Conf.UploadAllowed)
	Conf.UploadDenyInterface = parseIp(Conf.UploadDeny)

	fmt.Printf(sformat, "UploadAllowedInterface:", Conf.UploadAllowedInterface)
	fmt.Printf(sformat, "UploadDenyInterface:", Conf.UploadDenyInterface)

	if (len(Conf.Markers) > 0) {
		for k, v := range Conf.Markers {
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

	if len(Conf.LogDir) > 2 {
		os.MkdirAll(Conf.LogDir, 0777);
		logfile := Conf.LogDir + "/access.log"
		if _, err := os.Stat(logfile); os.IsNotExist(err) {
			_, _ = os.Create(logfile);
		}

		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file", "output:", err)
		}

		accessLog = log.New(file, "access: ", log.Ldate|log.Ltime)

		if (Conf.StorageType == "file") {
			store = initFile(Conf);
			fmt.Printf("-------------::::::::::use file as storage ::::::::--------------")
		} else if (Conf.StorageType == "weed") {
			store, _ = initWeed(Conf)
			fmt.Printf("-------------::::::::::use weed as storage ::::::::--------------")
		} else if (Conf.StorageType == "fastdfs") {
			store, _ = initFast(Conf)
			fmt.Printf("-------------::::::::::use fastdfs as storage ::::::::--------------")
		}

	}

	if len(Conf.FileCacheDir) > 2 {
		os.MkdirAll(Conf.FileCacheDir, 0777);
	}
}
