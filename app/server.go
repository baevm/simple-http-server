package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	"strings"
)

const ADDR = "0.0.0.0:3000"

var StaticDir string

func main() {
	flag.StringVar(&StaticDir, "directory", "", "directory for static files")
	flag.Parse()

	server := NewHttpServer()

	server.GET("/user-agent", func(c HttpContext) error {
		if userAgentHeader, isOk := c.req.Headers["User-Agent"]; isOk {
			return c.String(200, userAgentHeader[0])
		} else {
			return c.Error(404)
		}
	})

	server.GET("/", func(c HttpContext) error {
		return c.String(200, "")
	})

	server.GET("/echo/*", func(c HttpContext) error {
		path := strings.Split(c.req.Uri, "/")
		echoStr := path[len(path)-1]
		return c.String(200, echoStr)
	})

	server.GET("/get-user", func(c HttpContext) error {
		return c.Json(200, struct {
			User     string
			Password string
		}{
			User:     "test123",
			Password: "helloworld",
		})
	})

	server.GET("/files/*", func(c HttpContext) error {
		path := strings.Split(c.req.Uri, "/")
		filename := path[len(path)-1]

		filepath := fmt.Sprintf("%s/%s", StaticDir, filename)

		content, err := os.ReadFile(filepath)

		if err != nil {
			return c.Error(404)
		} else {
			return c.File(200, string(content))
		}
	})

	server.POST("/files/*", func(c HttpContext) error {
		path := strings.Split(c.req.Uri, "/")
		filename := path[len(path)-1]

		filepath := fmt.Sprintf("%s/%s", StaticDir, filename)

		err := os.WriteFile(filepath, c.req.Body, 0644)
		if err != nil {
			return c.Error(400)
		} else {
			return c.String(201, "")
		}
	})

	server.Listen(ADDR)
}
