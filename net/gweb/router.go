package gweb

import (
	"github.com/cryptowilliam/goutil/net/ghttp"
	"github.com/gin-gonic/gin"
)

type (
	Router struct {
		ng *gin.Engine
	}

	HandlerFunc func(*Ctx)
)

func NewRouter() *Router {
	gin.SetMode(gin.ReleaseMode)
	return &Router{ng: gin.Default()}
}

func (r *Router) Handle(m ghttp.Method, relativePath string, fn HandlerFunc) {
	fn2 := func(c *gin.Context) {
		fn(&Ctx{ctx: c})
	}
	r.ng.Handle(string(m), relativePath, fn2)
}

func (r *Router) Static(relativePath, root string) {
	r.ng.Static(relativePath, root)
}

func (r *Router) StaticFile(relativePath, filepath string) {
	r.ng.StaticFile(relativePath, filepath)
}

func (r *Router) Serve(addr string) error {
	return r.ng.Run(addr) // listen and serve on 0.0.0.0:8080
}
