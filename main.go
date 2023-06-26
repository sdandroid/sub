package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Server struct {
	Alias         string
	Address       string
	Port          int
	Password      string
	Protocol      string
	EncryptMethod string
}

func main() {

	resp, _ := http.Get("")
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	decodeString, _ := base64.StdEncoding.DecodeString(string(body))

	split := strings.Split(string(decodeString), "\n")

	var list []Server

	for i := 0; i < len(split); i++ {
		list = append(list, convert(split[i]))
	}
	var ServerList = ""

	for i := 0; i < len(list); i++ {
		ServerList += list[i].Address
		if !(i == len(list)-1) {
			ServerList += ";"
		}

	}
	cmd := exec.Command("/etc/xray/change.sh")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "ServerList="+ServerList, "Port="+strconv.Itoa(list[0].Port))
	err := cmd.Run()
	if err != nil {
		panic(err)
		return
	}

}
func convert(str string) Server {
	index := strings.Index(str, "://")
	Protocol := str[0:index]
	server := Server{}
	server.Protocol = Protocol

	if Protocol == "ss" {
		hashIndex := strings.Index(str, "#")
		serverStr, _ := base64.RawStdEncoding.DecodeString(str[index+3 : hashIndex])
		severSplit := strings.Split(string(serverStr), "@")
		split1 := strings.Split(severSplit[0], ":")
		split2 := strings.Split(severSplit[1], ":")
		server.EncryptMethod = split1[0]
		server.Password = split1[1]
		server.Address = split2[0]
		port, _ := strconv.Atoi(split2[1])
		server.Port = port
		server.Alias = str[hashIndex+1:]

	} else if Protocol == "vmess" {
		serverStr, _ := base64.RawStdEncoding.DecodeString(str[index+3:])

		jsonMap := map[string]any{}

		err := json.Unmarshal(serverStr, &jsonMap)
		if err != nil {
			panic(err)
			return Server{}
		}

		server.EncryptMethod = "auto"
		server.Password = jsonMap["id"].(string)
		server.Address = jsonMap["add"].(string)
		port, _ := strconv.Atoi(jsonMap["port"].(string))
		server.Port = port
		server.Alias = jsonMap["ps"].(string)

	}
	return server
}
