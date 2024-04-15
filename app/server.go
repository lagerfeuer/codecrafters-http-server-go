package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	httpVersion = "HTTP/1.1"
	endOfLine   = "\r\n"

	httpStatusCodes = map[int]string{
		200: "OK",
		404: "Not Found",
	}
)

// unused
type Request struct {
	method  string
	uri     string
	version string
	headers map[string]string
}

type Response struct {
	statusCode int
	headers    map[string]string
	body       string
}

func NewResponse() *Response {
	return &Response{
		statusCode: 404,
		headers:    map[string]string{},
		body:       "",
	}
}

func (r *Response) ToBytes() []byte {
	var sb strings.Builder
	sb.WriteString(
		fmt.Sprintf("%s %d %s%s", httpVersion, r.statusCode, httpStatusCodes[r.statusCode], endOfLine),
	)

	body := []byte(r.body)
	if r.body != "" {
		r.headers["Content-Type"] = "text/plain"
		r.headers["Content-Length"] = strconv.Itoa(len(body))
	}

	for key, val := range r.headers {
		sb.WriteString(fmt.Sprintf("%s: %s%s", key, val, endOfLine))
	}
	sb.WriteString(endOfLine)

	if r.body != "" {
		sb.WriteString(r.body)
	}

	return []byte(sb.String())
}

func handle(conn net.Conn) error {
	buffer := make([]byte, 1024)
	defer conn.Close()

	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading to connection: ", err.Error())
		return err
	}

	response := NewResponse()

	content := string(buffer)
	headers := strings.Split(content, endOfLine)
	startLine := strings.Split(headers[0], " ")

	if startLine[1] == "/" {
		response.statusCode = 200
	}
	if strings.HasPrefix(startLine[1], "/echo/") {
		response.statusCode = 200
		response.body = startLine[1][len("/echo/"):]
	}

	_, err = conn.Write(response.ToBytes())
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
