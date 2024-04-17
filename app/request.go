package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Request struct {
	method  string
	uri     string
	version string
	headers map[string]string
	body    []byte
}

func ParseRequest(raw []byte) Request {
	content := string(raw)
	parts := strings.Split(content, endOfLine+endOfLine)
	rawHeaders := strings.Split(parts[0], endOfLine)
	startLine := strings.Split(rawHeaders[0], " ")

	headers := make(map[string]string, 8)
	rawBody := parts[1]

	for _, line := range rawHeaders[1:] {
		splits := strings.SplitN(line, ":", 2)
		key := splits[0]
		val := strings.Trim(splits[1], " \r\n")
		headers[key] = val
	}

	body := []byte{}
	if headers["Content-Length"] != "" {
		length, err := strconv.Atoi(headers["Content-Length"])
		if err != nil {
			fmt.Println("Error converting 'Content-Length': ", err.Error())
			os.Exit(1)
		}
		body = []byte(rawBody)[:length]
	}

	return Request{
		method:  startLine[0],
		uri:     startLine[1],
		version: startLine[2],
		headers: headers,
		body:    body,
	}
}
