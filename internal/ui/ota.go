package ui

import (
	"log"
	"net/http"
)

// Deprecated
func (s *Server) serveOTA() {
	fs := http.FileServer(http.Dir(s.rootDir + "/ota"))
	http.Handle("/", fs)
	err := http.ListenAndServe(":"+s.otaPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Deprecated
func (s *Server) serveLOGS() {
	fs := http.FileServer(http.Dir(s.rootDir + "/logs"))
	http.Handle("/", fs)
	err := http.ListenAndServe(":"+s.otaPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
