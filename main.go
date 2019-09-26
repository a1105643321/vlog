package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//输出hello world

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte("hello world"))
}

func main() {
	//实现读取文件handler
	fileHandler := http.FileServer(http.Dir("./video"))


	http.Handle("/video/", http.StripPrefix("/video/", fileHandler))

	http.HandleFunc("/api/upload",uploadHandler)

	http.HandleFunc("/api/list",getFileListHandler)

	//注册进servermux 就是将不同url的请求交给对应的handler处理
	http.HandleFunc("/sayHello", sayHello)

	//启动web
	http.ListenAndServe(":8090", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	//限制上传文件大小
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//获取文件
	file, fileHeader, err := r.FormFile("uploadFile")

	//检查文件类型
	ret := strings.HasSuffix(fileHeader.Filename,".mp4" )
	if ret==false{
		http.Error(w,"not mp4",http.StatusInternalServerError)
		return
	}

	//获取随机名称
	md5Byte:=md5.Sum([]byte(fileHeader.Filename+time.Now().String()))
	md5Str:=fmt.Sprintf("%x",md5Byte)
	newFileName:=md5Str+".mp4"

	//写入文件
	dst,err:=os.Create("./video/"+newFileName)
	defer dst.Close()
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	defer file.Close()
	if _,err:=io.Copy(dst,file);err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	return
}

func getFileListHandler(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	files,_:=filepath.Glob("video/*")
	var ret[]string
	for _,file:=range files{
		ret=append(ret,"http://"+r.Host+"/video/"+filepath.Base(file))
	}
	retJson,_:=json.Marshal(ret)
	w.Write(retJson)
	return
}
