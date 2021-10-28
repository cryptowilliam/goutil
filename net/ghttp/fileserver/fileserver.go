package fileserver

// 来自 https://github.com/valyala/fasthttp/tree/master/examples/fileserver 2017-06-06 version
// 修改：把flag部分去掉，改成了FileServerOptions结构体传入所需参数

// Example static file server.
// Serves static files from the given directory.
// Exports various stats at /!stats .

import (
	"expvar"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/basic/glog"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/expvarhandler"
)

type FileServerOptions struct {
	Addr               string // Default "localhost:8080", TCP address to listen to
	AddrTLS            string // Default "", "TCP address to listen to TLS (aka SSL or HTTPS) requests. Leave empty for disabling TLS
	ByteRange          bool   // Default false, "Enables byte range requests if set to true
	CertFile           string // Default "./ssl-cert-snakeoil.pem", "Path to TLS certificate file
	Compress           bool   // Default false, "Enables transparent response compression if set to true
	Dir                string // Default "/usr/share/nginx/html", "Directory to serve static files from
	GenerateIndexPages bool   // Default true, "Whether to generate directory index pages
	KeyFile            string // Default "./ssl-cert-snakeoil.key", "Path to TLS key file
	Vhost              bool   // Default false, enables virtual hosting by prepending the requested path with the requested hostname
}

func Serve(fso *FileServerOptions, logger glog.Interface) {

	// Setup FS handler
	fs := &fasthttp.FS{
		Root:               fso.Dir,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: fso.GenerateIndexPages,
		Compress:           fso.Compress,
		AcceptByteRange:    fso.ByteRange,
	}
	if fso.Vhost {
		fs.PathRewrite = fasthttp.NewVHostPathRewriter(0)
	}
	fsHandler := fs.NewRequestHandler()

	// Create RequestHandler serving server stats on /!stats and files
	// on other requested paths.
	// /!stats output may be filtered using regexps. For example:
	//
	//   * /!stats?r=fs will show only stats (expvars) containing 'fs'
	//     in their names.
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/!stats":
			expvarhandler.ExpvarHandler(ctx)
		default:
			fsHandler(ctx)
			updateFSCounters(ctx, logger)
		}
	}

	// Start HTTP server.
	if len(fso.Addr) > 0 {
		logger.Infof("Starting HTTP server on " + fso.Addr)
		go func() {
			if err := fasthttp.ListenAndServe(fso.Addr, requestHandler); err != nil {
				logger.Erro(err)
			}
		}()
	}

	// Start HTTPS server.
	if len(fso.AddrTLS) > 0 {
		logger.Infof("Starting HTTPS server on " + fso.AddrTLS)
		go func() {
			if err := fasthttp.ListenAndServeTLS(fso.AddrTLS, fso.CertFile, fso.KeyFile, requestHandler); err != nil {
				logger.Erro(err)
			}
		}()
	}

	logger.Infof("Serving files from directory " + fso.Dir)
	logger.Infof("See stats at http://" + fso.Addr + "/!stats")

	// Wait forever.
	select {}
}

func updateFSCounters(ctx *fasthttp.RequestCtx, logger glog.Interface) {
	// Increment the number of fsHandler calls.
	fsCalls.Add(1)

	// Update other stats counters
	resp := &ctx.Response
	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		fsOKResponses.Add(1)
		fsResponseBodyBytes.Add(int64(resp.Header.ContentLength()))
		logger.Infof("File/Directory " + ctx.Request.URI().String() + " request succeed")
	case fasthttp.StatusNotModified:
		fsNotModifiedResponses.Add(1)
	case fasthttp.StatusNotFound:
		fsNotFoundResponses.Add(1)
		logger.Erro(gerrors.New("File/Directory " + ctx.Request.URI().String() + " request failed, path not found"))
	default:
		fsOtherResponses.Add(1)
	}
}

// Various counters - see https://golang.org/pkg/expvar/ for details.
var (
	// Counter for total number of fs calls
	fsCalls = expvar.NewInt("fsCalls")

	// Counters for various response status codes
	fsOKResponses          = expvar.NewInt("fsOKResponses")
	fsNotModifiedResponses = expvar.NewInt("fsNotModifiedResponses")
	fsNotFoundResponses    = expvar.NewInt("fsNotFoundResponses")
	fsOtherResponses       = expvar.NewInt("fsOtherResponses")

	// Total size in bytes for OK response bodies served.
	fsResponseBodyBytes = expvar.NewInt("fsResponseBodyBytes")
)
