package app

import (
	"fmt"
	"net"
	"time"
)

func newListener(cfg *ServerConfig) (net.Listener, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	// set actual port on config object (in case original port was 0)
	cfg.Port = ln.Addr().(*net.TCPAddr).Port
	return tcpListener{ln.(*net.TCPListener)}, nil
}

// custom TCP listener with keep-alive timeout
type tcpListener struct {
	*net.TCPListener
}

func (ln tcpListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
