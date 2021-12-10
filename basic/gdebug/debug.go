package gdebug

import (
	"github.com/cryptowilliam/goutil/basic/glog"
	"github.com/cryptowilliam/goutil/net/gnet"
	"net/http"
	"net/http/pprof"
)

func ListenAndServe(listen string) error {
	r := http.NewServeMux()

	// basic stats
	r.HandleFunc("/debug/stats", statsHandler)

	// text profile
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// visual profile
	us, err := gnet.ParseUrl(listen)
	if err != nil {
		return err
	}
	us.Host.Port++
	vp := newVisualizePprof(us.String(), glog.DefaultLogger)
	r.HandleFunc("/debug/visual-pprof/"+ProfileCPU.String(), vp.serveVisualPprof)
	r.HandleFunc("/debug/visual-pprof/"+ProfileHeap.String(), vp.serveVisualPprof)
	r.HandleFunc("/debug/visual-pprof/"+ProfileBlock.String(), vp.serveVisualPprof)
	r.HandleFunc("/debug/visual-pprof/"+ProfileMutex.String(), vp.serveVisualPprof)
	r.HandleFunc("/debug/visual-pprof/"+ProfileGoRoutine.String(), vp.serveVisualPprof)
	r.HandleFunc("/debug/visual-pprof/"+ProfileThreadCreate.String(), vp.serveVisualPprof)

	// index page
	http.Handle("/", r)

	return http.ListenAndServe(listen, nil)
}