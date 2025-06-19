package ui

import (
	"fmt"
	"net"
	"net/http"
)

func (s *Server) whoami(w http.ResponseWriter, r *http.Request) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer func() {
		_ = conn.Close()
	}()

	w.Header().Add("Content-Type", "application/json")

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Fprintf(
		w,
		`{"localaddr":"%s",
	"port":"%s",
	"wsPort":"%s"}`,
		localAddr.IP,
		s.port,
		s.websocketPort,
	)
}
