package gperf

import (
	"github.com/cryptowilliam/goutil/container/gnum"
	stats_api "github.com/fukata/golang-stats-api-handler"
	"net/http"
)

func Serve(port uint16) error {
	http.HandleFunc("/api/stats", stats_api.Handler)
	return http.ListenAndServe(":"+gnum.FormatUint16(port), nil)
}
