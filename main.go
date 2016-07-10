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
	http.HandleFunc("/uploadm", imageserverlib.UploadMultiHandler)
	http.HandleFunc("/favicon.ico", imageserverlib.FavHandle)
	http.HandleFunc("/", imageserverlib.GetHandler)
	http.HandleFunc("/hello", imageserverlib.HelloHandle)
	http.HandleFunc("/test", imageserverlib.UploadTestHandler)
	http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir("/pictest"))))

	err := http.ListenAndServe(":" + strconv.Itoa(imageserverlib.Conf.Port), nil)
	if err != nil {
		log.Panic("server start failed :", err)
		return
	} else {
		log.Printf("server start at port: " + strconv.Itoa(imageserverlib.Conf.Port))
	}
}




