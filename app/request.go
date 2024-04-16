package main

import (
	"strings"
)

type Request struct {
	method  string
	uri     string
	version string
	headers map[string]string
}

func ParseRequest(raw []byte) *Request {
	content := string(raw)
	lines := strings.Split(content, endOfLine)
	startLine := strings.Split(lines[0], " ")

	headers := make(map[string]string, 8)

	for _, line := range lines[1:] {
		if !strings.Contains(line, ":") {
			break
		}

		splits := strings.SplitN(line, ":", 2)
		key := splits[0]
		val := strings.Trim(splits[1], " \r\n")
		headers[key] = val
	}

	return &Request{
		method:  startLine[0],
		uri:     startLine[1],
		version: startLine[2],
		headers: headers,
	}
}
