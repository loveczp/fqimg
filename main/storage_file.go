package main

import (
	"io"
	"log"
	"os"
	"gopkg.in/ini.v1"
	"errors"
)

type FileStore struct {
	dir string
}

func (file FileStore ) storagePut(md5 string, src io.Reader) error {
	cDirs, err := file.getCompositDirs(md5);

	if err != nil {
		return err;
	}

	finalDir := file.dir + "/" + cDirs.first + "/" + cDirs.second;
	os.MkdirAll(finalDir, 0777);

	fileName := finalDir + "/" + md5

	if _, error := os.Stat(fileName); error != nil {

		fw, err := os.Create(fileName)
		if err != nil {
			log.Println("failed to create file in file storage:" + fileName + err.Error())
			return err
		}
		_, err = io.Copy(fw, src);

		if err != nil {
			log.Println("failed to copy file to file system:" + fileName, err.Error())
			return err;
		}

		fw.Close()
		return nil

	}else {
		return nil
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

func initFile(config *ini.File) FileStore {
	dir := config.Section("").Key("file.dir").MustString("/var/go_image_server")
	os.MkdirAll(dir, 0777);
	return FileStore{dir: dir}
}