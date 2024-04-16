package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Response struct {
	statusCode int
	headers    map[string]string
	body       string
}

func NewResponse() *Response {
	return &Response{
		statusCode: 404,
		headers:    make(map[string]string, 8),
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
