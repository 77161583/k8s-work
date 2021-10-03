package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"runtime"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		//3.Server 端记录访问日志包括客户端 IP
		LocalIP := GetLocalIP()
		fmt.Printf("client IP: %s\n", LocalIP)
		fmt.Printf("Request URL: %s\n", r.URL)
		fmt.Printf("User Agent: %s\n", r.Header.Get("User-Agent"))
		fmt.Printf("Request Header: %v\n", r.Header)

		//1.接收客户端 request，并将 request 中带的 header 写入 response header
		len := r.ContentLength           // 获取请求实体长度
		body := make([]byte, len)        // 创建存放请求实体的字节切片
		r.Body.Read(body)                // 调用 Read 方法读取请求实体并将返回内容存放到上面创建的字节切片
		io.WriteString(rw, string(body)) // 将请求实体作为响应实体返回
		rw.Header().Set("header_Data", string(body))

		//2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
		rw.Header().Set("version", runtime.Version())
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
