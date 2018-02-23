package lib

import (
	"io"
	"net/http"
	"log"
	"encoding/json"
	"github.com/loveczp/fqimg/store"
	"github.com/pkg/errors"
)

func UploadHandler(store store.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		err := req.ParseMultipartForm(1024);
		if err != nil {
			WriteErr(res, http.StatusInternalServerError, err)
			return
		}

		files := req.MultipartForm.File[Conf.UploadKey]

		if len(files) == 0 {
			WriteErr(res, http.StatusBadRequest, errors.New("found no image from the form with key:"+Conf.UploadKey))
			return
		}
		var md5List []string
		for i, _ := range files {
			tfile, _ := files[i].Open();
			key, err := store.Put(tfile)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				log.Fatalln("error ocurr when store to file", err)
			}
			md5List = append(md5List, Conf.ImageUrlPrefix+"pic/"+key);
		}
		restring, _ := json.Marshal(md5List);
		res.Header().Add("Content-Type", "application/json")
		io.WriteString(res, string(restring))
	}
}
