package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net"
)

type HttpContext struct {
	req  *Request
	conn net.Conn
}

func (c HttpContext) write(response string) error {
	c.conn.Write([]byte(response))
	return c.conn.Close()
}

func (c HttpContext) String(status int, body string) error {
	response := fmt.Sprintf("HTTP/1.1 %d %v\r\n", status, StatusText(status))

	encoding, isFound := c.req.Headers["Accept-Encoding"]

	if isFound {
		if ValidEncoding(encoding) {
			var b bytes.Buffer
			gz := gzip.NewWriter(&b)

			if _, err := gz.Write([]byte(body)); err != nil {
				return err
			}

			if err := gz.Close(); err != nil {
				return err
			}

			body = b.String()

			response += "Content-Encoding: gzip\r\n"
		}
	}

	response += "Content-Type: text/plain\r\n"
	response += fmt.Sprintf("Content-Length: %d\r\n", len(body))
	response += fmt.Sprintf("\r\n%s", body)

	return c.write(response)
}

func (c HttpContext) Json(status int, jsonBody any) error {
	response := fmt.Sprintf("HTTP/1.1 %d %v\r\n", status, StatusText(status))

	body, err := json.Marshal(jsonBody)

	if err != nil {
		return err
	}

	encoding, isFound := c.req.Headers["Accept-Encoding"]

	if isFound {
		if ValidEncoding(encoding) {
			var b bytes.Buffer
			gz := gzip.NewWriter(&b)

			if _, err := gz.Write(body); err != nil {
				return err
			}

			if err := gz.Close(); err != nil {
				return err
			}

			body = b.Bytes()

			response += "Content-Encoding: gzip\r\n"
		}
	}

	response += "Content-Type: application/json\r\n"
	response += fmt.Sprintf("Content-Length: %d\r\n", len(body))
	response += fmt.Sprintf("\r\n%s", string(body))

	return c.write(response)
}

func (c HttpContext) File(status int, body string) error {
	response := fmt.Sprintf("HTTP/1.1 %d %v\r\n"+
		"Content-Type: application/octet-stream"+"\r\n"+
		"Content-Length: %d"+"\r\n\r\n"+
		"%s",
		status,
		StatusText(status),
		len(body),
		body)

	return c.write(response)
}

func (c HttpContext) Error(status int) error {
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", status, StatusText(status))
	return c.write(response)
}
