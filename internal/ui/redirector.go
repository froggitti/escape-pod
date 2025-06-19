package ui

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath    string
	indexPath     string
	indexTemplate *template.Template
	Version       template.JS
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	path = filepath.Join(h.staticPath, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {

		if h.indexTemplate != nil {

			if err := h.indexTemplate.ExecuteTemplate(w, "indexTemplate", h); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if (filepath.Base(r.URL.Path) == "index.html" || filepath.Base(r.URL.Path) == "/") && h.indexTemplate != nil {
		if err := h.indexTemplate.Execute(w, h); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
