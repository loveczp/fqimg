package main

import (
	"io"
	"os"
	"net/http"
	"log"
	"encoding/json"
	"fmt"
)

func helloHandle(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world!")
}

func handleFav(w http.ResponseWriter, r *http.Request)  {
	if ffile ,err :=os.Open(favicon_path); err==nil{
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
		tfile,_:=files[i].Open();
		key,err:=store.storagePut(tfile)
		if err!=nil{
			log.Fatalln("error ocurr when store to file",err)
		}
		md5List=append(md5List,key);
	}
	restring, _ := json.Marshal(md5List);
	//log.Println(md5List.Len())
	io.WriteString(res, string(restring))
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

	key,err:=store.storagePut(req.Body)
	if err !=nil{
		log.Println("failed to save to the storage system:" + err.Error())
		io.WriteString(res,"failed to save to the storage system:" + err.Error())
		return
	}


	if err !=nil{
		log.Println("failed to remove fw2:" + err.Error())
		io.WriteString(res,"failed to remove fw2:" + err.Error())
		return
	}

	var remap map[string]string
	remap = map[string]string{"key": key, "msg": "ok"}
	restring, _ := json.Marshal(remap);
	io.WriteString(res, string(restring))
}



func uploadTestHandler(res http.ResponseWriter, req *http.Request) {
	data,_:=json.Marshal( req.Header)
	fmt.Println("header:",string(data))
	fmt.Println("body:",req.Body)
	io.WriteString(res, "ok")
	return
}

