package gdebug

import (
	"github.com/cryptowilliam/goutil/basic/glog"
	"html/template"
	"net/http"
	"net/http/pprof"
)

var (
	pathToVisualPprofCPU          = "/debug/visual-pprof/" + profileCPU.String()
	pathToVisualPprofHeap         = "/debug/visual-pprof/" + profileHeap.String()
	pathToVisualPprofBlock        = "/debug/visual-pprof/" + profileBlock.String()
	pathToVisualPprofMutex        = "/debug/visual-pprof/" + profileMutex.String()
	pathToVisualPprofAllocs       = "/debug/visual-pprof/" + profileAllocs.String()
	pathToVisualPprofGoRoutine    = "/debug/visual-pprof/" + profileGoRoutine.String()
	pathToVisualPprofThreadCreate = "/debug/visual-pprof/" + profileThreadCreate.String()
)

func serveIndexPage(w http.ResponseWriter, r *http.Request) {
	var v = struct {
		TextPprofIndex          string
		VisualPprofCPU          string
		VisualPprofHeap         string
		VisualPprofBlock        string
		VisualPprofMutex        string
		VisualPprofAllocs       string
		VisualPprofGoRoutine    string
		VisualPprofThreadCreate string
	}{
		TextPprofIndex:          "/debug/pprof/",
		VisualPprofCPU:          pathToVisualPprofCPU,
		VisualPprofHeap:         pathToVisualPprofHeap,
		VisualPprofBlock:        pathToVisualPprofBlock,
		VisualPprofMutex:        pathToVisualPprofMutex,
		VisualPprofAllocs:       pathToVisualPprofAllocs,
		VisualPprofGoRoutine:    pathToVisualPprofGoRoutine,
		VisualPprofThreadCreate: pathToVisualPprofThreadCreate,
	}
	if err := indexPageTmpl.Execute(w, &v); err != nil {
		glog.Erro(err, "tmpl.Execute")
	}
}

func ListenAndServe(listen string) error {
	r := http.NewServeMux()

	// index page
	r.HandleFunc("/", serveIndexPage)

	// basic stats
	r.HandleFunc("/debug/stats", statsHandler)

	// text profile
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// visual profile
	vp, err := newVisualizePprof(glog.DefaultLogger)
	if err != nil {
		return err
	}
	r.HandleFunc(pathToVisualPprofCPU, vp.serveVisualPprof)
	r.HandleFunc(pathToVisualPprofHeap, vp.serveVisualPprof)
	r.HandleFunc(pathToVisualPprofBlock, vp.serveVisualPprof)
	r.HandleFunc(pathToVisualPprofMutex, vp.serveVisualPprof)
	r.HandleFunc(pathToVisualPprofAllocs, vp.serveVisualPprof)
	r.HandleFunc(pathToVisualPprofGoRoutine, vp.serveVisualPprof)
	r.HandleFunc(pathToVisualPprofThreadCreate, vp.serveVisualPprof)

	// index page
	http.Handle("/", r)

	return http.ListenAndServe(listen, nil)
}

var indexPageTmpl = template.Must(template.New("").Parse(`<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title id="">Debug System</title>
  </head>
  <body>
	<div style="width:95%; height:95%; display: inline-block; vertical-align: top;">
		<p><a href="{{.TextPprofIndex}}" target="_blank">{{.TextPprofIndex}}</a></p>
		<p><a href="{{.VisualPprofCPU}}" target="_blank">{{.VisualPprofCPU}}</a></p>
		<p><a href="{{.VisualPprofHeap}}" target="_blank">{{.VisualPprofHeap}}</a></p>
		<p><a href="{{.VisualPprofBlock}}" target="_blank">{{.VisualPprofBlock}}</a></p>
		<p><a href="{{.VisualPprofMutex}}" target="_blank">{{.VisualPprofMutex}}</a></p>
		<p><a href="{{.VisualPprofAllocs}}" target="_blank">{{.VisualPprofAllocs}}</a></p>
		<p><a href="{{.VisualPprofGoRoutine}}" target="_blank">{{.VisualPprofGoRoutine}}</a></p>
		<p><a href="{{.VisualPprofThreadCreate}}" target="_blank">{{.VisualPprofThreadCreate}}</a></p>
	</div>
  </body>
</html>
`))
