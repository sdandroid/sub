package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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

	args := os.Args

	var ServerName = args[1]
	var url = args[2]
	var index, _ = strconv.Atoi(args[3])

	if len(strings.TrimSpace(url)) == 0 {
		os.Exit(0)
	}
	resp, e := http.Get(url)
	if e != nil {
		panic(e)
	}
	if resp == nil {
		fmt.Println("body is nil")
		os.Exit(0)
	}
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
	var addrBool = false
	var portBool = false
	var passwordBool = false
	for i := 0; i < len(list); i++ {
		var newIndex = strconv.Itoa(index + i)
		oldServer := execCommand("/sbin/uci", "get", ServerName+".@servers["+newIndex+"].server")
		//fmt.Printf("oldServer %v  %T \n", oldServer, oldServer)
		//fmt.Printf("newServer:  %v %T \n", list[i].Address, list[i].Address)
		//fmt.Printf("bool %t \n", oldServer == list[i].Address)
		//fmt.Printf("bool2 %t \n", oldServer != list[i].Address)
		//printStringByte(oldServer)
		//printStringByte(list[i].Address)
		if list[i].Address != oldServer {
			addrBool = true
			execCommand("/sbin/uci", "set", ServerName+".@servers["+newIndex+"].server="+list[i].Address)
		}
		oldPortStr := execCommand("/sbin/uci", "get", ServerName+".@servers["+newIndex+"].server_port")
		var oldPort, _ = strconv.Atoi(oldPortStr)
		if list[i].Port != oldPort {
			portBool = true
			var newPort = strconv.Itoa(list[i].Port)
			execCommand("/sbin/uci", "set", ServerName+".@servers["+newIndex+"].server_port="+newPort)
		}
		oldPassword := execCommand("/sbin/uci", "get", ServerName+".@servers["+newIndex+"].password")
		if list[i].Password != oldPassword {
			passwordBool = true
			execCommand("/sbin/uci", "set", ServerName+".@servers["+newIndex+"].password="+list[i].Password)
		}
	}
	fmt.Printf("addrBool %t portBool %t passwordBool %t \n", addrBool, portBool, passwordBool)
	if addrBool || portBool || passwordBool {
		execCommand("/sbin/uci", "commit", ServerName)
		execCommand("/etc/init.d/"+ServerName, "restart")
	}

	//var ServerList = ""
	//
	//for i := 0; i < len(list); i++ {
	//	ServerList += list[i].Address
	//	if !(i == len(list)-1) {
	//		ServerList += ";"
	//	}
	//
	//}
	//cmd := exec.Command(cmdStr)
	//cmd.Env = os.Environ()
	//cmd.Env = append(cmd.Env, "ServerList="+ServerList, "Port="+strconv.Itoa(list[0].Port), "Service="+ServerName)
	//result, err := cmd.Output()
	//if err != nil {
	//	panic(err)
	//	return
	//}
	//fmt.Println(string(result))
}
func printStringByte(str string) {
	b := []byte(str)
	fmt.Println(b)
}

func execCommand(command string, arg ...string) string {
	//fmt.Println("command: " + command)
	//for _, arg := range arg {
	//	fmt.Println("arg: " + arg)
	//}
	cmd := exec.Command(command, arg...)
	oldServerByte, err := cmd.Output()
	if err != nil {
		panic(err)
		return ""
	}
	if command == "/sbin/uci" && arg[0] == "get" {
		oldServerByte = oldServerByte[0 : len(oldServerByte)-1]
		//fmt.Println(oldServerByte)
		//return  string(oldServerByte)
	}
	var oldServer = string(oldServerByte)
	//println("result: " + oldServer)
	return oldServer
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
