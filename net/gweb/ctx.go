package gweb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"time"
)

type (
	Ctx struct {
		ctx *gin.Context
	}

	ErrFormatter func(err error) string

	ParamSource string

	Param struct {
		Value *string
	}
)

const (
	Head     = ParamSource("Head")
	UriSlice = ParamSource("UriSlice")
	Query    = ParamSource("Query")
	Body     = ParamSource("Body")
)

func (p Param) String() (string, bool) {
	if p.Value == nil {
		return "", false
	}
	return *p.Value, true
}

func (p Param) I64() (int64, bool) {
	if p.Value == nil {
		return 0, false
	}
	i64, err := strconv.ParseInt(*p.Value, 10, 64)
	if err != nil {
		return 0, false
	}
	return i64, true
}

func (p Param) MustString(defValIfNotExist string) string {
	if p.Value == nil {
		return defValIfNotExist
	}
	return *p.Value
}

func (p Param) MustI64(defValIfNotExist int64) int64 {
	if p.Value == nil {
		return defValIfNotExist
	}
	i64, err := strconv.ParseInt(*p.Value, 10, 64)
	if err != nil {
		return defValIfNotExist
	}
	return i64
}

// TODO
func (c Ctx) GetParam(src ParamSource, key string, idx int) Param {
	switch src {
	case Head:
		values, exist := c.ctx.Request.Header[key]
		if !exist || idx >= len(values) {
			return Param{}
		}
		s := values[idx]
		return Param{Value: &s}
	case UriSlice:
	case Query:
	case Body:
	}

	return Param{}
}

func (c Ctx) PrintHead() {
	fmt.Println("--- header for", c.ctx.Request.URL, " ---")
	for k, v := range c.ctx.Request.Header {
		fmt.Println(k, v)
	}
}

// TODO Bind似乎可以自动从Header、Uri、Query、Body中解析需要的参数，但测试下来不对

// Unmarshal HTTP body to json
func (c Ctx) Body2JSON(output interface{}) error {
	return c.ctx.BindJSON(output)
}

// receive file from HTTP body
func (c Ctx) Body2File(key string) ([]byte, error) {
	f, hd, err := c.ctx.Request.FormFile(key)
	if err != nil {
		return nil, gerrors.Wrap(err, fmt.Sprintf("file[%s]", key))
	}
	if hd.Size == 0 {
		return nil, gerrors.Wrap(err, fmt.Sprintf("empty file[%s]", key))
	}
	buf := bytes.NewBuffer(nil)
	sz := int64(0)
	for sz < hd.Size {
		n, err := buf.ReadFrom(f)
		if err != nil {
			return nil, err
		}
		sz += n
	}
	return buf.Bytes(), nil
}

// http://example.com/path1/path2?k1=v1&k2=v2
// [k1,v1] and [k2, v2] are query params
func (c Ctx) GetQueryParamString(key string) string {
	return c.ctx.Query(key)
}

// http://example.com/path1/path2?k1=v1&k2=v2
// [k1,v1] and [k2, v2] are query params
func (c Ctx) GetQueryParamInt(key string) (int, error) {
	i64, err := strconv.ParseInt(c.ctx.Query(key), 10, 64)
	return int(i64), err
}

// http://example.com/path1/path2?k1=v1&k2=v2
// path1 and path2 are router slices
func (c Ctx) GetUriSliceString(name string) string {
	return c.ctx.Params.ByName(name)
}

// http://example.com/path1/path2?k1=v1&k2=v2
// path1 and path2 are router slices
func (c Ctx) GetRUriSliceInt(name string) (int, error) {
	i64, err := strconv.ParseInt(c.ctx.Params.ByName(name), 10, 64)
	return int(i64), err
}

func (c Ctx) WriteString(code int, format string, values ...interface{}) {
	c.ctx.String(code, format, values...)
}

func (c Ctx) WriteMapJSON(code int, values map[string]interface{}) {
	c.ctx.JSON(code, values)
}

func (c Ctx) WriteStructJSON(code int, output interface{}, errFmt ErrFormatter) {
	buf, err := json.Marshal(output)
	if err != nil {
		c.WriteString(code, errFmt(err))
		return
	}
	c.ctx.String(code, string(buf))
}

func (c Ctx) WriteError(code int, err error, errFmt ErrFormatter) {
	c.WriteString(code, errFmt(err))
}

func (c Ctx) ServeDiskFile(filename, filepath string) {
	c.ctx.Writer.Header().Set("content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	http.ServeFile(c.ctx.Writer, c.ctx.Request, filepath)
}

func (c Ctx) ServeBufferFile(filename string, r io.ReadSeeker) {
	c.ctx.Writer.Header().Set("content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	http.ServeContent(c.ctx.Writer, c.ctx.Request, ".xlsx", time.Now(), r)
}
