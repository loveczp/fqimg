package main

import (
	"net/http"
	"strconv"
	"log"
	_ "net/http/pprof"
	"github.com/loveczp/go_image_server/imageserverlib"
)

func main() {
	http.HandleFunc("/upload", imageserverlib.UploadBinHandler)
	http.HandleFunc("/uploadm", UploadMultiHandler)
	http.HandleFunc("/favicon.ico", FavHandle)
	http.HandleFunc("/", GetHandler)
	http.HandleFunc("/hello", HelloHandle)
	http.HandleFunc("/test", UploadTestHandler)
	http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir("/pictest"))))

	err := http.ListenAndServe(":" + strconv.Itoa(conf.Port), nil)
	if err != nil {
		log.Panic("server start failed :", err)
		return
	} else {
		log.Printf("server start at port: " + strconv.Itoa(Conf.Port))
	}
}




