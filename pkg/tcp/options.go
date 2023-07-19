package tcp

import "time"

type ServerOption func(server *Server)

func WithAddr(addr string) ServerOption {
	return func(server *Server) {
		server.addr = addr
	}
}

func WithHandler(handler ClientHandler) ServerOption {
	return func(server *Server) {
		server.handler = handler
	}
}

func WithTimeout(dur time.Duration) ServerOption {
	return func(server *Server) {
		if dur.Milliseconds() < ServerMinTimeoutMs {
			dur = time.Duration(ServerDefaultTimeoutMs) * time.Millisecond
		}
		server.timeout = dur
	}
}
