package imageserverlib

import (
	"io"
	"log"
	"os"
	"errors"
	"fmt"
	"io/ioutil"
	"encoding/hex"
	"crypto/md5"
	"bytes"
)

type FileStore struct {
	dir string
}

func (fileStorage FileStore ) storagePut(src io.Reader) (string,error) {
	srcBuffer, err := ioutil.ReadAll(src)
	if err != nil {
		log.Fatalln("failed to read temp file")
		return "",err
	}
	md5value := hex.EncodeToString(byte2string(md5.Sum(srcBuffer)))
	log.Println(md5value);

	cDirs, err := fileStorage.getCompositDirs(md5value);
	if err != nil {
		return "",err
	}

	finalDir := fileStorage.dir + "/" + cDirs.first + "/" + cDirs.second;
	os.MkdirAll(finalDir, 0777);

	destFileName := finalDir + "/" + md5value

	if _, error := os.Stat(destFileName); error != nil {
		destFileWriter, err := os.Create(destFileName)
		if err != nil {
			log.Fatalln("failed to create file in file storage:" + destFileName + err.Error())
			return "",err
		}
		bufferReader := bytes.NewReader(srcBuffer)
		bytenum, err := io.Copy(destFileWriter, bufferReader);
		log.Println("storage_file copied byte:",bytenum)


		if err != nil {
			log.Fatalln("failed to copy file to file system:" + destFileName, err.Error())
			return "",err;
		}

		destFileWriter.Close()
		return md5value,nil

	}else {
		return md5value,nil
	}
}

func (file FileStore ) storageGet(md5 string) (io.Reader, error) {
	cDirs, err := file.getCompositDirs(md5);

	if err != nil {
		return nil, err;
	}
	finalDir := file.dir + "/" + cDirs.first + "/" + cDirs.second;
	fileName := finalDir + "/" + md5

	if md5Reader, error := os.Open(fileName); error == nil {
		return md5Reader, nil
	}else {
		return nil, error;
	}
}

//linux的ext3，4中存放超过1000个文件后，就会变得读取效率非常低。
type CompositDirs struct {
	first  string
	second string
	md5    string
}

func (file FileStore ) getCompositDirs(md5 string) (CompositDirs, error) {
	if (len(md5) < 5) {
		return CompositDirs{}, errors.New("this is a new error")
	}else {
		return CompositDirs{first:md5[0:2], second:md5[2:4], md5:md5}, nil
	}
}

func initFile(config  Config) FileStore {
	if(config.FileDir==""){
		config.FileDir="/var/go_image_server"
	}
	fmt.Printf(sformat,"file_dir:",config.FileDir)

	os.MkdirAll(config.FileDir, 0777);
	return FileStore{dir: config.FileDir}
}