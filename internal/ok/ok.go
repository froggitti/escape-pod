package ok

import (
	"fmt"
	"net/http"
)

// Start starts the http server
func Start() {
	http.HandleFunc("/ok", okay)
	_ = http.ListenAndServe(":8080", nil)
}

func okay(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}
