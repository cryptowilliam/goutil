package gdebug

import (
	"github.com/cryptowilliam/goutil/basic/glog"
	"html/template"
	"net/http"
	"net/http/pprof"
)

type (
	allPaths struct {
		BasicStats              string
		TextPprofIndex          string
		TextPprofAllocs         string
		TextPprofBlock          string
		TextPprofCmdline        string
		TextPprofGoroutine      string
		TextPprofHeap           string
		TextPprofMutex          string
		TextPprofProfile        string
		TextPprofThreadCreate   string
		TextPprofTrace          string
		TextPprofSymbol         string
		FullGoroutineStackDump  string
		VisualPprofCPU          string
		VisualPprofHeap         string
		VisualPprofBlock        string
		VisualPprofMutex        string
		VisualPprofAllocs       string
		VisualPprofGoroutine    string
		VisualPprofThreadCreate string
	}
)

var (
	aps = allPaths{
		BasicStats:              "/debug/stats",
		TextPprofIndex:          "/debug/pprof/",
		TextPprofAllocs:         "/debug/pprof/allocs?debug=1",
		TextPprofBlock:          "/debug/pprof/block?debug=1",
		TextPprofCmdline:        "/debug/pprof/cmdline?debug=1",
		TextPprofGoroutine:      "/debug/pprof/goroutine?debug=1",
		TextPprofHeap:           "/debug/pprof/heap?debug=1",
		TextPprofMutex:          "/debug/pprof/mutex?debug=1",
		TextPprofProfile:        "/debug/pprof/profile?debug=1",
		TextPprofThreadCreate:   "/debug/pprof/threadcreate?debug=1",
		TextPprofTrace:          "/debug/pprof/trace?debug=1",
		TextPprofSymbol:         "/debug/pprof/symbol?debug=1",
		FullGoroutineStackDump:  "/debug/pprof/goroutine?debug=2",
		VisualPprofCPU:          "/debug/visual-pprof/" + profileCPU.String(),
		VisualPprofHeap:         "/debug/visual-pprof/" + profileHeap.String(),
		VisualPprofBlock:        "/debug/visual-pprof/" + profileBlock.String(),
		VisualPprofMutex:        "/debug/visual-pprof/" + profileMutex.String(),
		VisualPprofAllocs:       "/debug/visual-pprof/" + profileAllocs.String(),
		VisualPprofGoroutine:    "/debug/visual-pprof/" + profileGoRoutine.String(),
		VisualPprofThreadCreate: "/debug/visual-pprof/" + profileThreadCreate.String(),
	}
)

func serveIndexPage(w http.ResponseWriter, r *http.Request) {
	if err := indexPageTmpl.Execute(w, &aps); err != nil {
		glog.Erro(err, "tmpl.Execute")
	}
}

func ListenAndServe(listen string) error {
	r := http.NewServeMux()

	// index page
	r.HandleFunc("/", serveIndexPage)

	// basic stats
	r.HandleFunc(aps.BasicStats, statsHandler)

	// text profile
	r.HandleFunc(aps.TextPprofIndex, pprof.Index)
	r.HandleFunc(aps.TextPprofCmdline, pprof.Cmdline)
	r.HandleFunc(aps.TextPprofProfile, pprof.Profile)
	r.HandleFunc(aps.TextPprofSymbol, pprof.Symbol)
	r.HandleFunc(aps.TextPprofTrace, pprof.Trace)

	// visual profile
	vp, err := newVisualizePprof(glog.DefaultLogger)
	if err != nil {
		return err
	}
	r.HandleFunc(aps.VisualPprofCPU, vp.serveVisualPprof)
	r.HandleFunc(aps.VisualPprofHeap, vp.serveVisualPprof)
	r.HandleFunc(aps.VisualPprofBlock, vp.serveVisualPprof)
	r.HandleFunc(aps.VisualPprofMutex, vp.serveVisualPprof)
	r.HandleFunc(aps.VisualPprofAllocs, vp.serveVisualPprof)
	r.HandleFunc(aps.VisualPprofGoroutine, vp.serveVisualPprof)
	r.HandleFunc(aps.VisualPprofThreadCreate, vp.serveVisualPprof)

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
		<p><a href="{{.BasicStats}}" target="_blank">Basic stats information</a></p>
		<p><a href="{{.TextPprofIndex}}" target="_blank">Text pprof index</a></p>
		<p><a href="{{.TextPprofAllocs}}" target="_blank">Text pprof allocs</a></p>
		<p><a href="{{.TextPprofBlock}}" target="_blank">Text pprof block</a></p>
		<p><a href="{{.TextPprofCmdline}}" target="_blank">Text pprof cmdline</a></p>
		<p><a href="{{.TextPprofGoroutine}}" target="_blank">Text pprof goroutine</a></p>
		<p><a href="{{.TextPprofHeap}}" target="_blank">Text pprof heap</a></p>
		<p><a href="{{.TextPprofMutex}}" target="_blank">Text pprof mutex</a></p>
		<p><a href="{{.TextPprofProfile}}" target="_blank">Text pprof profile</a></p>
		<p><a href="{{.TextPprofThreadCreate}}" target="_blank">Text pprof threadcreate</a></p>
		<p><a href="{{.TextPprofTrace}}" target="_blank">Text pprof trace</a></p>
		<p><a href="{{.TextPprofSymbol}}" target="_blank">Text pprof symbol</a></p>
		<p><a href="{{.VisualPprofCPU}}" target="_blank">Visual pprof CPU (wait 10+ seconds)</a></p>
		<p><a href="{{.VisualPprofHeap}}" target="_blank">Visual pprof heap</a></p>
		<p><a href="{{.VisualPprofBlock}}" target="_blank">Visual pprof block</a></p>
		<p><a href="{{.VisualPprofMutex}}" target="_blank">Visual pprof mutex</a></p>
		<p><a href="{{.VisualPprofAllocs}}" target="_blank">Visual pprof allocs</a></p>
		<p><a href="{{.VisualPprofGoRoutine}}" target="_blank">Visual pprof Go routine</a></p>
		<p><a href="{{.VisualPprofThreadCreate}}" target="_blank">Visual pprof thread create</a></p>
	</div>
  </body>
</html>
`))
