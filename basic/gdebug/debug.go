package gdebug

import (
	"github.com/cryptowilliam/goutil/basic/glog"
	"html/template"
	"net/http"
	"net/http/pprof"
)

var (
	pathToBasicStats              = "/debug/stats"
	pathToTextPprofIndex          = "/debug/pprof/"
	pathToTextPprofCmdline        = "/debug/pprof/cmdline"
	pathToTextPprofProfile        = "/debug/pprof/profile"
	pathToTextPprofSymbol         = "/debug/pprof/symbol"
	pathToTextPprofTrace          = "/debug/pprof/trace"
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
		BasicStats              string
		TextPprofIndex          string
		TextPprofCmdline        string
		TextPprofProfile        string
		TextPprofSymbol         string
		TextPprofTrace          string
		VisualPprofCPU          string
		VisualPprofHeap         string
		VisualPprofBlock        string
		VisualPprofMutex        string
		VisualPprofAllocs       string
		VisualPprofGoRoutine    string
		VisualPprofThreadCreate string
	}{
		BasicStats:              pathToBasicStats,
		TextPprofIndex:          pathToTextPprofIndex,
		TextPprofCmdline:        pathToTextPprofCmdline,
		TextPprofProfile:        pathToTextPprofProfile,
		TextPprofSymbol:         pathToTextPprofSymbol,
		TextPprofTrace:          pathToTextPprofTrace,
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
	r.HandleFunc(pathToBasicStats, statsHandler)

	// text profile
	r.HandleFunc(pathToTextPprofIndex, pprof.Index)
	r.HandleFunc(pathToTextPprofCmdline, pprof.Cmdline)
	r.HandleFunc(pathToTextPprofProfile, pprof.Profile)
	r.HandleFunc(pathToTextPprofSymbol, pprof.Symbol)
	r.HandleFunc(pathToTextPprofTrace, pprof.Trace)

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
		<p><a href="{{.BasicStats}}" target="_blank">Basic stats information -> {{.BasicStats}}</a></p>
		<p><a href="{{.TextPprofIndex}}" target="_blank">Text pprof index -> {{.TextPprofIndex}}</a></p>
		<p><a href="{{.TextPprofCmdline}}" target="_blank">Text pprof cmdline -> {{.TextPprofCmdline}}</a></p>
		<p><a href="{{.TextPprofProfile}}" target="_blank">Text pprof profile -> {{.TextPprofProfile}}</a></p>
		<p><a href="{{.TextPprofSymbol}}" target="_blank">Text pprof symbol -> {{.TextPprofSymbol}}</a></p>
		<p><a href="{{.TextPprofTrace}}" target="_blank">Text pprof trace -> {{.TextPprofTrace}}</a></p>
		<p><a href="{{.VisualPprofCPU}}" target="_blank">Visual pprof CPU (wait 10+ seconds) -> {{.VisualPprofCPU}}</a></p>
		<p><a href="{{.VisualPprofHeap}}" target="_blank">Visual pprof heap -> {{.VisualPprofHeap}}</a></p>
		<p><a href="{{.VisualPprofBlock}}" target="_blank">Visual pprof block -> {{.VisualPprofBlock}}</a></p>
		<p><a href="{{.VisualPprofMutex}}" target="_blank">Visual pprof mutex -> {{.VisualPprofMutex}}</a></p>
		<p><a href="{{.VisualPprofAllocs}}" target="_blank">Visual pprof allocs -> {{.VisualPprofAllocs}}</a></p>
		<p><a href="{{.VisualPprofGoRoutine}}" target="_blank">Visual pprof Go routine -> {{.VisualPprofGoRoutine}}</a></p>
		<p><a href="{{.VisualPprofThreadCreate}}" target="_blank">Visual pprof thread create -> {{.VisualPprofThreadCreate}}</a></p>
	</div>
  </body>
</html>
`))
