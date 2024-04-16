package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Response struct {
	statusCode int
	headers    map[string]string
	body       []byte
}

func NewResponse() *Response {
	return &Response{
		statusCode: 404,
		headers:    make(map[string]string, 8),
		body:       []byte{},
	}
}

func (r *Response) ToBytes() []byte {
	var sb strings.Builder
	sb.WriteString(
		fmt.Sprintf("%s %d %s%s", httpVersion, r.statusCode, httpStatusCodes[r.statusCode], endOfLine),
	)

	if len(r.body) > 0 {
		r.headers["Content-Length"] = strconv.Itoa(len(r.body))
	}

	for key, val := range r.headers {
		sb.WriteString(fmt.Sprintf("%s: %s%s", key, val, endOfLine))
	}
	sb.WriteString(endOfLine)

	if len(r.body) > 0 {
		sb.Write(r.body)
	}

	return []byte(sb.String())
}
