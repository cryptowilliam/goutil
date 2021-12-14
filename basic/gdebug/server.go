package gdebug

import (
	"github.com/cryptowilliam/goutil/basic/glog"
	"github.com/cryptowilliam/goutil/container/gternary"
	"github.com/cryptowilliam/goutil/net/gnet"
	"html/template"
	"net/http"
	"net/http/pprof"
	"runtime"
)

type (
	allPaths struct {
		BasicStats                        string
		TextIndex                    string
		TextAllocs                   string
		TextBlock                    string
		TextCmdline                  string
		TextGoroutine                string
		TextHeap                     string
		TextMutex                    string
		TextProfile                  string
		TextThreadCreate             string
		TextTrace                    string
		TextSymbol                   string
		TextFullGoroutineStackDump   string
		VisualProfile                string
		VisualHeap                   string
		VisualBlock                  string
		VisualMutex                  string
		VisualAllocs                 string
		VisualGoroutine              string
		VisualThreadCreate           string
	}
)

var (
	aps = allPaths{
		BasicStats:                        "/debug/stats",
		TextIndex:                    "/debug/pprof/",
		TextAllocs:                   "/debug/pprof/allocs?debug=1",
		TextBlock:                    "/debug/pprof/block?debug=1",
		TextCmdline:                  "/debug/pprof/cmdline?debug=1",
		TextGoroutine:                "/debug/pprof/goroutine?debug=1",
		TextHeap:                     "/debug/pprof/heap?debug=1",
		TextMutex:                    "/debug/pprof/mutex?debug=1",
		TextProfile:                  "/debug/pprof/profile?debug=1",
		TextThreadCreate:             "/debug/pprof/threadcreate?debug=1",
		TextTrace:                    "/debug/pprof/trace?debug=1",
		TextSymbol:                   "/debug/pprof/symbol?debug=1",
		TextFullGoroutineStackDump:   "/debug/pprof/goroutine?debug=2",
		VisualProfile:                "/debug/visual-pprof/profile",
		VisualHeap:                   "/debug/visual-pprof/heap",
		VisualBlock:                  "/debug/visual-pprof/block",
		VisualMutex:                  "/debug/visual-pprof/mutex",
		VisualAllocs:                 "/debug/visual-pprof/allocs",
		VisualGoroutine:              "/debug/visual-pprof/goroutine",
		VisualThreadCreate:           "/debug/visual-pprof/threadcreate",
	}
)

func serveIndexPage(w http.ResponseWriter, r *http.Request) {
	if err := indexPageTmpl.Execute(w, &aps); err != nil {
		glog.Erro(err, "tmpl.Execute")
	}
}

// control part profiles with runtime interfaces.
// CPU profile is not included but only started when user want, because
// CPU profile requires an opened file stream to store the data.
func enableProfile(enable bool) {
	// SetBlockProfileRate controls the fraction of goroutine blocking events
	// that are reported in the blocking profile. The profiler aims to sample
	// an average of one blocking event per rate nanoseconds spent blocked.
	//
	// To include every blocking event in the profile, pass rate = 1.
	// To turn off profiling entirely, pass rate <= 0.
	blockRate := gternary.If(enable).Int(1, 0)
	runtime.SetBlockProfileRate(blockRate)

	mutexFraction := gternary.If(enable).Int(1, 0)
	runtime.SetMutexProfileFraction(mutexFraction)
}

