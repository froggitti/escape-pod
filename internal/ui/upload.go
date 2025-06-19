package ui

import (
	"fmt"
	"io"
	"net/http"
	"os"

	log "log"
)

func (s *Server) upload(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseMultipartForm(400 << 20)

	file, handler, err := r.FormFile("file")
	fileName := r.FormValue("file_name")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	f, err := os.OpenFile(
		fmt.Sprintf("%s/%s/%s", s.rootDir, s.otaDir, handler.Filename),
		os.O_WRONLY|os.O_CREATE,
		0600,
	)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer f.Close()

	if _, err = io.WriteString(w, "File "+fileName+" Uploaded successfully"); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err = io.Copy(f, file); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
