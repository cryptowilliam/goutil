package gdebug

import (
	"net/http"
	"net/http/pprof"
	"time"
)

func ListenAndServe(listen string, duration time.Duration) error {
	r := http.NewServeMux()

	// basic stats
	r.HandleFunc("/stats", statsHandler)

	// text profile
	r.HandleFunc("/pprof/", pprof.Index)
	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/pprof/trace", pprof.Trace)

	// index page
	r.Handle("/", r)

	return http.ListenAndServe(listen, r)
}
