package ghttp

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// if you want to treat a tcp net.Conn as http, and give dest a valid http response,
// please use this co create a response.
// some Response member is required for browser to understand it.
func NewResponse(content []byte) *http.Response {
	t := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBuffer(content)),
		ContentLength: int64(len(content)),
		Header:        make(http.Header, 0),
	}
	return t
}
