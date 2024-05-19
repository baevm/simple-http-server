package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"strconv"
	"strings"
)

type Headers map[string][]string

type Request struct {
	Method  string
	Uri     string
	Proto   string
	Headers Headers
	Body    []byte
}

type Response struct {
	Proto      string
	StatusCode int
	StatusText string
	Headers    Headers
	Body       string
}

var SEPARATOR = []byte("\r\n")
var NEW_LINE = byte('\n')

var HEADER_SEPARATOR = []byte(": ")

func NewRequest(conn net.Conn) *Request {
	reader := bufio.NewReader(conn)
	req := &Request{
		Headers: make(Headers),
	}

	/* Read request line */
	req_line, err := reader.ReadString(NEW_LINE)

	if err != nil {
		return &Request{}
	}

	req_line_slice := strings.Fields(req_line)

	if len(req_line_slice) < 3 {
		return &Request{}
	}

	req.Method = req_line_slice[0]
	req.Uri = req_line_slice[1]
	req.Proto = req_line_slice[2]

	/* Read headers */
	for {
		line, _, err := reader.ReadLine()

		if err != nil {
			break
		}

		if len(line) == 0 {
			break
		}

		header := bytes.Split(line, HEADER_SEPARATOR)
		headerKey := string(header[0])
		headerValue := string(header[1])

		req.Headers[headerKey] = strings.Split(headerValue, ", ")
	}

	if req.Method == "GET" || req.Method == "HEAD" {
		return req
	}

	/* Read body */
	contentLengthStr, isOk := req.Headers["Content-Length"]

	if !isOk {
		return req
	}

	contentLength, err := strconv.Atoi(contentLengthStr[0])

	if err != nil {
		return req
	}

	if contentLength == 0 {
		return req
	}

	body := make([]byte, contentLength)
	_, err = io.ReadFull(reader, body)

	if err != nil {
		return req
	}

	req.Body = body

	return req
}
