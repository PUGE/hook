package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

var path = "C:\\Users\\my\\Documents\\GitHub\\OSS"

func hello(w http.ResponseWriter, r *http.Request) {
	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("read data error!")
	}
	log.Println("get message:", string(result))

	log.Println("cd " + path + ";git pull")
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("cmd", "/c", "cd "+path+" & dir")

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err = cmd.Run()
	if err != nil {
		log.Println(err)
		log.Println("Execution failure!")
	}

	w.Write([]byte(out.String()))
}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8001", nil)
}
