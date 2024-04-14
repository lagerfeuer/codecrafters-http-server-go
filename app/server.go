package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	httpVersion = "HTTP/1.1"
	endOfLine   = "\r\n\r\n"

	httpStatusCodes = map[int]string{
		200: "OK",
		404: "Not Found",
	}
)

func makeResponse(statusCode int) []byte {
	s := fmt.Sprintf("%s %d %s %s", httpVersion, statusCode, httpStatusCodes[statusCode], endOfLine)
	return []byte(s)
}

func handle(conn net.Conn) error {
	buffer := make([]byte, 1024)
	defer conn.Close()

	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading to connection: ", err.Error())
		return err
	}

	response := makeResponse(404)

	content := string(buffer)
	headers := strings.Split(content, endOfLine)
	startLine := strings.Split(headers[0], " ")

	if startLine[1] == "/" {
		response = makeResponse(200)
	}

	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		return err
	}

	return nil
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221: ", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	err = handle(conn)
	if err != nil {
		fmt.Println("Error while handling connection: ", err.Error())
		os.Exit(1)
	}
}
