package store

import (
	"io"
	"fmt"
	"errors"
	"github.com/weilaihui/fdfs_client"
	"os"
	"strconv"
	"log"
	"math/rand"
	"bytes"
	"strings"
)

type FastStore struct {
	client fdfs_client.FdfsClient
}

func (fast FastStore) Put(src io.Reader) (string, error) {
	tempfileName := os.TempDir() + "/" + strconv.Itoa(rand.Int())
	fw, err := os.Create(tempfileName)
	if err != nil {
		log.Fatalln("failed to create temp file:" + err.Error())
		return "", err
	}

	_, err = io.Copy(fw, src);
	if err != nil {
		log.Fatalln(tempfileName, err.Error())
		return "", err
	}

	fw.Close()
	defer os.Remove(tempfileName)
	resp, err := fast.client.UploadByFilename(tempfileName)
	log.Println("fastdfs upload resp.RemoteFileId:", resp.RemoteFileId)
	log.Println("fastdfs upload resp.GroupName:", resp.GroupName)

	//
	key := strings.Replace(resp.RemoteFileId, "\\", "/", 1)
	//key := strings.Replace(resp.RemoteFileId, "\\", "", 0)

	/*	escapeId,err:=url.QueryUnescape( resp.RemoteFileId)
		log.Println("escapeId:",escapeId)
		if err!=nil{
			log.Fatalln("url.QueryUnescape( resp.RemoteFileId):", err.Error())
			return "",err
		}*/
	return key, err;
}

func (fast FastStore) Get(key string) (io.Reader, error) {
	/*	log.Println("fastdfs download input  key:",key)

		unescapedId,err:= url.QueryUnescape(key);
		log.Println("fastdfs download input  unescaped key:",unescapedId)*/
	resp, err := fast.client.DownloadToBuffer(key, 0, 0)
	if err != nil {
		//log.Fatalln("can not get file from fastdfs key:",key)
		log.Println("fastdfs download error:", err.Error())
		return nil, err;
	}
	bys := resp.Content.([]byte)
	reader := bytes.NewReader(bys)
	return reader, nil
}

func InitFast(fastPath string) (FastStore, error) {
	//hasUrl := config.Get("fastdfs.config_file_path")
	if (fastPath == "") {
		panic("fastdfs.configPath does not exsit!")
		return FastStore{}, errors.New("fastdfs.config_file_path does not exsit!");
	}
	//fastdfsPath:=config.Get("fastdfs.config_file_path").(string)
	fmt.Printf("%-20s%-20s\n", "fastdfs.config_file_path:", fastPath)
	c, err := fdfs_client.NewFdfsClient(fastPath)
	if err != nil {
		panic("server start failed :" + err.Error())
	}

	return FastStore{client: *c}, err
}
