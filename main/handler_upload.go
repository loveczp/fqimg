package main

import (
	"io"
	"os"
	"net/http"
	"strconv"
	"math/rand"
	"log"
	"encoding/hex"
	"crypto/md5"
	"io/ioutil"
	"encoding/json"
)

func helloHandle(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world!")
}

func handleFav(w http.ResponseWriter, r *http.Request)  {
	if ffile ,err :=os.Open(favicoPath); err==nil{
		io.Copy(w,ffile);
	}
	return ;
}

const uploadhtml =
`<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<html>
    <head>
        <script>
            function changetext()
            {
                var newItemText=document.createTextNode("Choose file:")
                var newItemInput=document.createElement("input")
                newItemInput.name="userfile"
                newItemInput.type="file"
                var addfile=document.getElementById("bt_addfile")
                var submit=document.getElementById("bt_submit")
                var newItemBr=document.createElement("br")
                var myform=document.getElementById("upform")
                myform.appendChild(newItemText);
                myform.appendChild(newItemInput);
                myform.appendChild(addfile);
                myform.appendChild(newItemBr);
                myform.appendChild(submit);
            }
        </script>
    </head>
    <h1>Welcome to go_image_server world!</h1>
    <p>Upload image(s) to go_image_server:</p>
    <form enctype="multipart/form-data" action="uploadm" method=post target=_blank id="upform">
        Choose file:<input name="userfile" type="file">
        <input type="button" value="+" onclick="changetext()" id="bt_addfile">
        </br>
        <input type="submit" value="upload" id="bt_submit">
    </form>
</html>`

func uploadMultiHandler(res http.ResponseWriter, req *http.Request)  {
	err:=req.ParseMultipartForm(1024);
	if err!=nil{
		io.WriteString(res,"error when parse multipart form file");
		return
	}

	files:=req.MultipartForm.File["userfile"]
	var md5List []string
	for i, _ := range files {
		tfile,err:=files[i].Open();
		Buf, err := ioutil.ReadAll(tfile)
		if err != nil {
			log.Println(files[i], err)
			io.WriteString(res,"failed to read temp file:" + err.Error())
			return
		}
		md5value := hex.EncodeToString(byte2string(md5.Sum(Buf)))
		log.Println(md5value);
		tfile.Close();


		//l[i]=(md5value);
		md5List=append(md5List,md5value);
		//log.Println(l.Len())

		fw2, err := files[i].Open();
		if err !=nil{
			log.Println("failed to open saved temp file:" + err.Error())
			io.WriteString(res,"failed to open saved temp file:" + err.Error())
			return
		}

		err=store.storagePut(md5value,fw2)
		if err !=nil{
			log.Println("failed to save to the storage system:" + err.Error())
			io.WriteString(res,"failed to save to the storage system:" + err.Error())
			return
		}


		err=fw2.Close()
		if err !=nil{
			log.Println("failed to close fw2:" + err.Error())
			io.WriteString(res,"failed to close fw2:" + err.Error())
			return
		}
	}



	restring, _ := json.Marshal(md5List);
	//log.Println(md5List.Len())
	io.WriteString(res, string(restring))
}

type reValue struct{
	key string
	msg string;
}

func uploadBinHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method =="GET" {
		res.Header().Add("Content-Type","text/html; charset=utf-8")
		io.WriteString(res,uploadhtml)
		return
	}

	if store == nil {
		log.Println("store entity is null");
		io.WriteString(res,"store entity is null")
		return
	}

	tempfileName := os.TempDir() + "/"+strconv.Itoa(rand.Int())
	fw, err := os.Create(tempfileName)
	if err != nil {
		log.Println("failed to create temp file:" + err.Error())
		io.WriteString(res,"failed to create temp file:" + err.Error())
		return
	}

	_, err = io.Copy(fw, req.Body);
	if err != nil {
		log.Println(tempfileName, err.Error())
		io.WriteString(res,"failed to copy temp file:" + err.Error())
		return
	}
	fw.Close()


	Buf, err := ioutil.ReadFile(tempfileName)
	if err != nil {
		log.Println(tempfileName, err)
		io.WriteString(res,"failed to read temp file:" + err.Error())
		return
	}

	md5value := hex.EncodeToString(byte2string(md5.Sum(Buf)))

	fw2, err := os.Open(tempfileName)
	if err !=nil{
		log.Println("failed to open saved temp file:" + err.Error())
		io.WriteString(res,"failed to open saved temp file:" + err.Error())
		return
	}

	err=store.storagePut(md5value,fw2)
	if err !=nil{
		log.Println("failed to save to the storage system:" + err.Error())
		io.WriteString(res,"failed to save to the storage system:" + err.Error())
		return
	}


	err=fw2.Close()
	if err !=nil{
		log.Println("failed to close fw2:" + err.Error())
		io.WriteString(res,"failed to close fw2:" + err.Error())
		return
	}
	err = os.Remove(tempfileName)

	if err !=nil{
		log.Println("failed to remove fw2:" + err.Error())
		io.WriteString(res,"failed to remove fw2:" + err.Error())
		return
	}

	var remap map[string]string
	remap = map[string]string{"md5": md5value, "msg": "ok"}
	restring, _ := json.Marshal(remap);
	io.WriteString(res, string(restring))
}

