package ui

import (
	"crypto/tls"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/DDLbots/escape-pod/internal/ui/logwebsocket"
	"github.com/DDLbots/escape-pod/internal/ui/statuswebsocket"
	"github.com/DDLbots/escape-pod/internal/version"
	"github.com/digital-dream-labs/vector-bluetooth/ble"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	timeout = 600
)

type ServeHTTPFunc func(rw http.ResponseWriter, r *http.Request)

func (f ServeHTTPFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	f(rw, r)
}

func Middle(in http.Handler) http.Handler {
	return ServeHTTPFunc(func(rw http.ResponseWriter, r *http.Request) {

	})
}

// Server stores the config
type Server struct {
	rootDir string
	otaDir  string
	logsDir string
	uiDir   string

	port string

	websocketPort string

	Router *mux.Router
	crt    tls.Certificate

	dlCh   chan ble.StatusChannel
	dlWs   *statuswebsocket.StatusStreamer
	parsed chan interface{}

	// DEPRECATED
	otaPort string
}

// New returns a new server
func New(opts ...Option) (*Server, error) {

	r := mux.NewRouter()

	s := &Server{
		rootDir: "/usr/lib/escape-pod",
		otaDir:  "ota",
		uiDir:   "dist",
		logsDir: "logs",
		Router:  r,
		parsed:  make(chan interface{}),
	}

	for _, opt := range opts {
		if err := opt.applyOption(s); err != nil {
			return nil, err
		}
	}
	// TODO fix this
	// r.Use(loggingMiddleware)
	r.HandleFunc("/api/v1/ota/whoami", s.whoami).Methods(http.MethodGet)

	r.HandleFunc("/api/v1/ota", s.listOTA).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/ota", s.upload).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/ota/{file}", s.delete).Methods(http.MethodDelete)

	r.HandleFunc("/api/v1/ota/{filename}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filename, ok := vars["filename"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		http.ServeFile(w, r, filepath.Join(s.rootDir, s.otaDir, filename))
	}).Methods(http.MethodGet)

	r.HandleFunc("/api/v1/logs/{filename}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filename, ok := vars["filename"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		http.ServeFile(w, r, filepath.Join(s.rootDir, s.logsDir, filename))
	}).Methods(http.MethodGet)

	r.HandleFunc("/api/v1/logs", s.listLOGS).Methods(http.MethodGet)

	r.HandleFunc("/v1/downloadstatus", s.dlWs.DownloadStatus)

	// note: these are not logs. this is the intent coming out of the robot.
	r.HandleFunc("/v1/logs", logwebsocket.ParsedIntent(s.parsed))

	r.HandleFunc("/api/v1/version", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		w := json.NewEncoder(rw)
		if err := w.Encode(version.DefaultVersionResponse); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(http.MethodGet)

	f, err := os.Open(filepath.Join(s.rootDir, s.uiDir, "index.html"))
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	indexTemplate, err := template.New("indexTemplate").Parse(string(b))
	if err != nil {
		return nil, err
	}

	spa := spaHandler{
		staticPath:    filepath.Join(s.rootDir, s.uiDir),
		indexPath:     "index.html",
		indexTemplate: indexTemplate,
		Version:       template.JS(version.Version),
	}

	r.Handle("/", spa)
	r.Handle("/{(?!api|v1)}", spa)

	return s, nil
}

func (s *Server) Handle(pattern string, handler http.Handler) {
	s.Router.PathPrefix("/api").Handler(handler)
}

func (s *Server) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	corsHandler := handlers.CORS(
		handlers.AllowedMethods(
			[]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"},
		),
		handlers.AllowedOrigins(
			[]string{"*"},
		),
	)(s.Router)

	corsHandler.ServeHTTP(wr, r)
}

// TODO: turn this in to a logging middleware that works with our structured logger interface
func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

// Serve serves the static content
func (s *Server) Serve() {
	srv := &http.Server{
		Addr:         ":" + s.port,
		ReadTimeout:  timeout * time.Second,
		WriteTimeout: timeout * time.Second,
		Handler:      s,
	}

	log.Fatal(
		srv.ListenAndServe(
		// tlsListener,
		),
	)
}
