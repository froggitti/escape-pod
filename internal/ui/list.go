package ui

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type listResp struct {
	Files []string `json:"files,omitempty"`
}

func (s *Server) listOTA(w http.ResponseWriter, r *http.Request) {
	list, err := ioutil.ReadDir(filepath.Join(s.rootDir, s.otaDir))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	f := listResp{}
	for _, t := range list {
		f.Files = append(f.Files, t.Name())
	}

	resp, err := json.Marshal(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) listLOGS(w http.ResponseWriter, r *http.Request) {
	list, err := ioutil.ReadDir(filepath.Join(s.rootDir, s.logsDir))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	f := listResp{}
	for _, t := range list {
		f.Files = append(f.Files, t.Name())
	}

	resp, err := json.Marshal(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
