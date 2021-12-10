package gdebug

import (
	"net/http"
	"net/http/pprof"
	"time"
)

func ListenAndServe(listen string, duration time.Duration) error {
	// basic stats
	http.HandleFunc("/stats", statsHandler)

	// text profile
	http.HandleFunc("/pprof/", pprof.Index)
	http.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	http.HandleFunc("/pprof/profile", pprof.Profile)
	http.HandleFunc("/pprof/symbol", pprof.Symbol)
	http.HandleFunc("/pprof/trace", pprof.Trace)

	return http.ListenAndServe(listen, nil)
}
