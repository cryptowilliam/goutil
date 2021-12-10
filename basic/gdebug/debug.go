package gdebug

import (
	"net/http"
	"net/http/pprof"
	"time"
)

func ListenAndServe(listen string, duration time.Duration) error {
	r := http.NewServeMux()

	// basic stats
	r.HandleFunc("/debug/stats", statsHandler)

	// text profile
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// index page
	http.Handle("/", r)

	return http.ListenAndServe(listen, nil)
}
