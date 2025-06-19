package debug

import (
	"expvar"
	"net/http"
	_ "net/http/pprof" // Register the pprof handlers
	"runtime"
)

type DebugServer struct {
	*http.Server
}

// NewDebugServer provides new debug http server
func NewDebugServer(address string) *DebugServer {
	return &DebugServer{
		&http.Server{
			Addr:    address,
			Handler: http.DefaultServeMux,
		},
	}
}

var m = struct {
	gr  *expvar.Int
	req *expvar.Int
}{
	gr:  expvar.NewInt("goroutines"),
	req: expvar.NewInt("requests"),
}

// Metrics updates program counters.
func Metrics(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		// Increment the request counter.
		m.req.Add(1)

		// Update the count for the number of active goroutines every 100 requests.
		if m.req.Value()%100 == 0 {
			m.gr.Set(int64(runtime.NumGoroutine()))
		}
	}

	return http.HandlerFunc(fn)
}
