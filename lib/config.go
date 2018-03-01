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
	HttpPort              int
	FaviconPath           string
	DefaultAction         string
	Headers               map[string]string
	LogDir                string
	Markers               map[string]string
	UploadIpAllowed       []string
	UploadIpDeny          []string
	UploadIpLookups       []string
	CorsAllow             bool
	ImageUrlPrefix        string
	UploadKey             string
	UploadFileSizeLimit   int
	UploadFileNmuLimit    int
	FileCacheDir          string
	FileCacheSize         int
	GetForceAction        string
	UploadThrottlePerIp   int
	UploadThrottleTotal   int
}

var markHash = make(map[string]image.Image)

func InitConfig() {

	sformat := "%-30s%-20s\n"
	dformat := "%-30s%-20d\n"
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
	if Conf.FileDir != "" {
		fmt.Printf(sformat, "file_dir:", Conf.FileDir)
	}
	if Conf.WeedMasterUrl != "" {
		fmt.Printf(sformat, "weed_master_url:", Conf.WeedMasterUrl)
	}
	if Conf.FastdfsConfigFilePath != "" {
		fmt.Printf(sformat, "fastdfs_config_file_path:", Conf.FastdfsConfigFilePath)
	}

	if (Conf.StorageType == "file") {
		Storage_instance = store.InitFile(Conf.FileDir);
		fmt.Printf("-------------::::::::::use file as storage ::::::::--------------\n")
	} else if (Conf.StorageType == "weed") {
		Storage_instance, _ = store.InitWeed(Conf.WeedMasterUrl)
		fmt.Printf("-------------::::::::::use weed as storage ::::::::--------------")
	} else if (Conf.StorageType == "fastdfs") {
		Storage_instance, _ = store.InitFast(Conf.FastdfsConfigFilePath)
		fmt.Printf("-------------::::::::::use fastdfs as storage ::::::::--------------")
	} else {
		log.Panic("storage config not found")
	}

	fmt.Printf("%-30s%-20d\n", "port:", Conf.HttpPort)
	if Conf.HttpPort == 0 {
		Conf.HttpPort = 80
	}
	//fmt.Printf(sformat, "favicon_path:", Conf.FaviconPath)
	if Conf.DefaultAction != "" {
		fmt.Printf(sformat, "default_action:", Conf.DefaultAction)
	}
	if len(Conf.Headers) > 0 {
		fmt.Printf(sformat, "Headers:", Conf.Headers)
	}
	if Conf.LogDir != "" {
		fmt.Printf(sformat, "log_dir:", Conf.LogDir)
	}

	if len(Conf.Markers) > 0 {
		fmt.Printf(sformat, "Markers:", Conf.Markers)
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
	}
	if len(Conf.UploadIpAllowed) > 0 {
		fmt.Printf(sformat, "UploadIpAllowed:", Conf.UploadIpAllowed)
	}
	if len(Conf.UploadIpDeny) > 0 {
		fmt.Printf(sformat, "UploadIpDeny:", Conf.UploadIpDeny)
	}

	if len(Conf.UploadIpDeny) > 0 && len(Conf.UploadIpAllowed) > 0 {
		log.Panic("uploadDeny and uploadAllowed and not be set at same time, please  delete uploadDeny or uploadAllowed ")
	}

	if Conf.CorsAllow {
		fmt.Printf("%-30s%-20t\n", "CorsAllow:", Conf.CorsAllow)
	}

	if Conf.ImageUrlPrefix != "" {
		fmt.Printf(sformat, "ImageUrlPrefix:", Conf.ImageUrlPrefix)
	} else {
		log.Panic("ImageUrlPrefix  must be set")
	}

	if Conf.UploadKey != "" {
		fmt.Printf(sformat, "UploadKey:", Conf.UploadKey)
	} else {
		Conf.UploadKey = "file"
		fmt.Printf(sformat, "UploadKey(use default value):", Conf.UploadKey)
	}

	if Conf.FileCacheDir != "" {
		fmt.Printf(sformat, "FileCacheDir:", Conf.FileCacheDir)
	} else {
		Conf.FileCacheDir = "/var/fqimg/cache"
		fmt.Printf(sformat, "FileCacheDir(use default value): ", Conf.FileCacheDir)
	}

	if Conf.FileCacheSize != 0 {
		fmt.Printf(dformat, "FileCacheSize:", Conf.FileCacheSize)
	} else {
		Conf.FileCacheSize = 10000
		fmt.Printf(dformat, "FileCacheSize(use default value):", Conf.FileCacheSize)
	}

	if Conf.GetForceAction != "" {
		fmt.Printf(sformat, "GetForceAction:", Conf.GetForceAction)
	}
	if Conf.UploadThrottlePerIp != 0 {
		fmt.Printf(dformat, "UploadThrottlePerIp:", Conf.GetForceAction)
	}

	if Conf.UploadThrottleTotal != 0 {
		fmt.Printf(dformat, "UploadThrottleTotal:", Conf.GetForceAction)
	}
}
