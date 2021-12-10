package gdebug

import (
	"net/http"
	"time"
)

func ListenAndServe(listen string, duration time.Duration) error {
	r := http.NewServeMux()

	// basic stats
	r.HandleFunc("/stats", statsHandler)

	// index page
	http.Handle("/", r)

	return http.ListenAndServe(listen, nil)
}
