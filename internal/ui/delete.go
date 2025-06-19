package ui

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func (s *Server) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file, ok := vars["file"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := os.Remove(
		fmt.Sprintf("%s/%s/%s", s.rootDir, s.otaDir, file),
	); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
