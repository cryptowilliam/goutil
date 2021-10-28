package gfileserver

// reference:
// https://github.com/chinglinwen/fileserver2 current used version
// https://github.com/jpillora/uploader code simple and beautiful

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/glog"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	path string
)

func NewFileServer(localAddr, dir string) {
	path = dir
	http.HandleFunc("/", detector)
	err := http.ListenAndServe(localAddr, nil)
	if err != nil {
		glog.Erro(err)
		os.Exit(1)
	}
}

func detector(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.RequestURI, "uploadapi") {
		uploadHandler(w, r)
		return
	}
	// print logs
	ip := strings.Split(r.RemoteAddr, ":")[0]
	log.Println(ip, r.RequestURI, "visited")

	if strings.HasSuffix(r.RequestURI, "upload") {
		uploadPageHandler(w, r)
		return
	}
	http.FileServer(http.Dir(path)).ServeHTTP(w, r)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
