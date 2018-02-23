package plugin

import (
	"net/http"
	"github.com/hashicorp/golang-lru"
	"log"
	"os"
	"io"
	"encoding/base64"
	"github.com/loveczp/fqimg/lib"
)

type fileResponseWriter struct {
	file  io.Writer
	resp  http.ResponseWriter
	multi io.Writer
}

func newFileResponseWriter(file io.Writer, resp http.ResponseWriter) http.ResponseWriter {
	multi := io.MultiWriter(file, resp)
	return &fileResponseWriter{
		file:  file,
		resp:  resp,
		multi: multi,
	}
}

// implement http.ResponseWriter
// https://golang.org/pkg/net/http/#ResponseWriter
func (w *fileResponseWriter) Header() http.Header {
	return w.resp.Header()
}

func (w *fileResponseWriter) Write(b []byte) (int, error) {
	return w.multi.Write(b)
}

func (w *fileResponseWriter) WriteHeader(i int) {
	w.resp.WriteHeader(i)
}

type imageCacheItem struct {
	key      string
	filePath string
}

func removeFile(key interface{}, value interface{}) {
	citem := value.(imageCacheItem)
	log.Println("romve cache item :", citem.key);
	go func(filePath string) {
		//log.Println("romve cache item go routin in :",filePath);
		if _, err := os.Stat(filePath); err == nil {
			err := os.Remove(filePath)
			if err != nil {
				log.Panic("remove cache item error:", err.Error())
			}
		} else {
			log.Println("romve cache item ,file to be removed is not exsit :", filePath);
		}

	}(citem.filePath)
}

func Plugin_get_filecache(h http.HandlerFunc, conf lib.Config) http.HandlerFunc {
	var fileCache *(lru.Cache)
	if fileCache == nil {
		var err error;
		fileCache, err = lru.NewWithEvict(conf.FileCacheSize, removeFile)
		if err != nil {
			log.Panic("cache create error :", err)
		}
	}

	if len(conf.FileCacheDir) > 2 {
		os.MkdirAll(conf.FileCacheDir, 0777);
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if item, ok := (*fileCache).Get(r.URL.String()); ok {
			citem := item.(imageCacheItem)
			if cfile, err := os.Open(citem.filePath); err == nil {
				io.Copy(w, cfile);
				log.Println("data from file cache:", r.URL.String());
				defer cfile.Close()
				return;
			} else {
				(*fileCache).Remove(r.URL.String())
				log.Println("data error from file cache:", r.URL.String());
				h.ServeHTTP(w, r)
			}
		} else {
			cPath := conf.FileCacheDir + base64.StdEncoding.EncodeToString([]byte(r.URL.String()));
			var tempFile *(os.File)
			if _, err := os.Stat(cPath); !os.IsExist(err) {
				tempFile, err = os.Create(cPath);
				defer (*tempFile).Close()
			}
			log.Println("cache set data , key ", r.URL.String())
			citem := imageCacheItem{filePath: cPath, key: r.URL.String()}
			(*fileCache).Add(r.URL.String(), citem)
			mw := newFileResponseWriter(tempFile, w)
			h.ServeHTTP(mw, r)
		}
	}
}
