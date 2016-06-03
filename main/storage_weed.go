package main

import (
	"io"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"mime/multipart"
	"os"
	"strconv"
	"math/rand"
	"strings"
	"errors"
)

type WeedStore struct {
	masterUrl string
}






func (weed WeedStore ) storagePut(src io.Reader) (string,error){

	//{"count":1,"fid":"3,01637037d6","url":"127.0.0.1:8080","publicUrl":"localhost:8080"}
	resp, err := http.Get(weed.masterUrl+"/dir/assign")
	if err != nil {
		log.Fatalln("failed to create temp file:" + err.Error())
		return "",err
	}
	defer resp.Body.Close()
	assignBody, err := ioutil.ReadAll(resp.Body)


	log.Println("assignBody:",string(assignBody))
	var dat map[string]interface{}
	if err := json.Unmarshal(assignBody, &dat); err != nil {
		panic(err)
	}

	fid := dat["fid"].(string)
	url := dat["url"].(string)
	fullurl := "http://"+url+"/"+fid



	tempfileName := os.TempDir() + "/"+strconv.Itoa(rand.Int())
	fw, err := os.Create(tempfileName)
	if err != nil {
		log.Fatalln("failed to create temp file:" + err.Error())
		return "",err
	}

	_, err = io.Copy(fw, src);
	if err != nil {
		log.Fatalln(tempfileName, err.Error())
		return "",err
	}

	fw.Close()



















	reqBody := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(reqBody)
	fileWriter, err := multipartWriter.CreateFormFile("file",tempfileName)
	if err != nil {
		log.Fatalln(err.Error())
		return  "",err
	}

	// open file handle
	fh, err := os.Open(tempfileName)
	if err != nil {
		fmt.Println("error opening temp file")
		return "",err
	}


	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		log.Fatalln(err.Error())
		return  "",err
	}
	contentType := multipartWriter.FormDataContentType()
	err = multipartWriter.Close()

	//fullurl="http://127.0.0.1:12345/test"
	finalResp ,err:=http.Post( fullurl,contentType, reqBody)
	if err != nil {
		log.Fatalln(err.Error())
		return  "",err
	}



	finalRespBody, err := ioutil.ReadAll(finalResp.Body)
	log.Println("finalRespBody:",string(finalRespBody))
	if err != nil {
		log.Fatalln("finalRespBody Error:",err.Error())
		return  "",err
	}

	defer os.Remove(tempfileName)


	return fid,err


}







func (weed WeedStore ) storageGet(key string) (io.Reader, error) {
	//http://localhost:9333/dir/lookup?volumeId=3
	//{"locations":[{"publicUrl":"localhost:8080","url":"localhost:8080"}]}

	volumeId := key[0:strings.Index(key,",")]
	resp, err := http.Get(weed.masterUrl+"/dir/lookup?volumeId="+volumeId)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	assignBody, err := ioutil.ReadAll(resp.Body)


	log.Println(string(assignBody))

	var dat map[string]interface{}
	if err := json.Unmarshal(assignBody, &dat); err != nil {
		log.Fatalln(err)
	}

	locations := dat["locations"].([]interface{})

	location:=locations[0].( map[string]interface{})
	publicUrl:=location["publicUrl"].(string)

	fullurl := "http://"+publicUrl+"/"+key


	finalResp, err := http.Get(fullurl)
	if err != nil {
		log.Fatalln("get pic data from :",fullurl,". error occur. error:"+err.Error())
	}



	return finalResp.Body,nil
}


func initWeed(config  Config) ( WeedStore ,error){
	//hasUrlGet := config.Get("weed.master_url")
	if(config.WeedMasterUrl==""){
		log.Panic("weed.master_url does not exsit!")

		return WeedStore{},errors.New("weed.master_url does not exsit!");
	}
	//master_url:=config.Get("weed.master_url").(string)
	fmt.Printf(sformat,"weed.master_url:",config.WeedMasterUrl)

	return WeedStore{masterUrl: config.WeedMasterUrl},nil
}