// ListenAndServe starts a debug server with web ui.
// Visit http://listen to see it.
// Note: don't start it if not necessary.
func ListenAndServe(listen string) error {
	us, err := gnet.ParseUrl(listen)
	if err != nil {
		return err
	}

	r := http.NewServeMux()
	enableProfile(true)
	defer enableProfile(false)

	// index page
	r.HandleFunc("/", serveIndexPage)

	// basic stats
	r.HandleFunc(aps.BasicStats, statsHandler)

	// text profile
	r.HandleFunc(aps.TextIndex, pprof.Index)
	r.HandleFunc(aps.TextCmdline, pprof.Cmdline)
	r.HandleFunc(aps.TextProfile, pprof.Profile)
	r.HandleFunc(aps.TextSymbol, pprof.Symbol)
	r.HandleFunc(aps.TextTrace, pprof.Trace)

	// visual profile
	// Note: Graphviz required
	vp, err := newSvgServer(us.Host.Port, glog.DefaultLogger)
	if err != nil {
		return err
	}
	r.HandleFunc(aps.VisualProfile, vp.serve)
	r.HandleFunc(aps.VisualHeap, vp.serve)
	r.HandleFunc(aps.VisualBlock, vp.serve)
	r.HandleFunc(aps.VisualMutex, vp.serve)
	r.HandleFunc(aps.VisualAllocs, vp.serve)
	r.HandleFunc(aps.VisualGoroutine, vp.serve)
	r.HandleFunc(aps.VisualThreadCreate, vp.serve)

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
	<div style="width:45%; height:95%; display: inline-block; vertical-align: top;">
		<p><a href="{{.TextIndex}}" target="_blank">Text index</a></p>
		<p><a href="{{.TextSymbol}}" target="_blank">Text symbol</a></p>
		<p><a href="{{.TextCmdline}}" target="_blank">Text cmdline</a></p>
		<p><a href="{{.TextTrace}}" target="_blank">Text trace</a></p>

		<p><a href="{{.TextProfile}}" target="_blank">Text profile</a></p>
		<p><a href="{{.TextHeap}}" target="_blank">Text heap</a></p>
		<p><a href="{{.TextBlock}}" target="_blank">Text block</a></p>
		<p><a href="{{.TextMutex}}" target="_blank">Text mutex</a></p>
		<p><a href="{{.TextAllocs}}" target="_blank">Text allocs</a></p>
		<p><a href="{{.TextGoroutine}}" target="_blank">Text goroutine</a></p>
		<p><a href="{{.TextThreadCreate}}" target="_blank">Text threadcreate</a></p>
		<p><a href="{{.TextFullGoroutineStackDump}}" target="_blank">Text full goroutine stack dump</a></p>
	</div>
	<div style="width:45%; height:95%; display: inline-block; vertical-align: top;">
		<p><a href="{{.VisualProfile}}" target="_blank">Visual profile</a></p>
		<p><a href="{{.VisualHeap}}" target="_blank">Visual heap</a></p>
		<p><a href="{{.VisualBlock}}" target="_blank">Visual block</a></p>
		<p><a href="{{.VisualMutex}}" target="_blank">Visual mutex</a></p>
		<p><a href="{{.VisualAllocs}}" target="_blank">Visual allocs</a></p>
		<p><a href="{{.VisualGoroutine}}" target="_blank">Visual goroutine</a></p>
		<p><a href="{{.VisualThreadCreate}}" target="_blank">Visual threadcreate</a></p>

		<p><a href="{{.BasicStats}}" target="_blank">Basic stats information</a></p>
	</div>
	<br>
	<p>
	Profile Descriptions:
	<ul>
	<li><div class=profile-name>allocs: </div> A sampling of all past memory allocations</li>
	<li><div class=profile-name>block: </div> Stack traces that led to blocking on synchronization primitives</li>
	<li><div class=profile-name>cmdline: </div> The command line invocation of the current program</li>
	<li><div class=profile-name>goroutine: </div> Stack traces of all current goroutines</li>
	<li><div class=profile-name>heap: </div> A sampling of memory allocations of live objects. You can specify the gc GET parameter to run GC before taking the heap sample.</li>
	<li><div class=profile-name>mutex: </div> Stack traces of holders of contended mutexes</li>
	<li><div class=profile-name>profile: </div> CPU profile. You can specify the duration in the seconds GET parameter. After you get the profile file, use the go tool pprof command to investigate the profile.</li>
	<li><div class=profile-name>threadcreate: </div> Stack traces that led to the creation of new OS threads</li>
	<li><div class=profile-name>trace: </div> A trace of execution of the current program. You can specify the duration in the seconds GET parameter. After you get the trace file, use the go tool trace command to investigate the trace.</li>
	</ul>
	</p>
  </body>
</html>
`))

var mainPageTmpl = template.Must(template.New("").Parse(`<html>
<head>
<title>/debug/pprof/</title>
<style>
.profile-name{
	display:inline-block;
	width:6rem;
}
</style>
</head>
<body>
/debug/pprof/<br>
<br>
Types of profiles available:
<table>
<thead><td>Count</td><td>Profile</td></thead>
<tr><td>86</td><td><a href='allocs?debug=1'>allocs</a></td></tr>
<tr><td>19</td><td><a href='block?debug=1'>block</a></td></tr>
<tr><td>0</td><td><a href='cmdline?debug=1'>cmdline</a></td></tr>
<tr><td>114</td><td><a href='goroutine?debug=1'>goroutine</a></td></tr>
<tr><td>86</td><td><a href='heap?debug=1'>heap</a></td></tr>
<tr><td>0</td><td><a href='mutex?debug=1'>mutex</a></td></tr>
<tr><td>0</td><td><a href='profile?debug=1'>profile</a></td></tr>
<tr><td>13</td><td><a href='threadcreate?debug=1'>threadcreate</a></td></tr>
<tr><td>0</td><td><a href='trace?debug=1'>trace</a></td></tr>
</table>
<a href="goroutine?debug=2">full goroutine stack dump</a>
<br>
<p>
Profile Descriptions:
<ul>
<li><div class=profile-name>allocs: </div> A sampling of all past memory allocations</li>
<li><div class=profile-name>block: </div> Stack traces that led to blocking on synchronization primitives</li>
<li><div class=profile-name>cmdline: </div> The command line invocation of the current program</li>
<li><div class=profile-name>goroutine: </div> Stack traces of all current goroutines</li>
<li><div class=profile-name>heap: </div> A sampling of memory allocations of live objects. You can specify the gc GET parameter to run GC before taking the heap sample.</li>
<li><div class=profile-name>mutex: </div> Stack traces of holders of contended mutexes</li>
<li><div class=profile-name>profile: </div> CPU profile. You can specify the duration in the seconds GET parameter. After you get the profile file, use the go tool pprof command to investigate the profile.</li>
<li><div class=profile-name>threadcreate: </div> Stack traces that led to the creation of new OS threads</li>
<li><div class=profile-name>trace: </div> A trace of execution of the current program. You can specify the duration in the seconds GET parameter. After you get the trace file, use the go tool trace command to investigate the trace.</li>
</ul>
</p>
<p>Visualized pprof requires Graphviz.</p>
</body>
</html>`))
