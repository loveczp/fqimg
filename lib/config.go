package lib

import (
	"log"
	"image"
	"flag"
	"fmt"
	"os"
	"github.com/disintegration/imaging"
	"github.com/BurntSushi/toml"
	"github.com/loveczp/fqimg/store"
)

var (
	Conf             Config
	configPath       string
	Storage_instance store.Storage
	getAlias         = "get"
)

type Config struct {
	StorageType           string
	FileDir               string
	WeedMasterUrl         string
	FastdfsConfigFilePath string
	Port                  int
	FaviconPath           string
	DefaultAction         string
	Headers               map[string]string
	LogDir                string
	Markers               map[string]string
	UploadAllowed         []string
	UploadDeny            []string
	CorsAllow             bool
	ImageUrlPrefix        string
	UploadKey             string
	FileCacheDir          string
	FileCacheSize         int
}

var markHash = make(map[string]image.Image)

func init() {

	sformat := "%-30s%-20s\n"
	flag.StringVar(&configPath, "c", "./config.conf", "config file path")
	flag.Parse()
	fmt.Printf(sformat, "configPath:", configPath)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Panic("config file not exsit error.\n", err)
	}

	if _, err := toml.DecodeFile(configPath, &Conf); err != nil {
		log.Panic("config file decode error.\n", err)
	}
	fmt.Printf(sformat, "storage_type:", Conf.StorageType)
	fmt.Printf(sformat, "file_dir:", Conf.FileDir)
	fmt.Printf(sformat, "weed_master_url:", Conf.WeedMasterUrl)
	fmt.Printf(sformat, "fastdfs_config_file_path:", Conf.FastdfsConfigFilePath)
	fmt.Printf("%-30s%-20d\n", "port:", Conf.Port)
	fmt.Printf(sformat, "favicon_path:", Conf.FaviconPath)
	fmt.Printf(sformat, "default_action:", Conf.DefaultAction)
	fmt.Printf(sformat, "Headers:", Conf.Headers)
	fmt.Printf(sformat, "log_dir:", Conf.LogDir)
	fmt.Printf(sformat, "Markers:", Conf.Markers)
	fmt.Printf(sformat, "UploadAllowed:", Conf.UploadAllowed)
	fmt.Printf(sformat, "UploadDeny:", Conf.UploadDeny)
	fmt.Printf("%-30s%-20t\n", "CorsAllow:", Conf.CorsAllow)
	fmt.Printf(sformat, "ImageUrlPrefix:", Conf.ImageUrlPrefix)
	fmt.Printf(sformat, "UploadKey:", Conf.UploadKey)
	fmt.Printf(sformat, "FileCacheDir:", Conf.FileCacheDir)
	fmt.Printf("%-30s%-20d\n", "FileCacheSize:", Conf.FileCacheSize)

	if (Conf.FileCacheSize < 3 || Conf.FileCacheSize > 10000) {
		Conf.FileCacheSize = 10000
	}

	if (len(Conf.FileCacheDir) < 2) {
		Conf.FileCacheDir = "/var/go_image_server/"
	}

	if (Conf.UploadDeny != nil && Conf.UploadAllowed != nil) {
		log.Panic("uploadDeny and uploadAllowed and not be set at same time, please  delete uploadDeny or uploadAllowed ")
	}

	if Conf.UploadKey == "" {
		Conf.UploadKey = "file"
	}

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

	if (Conf.StorageType == "file") {
		Storage_instance = store.InitFile(Conf.FileDir);
		fmt.Printf("-------------::::::::::use file as storage ::::::::--------------")
	} else if (Conf.StorageType == "weed") {
		Storage_instance, _ = store.InitWeed(Conf.WeedMasterUrl)
		fmt.Printf("-------------::::::::::use weed as storage ::::::::--------------")
	} else if (Conf.StorageType == "fastdfs") {
		Storage_instance, _ = store.InitFast(Conf.FastdfsConfigFilePath)
		fmt.Printf("-------------::::::::::use fastdfs as storage ::::::::--------------")
	} else {
		log.Panic("storage config not found")
	}
}
