package main

import (
	"log"
	"net"
	"slices"
	"strings"
	"time"
)

type Handler func(c HttpContext) error

type HttpServer struct {
	Routes map[string]Handler
}

func NewHttpServer() HttpServer {
	return HttpServer{
		Routes: make(map[string]Handler),
	}
}

func (s HttpServer) Listen(addr string) {
	log.Printf("Starting HTTP server on %s\n", addr)
	l, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v\n", err)
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *HttpServer) GET(path string, h Handler) {
	pathWithMethod := "GET" + path
	s.Routes[pathWithMethod] = h
}

func (s *HttpServer) POST(path string, h Handler) {
	pathWithMethod := "POST" + path
	s.Routes[pathWithMethod] = h
}

func (s *HttpServer) handleConnection(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
	httpReq := NewRequest(conn)

	ctx := HttpContext{
		req:  httpReq,
		conn: conn,
	}

	for path, handler := range s.Routes {
		pathWithMethod := httpReq.Method + httpReq.Uri

		// exact path
		if path == pathWithMethod {
			handler(ctx)
			break
		}

		// any path
		if path[len(path)-1] == '*' {
			if strings.HasPrefix(pathWithMethod, path[0:len(path)-2]) {
				handler(ctx)
				break
			}
		}
	}

	ctx.Error(404)
}

func StatusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Not Found"
	}
}

func ValidEncoding(encoding []string) bool {
	return slices.Contains(encoding, "gzip")
}
