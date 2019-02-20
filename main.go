package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
)

// 接收数据格式
type res struct {
	Command []string `json:"command"`
}

func creatSendText(errCode int, message string) []byte {
	// 发送数据格式
	type send struct {
		Error   int    `json:"error"`
		Message string `json:"message"`
	}

	var sendData send
	sendData.Error = errCode
	sendData.Message = message
	sendText, _ := json.Marshal(sendData)
	return sendText
}

func runCommand(w http.ResponseWriter, r *http.Request) {
	var rasData res
	log.Println("[" + r.RemoteAddr + "] request url: " + r.RequestURI)

	// 取出Post过来的数据
	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write(creatSendText(1, "解析数据出错!"))
		return
	}
	log.Println("[" + r.RemoteAddr + "] send message: " + string(result))

	// 格式化Json数据
	if err := json.Unmarshal(result, &rasData); err != nil {
		w.Write(creatSendText(1, "请求数据不是正确的Json格式!"))
		return
	}
	if len(rasData.Command) == 0 {
		log.Println("数据中不包含Command,尝试从url中解析命令!")
		urlCommand := strings.Split(r.URL.RawQuery, "&")
		for ind, base64Data := range urlCommand {
			decoded, err := base64.StdEncoding.DecodeString(base64Data)
			if err != nil {
				w.Write(creatSendText(1, "参数["+string(base64Data)+"]不是正确的base64编码数据!"))
				log.Println("decode error: ", base64Data)
				return
			}
			urlCommand[ind] = string(decoded)
		}
		if len(urlCommand) == 0 {
			log.Println("找不到Command，忽略请求!")
			return
		}
		rasData.Command = urlCommand
	}

	log.Println("执行命令组:", rasData.Command)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		comm := strings.Join(rasData.Command, " & ")
		log.Println(comm)
		cmd = exec.Command("cmd", "/c", comm)
	case "linux":
		comm := strings.Join(rasData.Command, ";")
		log.Println(comm)
		cmd = exec.Command("sh", "-c", comm)
	}

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err = cmd.Run()
	if err != nil {
		w.Write(creatSendText(1, "无法运行命令!"))
		return
	}
	w.Write(creatSendText(0, out.String()))
}

func main() {
	http.HandleFunc("/hook", runCommand)
	log.Printf("server is running at 0.0.0.0:%s", "8888")
	http.ListenAndServe(":8888", nil)
}
