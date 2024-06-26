package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	httpVersion = "HTTP/1.1"
	endOfLine   = "\r\n"

	httpStatusCodes = map[int]string{
		200: "OK",
		201: "Created",
		404: "Not Found",
	}
)

func readFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	return content, nil
}

func writeFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0666)
}

func handleGet(request Request, path string) Response {
	response := NewResponse()

	if request.uri == "/" {
		response.statusCode = 200
	} else if strings.HasPrefix(request.uri, "/echo/") {
		response.statusCode = 200
		response.headers["Content-Type"] = "text/plain"
		response.body = []byte(request.uri[len("/echo/"):])
	} else if request.uri == "/user-agent" {
		response.statusCode = 200
		response.headers["Content-Type"] = "text/plain"
		response.body = []byte(request.headers["User-Agent"])
	} else if strings.HasPrefix(request.uri, "/files/") {
		filename := request.uri[len("/files/"):]
		content, err := readFile(path + "/" + filename)
		if err != nil {
			response.statusCode = 404
		} else {
			response.statusCode = 200
			response.headers["Content-Type"] = "application/octet-stream"
			response.body = content
		}
	}

	return response
}

func handlePost(request Request, path string) Response {
	response := NewResponse()

	if strings.HasPrefix(request.uri, "/files/") {
		filename := request.uri[len("/files/"):]
		err := writeFile(path+"/"+filename, request.body)
		if err != nil {
			response.statusCode = 500
			response.body = []byte(err.Error())
		} else {
			response.statusCode = 201
		}
	}

	return response
}

func handle(conn net.Conn, path string) error {
	buffer := make([]byte, 1024)
	defer conn.Close()

	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading to connection: ", err.Error())
		return err
	}

	request := ParseRequest(buffer)
	response := Response{}
	if request.method == "GET" {
		response = handleGet(request, path)
	} else if request.method == "POST" {
		response = handlePost(request, path)
	}

	_, err = conn.Write(response.ToBytes())
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		return err
	}

	return nil
}

func main() {
	args := os.Args
	path := ""
	if len(args) == 3 && args[1] == "--directory" {
		path = args[2]
	}

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221: ", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handle(conn, path)
	}
}
