package lib

import (
	"io"
	"net/http"
	"log"
	"encoding/json"
	"github.com/loveczp/fqimg/store"
	"github.com/pkg/errors"
	"strconv"
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

		if len(files)>Conf.UploadFileNmuLimit{
			WriteErr(res, http.StatusBadRequest, errors.New("number of uploaded file exceed the limit :"+strconv.Itoa(Conf.UploadFileNmuLimit)))
			return
		}

		for _,f := range files {
			if f.Size > int64(Conf.UploadFileSizeLimit*1024){
				WriteErr(res, http.StatusBadRequest, errors.New("size of uploaded file exceed the limit :"+strconv.Itoa(Conf.UploadFileSizeLimit)+"kb"))
				return
			}
		}


		var md5List []string
		for _,f := range files {
			tfile, _ := f.Open();
			key, err := store.Put(tfile)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				log.Fatalln("error ocurr when store to file", err)
			}
			md5List = append(md5List, Conf.ImageUrlPrefix+ getAlias+"/"+key);
		}


		restring, _ := json.Marshal(md5List);
		res.Header().Add("Content-Type", "application/json")
		io.WriteString(res, string(restring))
	}
}
