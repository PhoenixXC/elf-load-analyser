package web

import (
    "embed"
    "github.com/benbjohnson/hashfs"
    "net/http"
    "net/url"
    "strings"
)

//go:embed template/dist
var reactDist embed.FS
var fsys = hashfs.NewFS(reactDist)
var apiPrefix = "/api"

func FrontedService() http.Handler {
    return autoAppendPrefix("/template/dist", hashfs.FileServer(fsys))
}

func autoAppendPrefix(prefix string, h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request){
        if strings.HasPrefix(req.URL.Path, apiPrefix) {
            h.ServeHTTP(w, req)
            return
        }

        var index string
        if "/" == req.URL.Path {
            index = "/index.html"
        }

        newR := new(http.Request)
        *newR = *req
        newR.URL = new(url.URL)
        *newR.URL = *req.URL
        newR.URL.Path = prefix + req.URL.Path + index
        newR.URL.RawPath = prefix + req.URL.RawPath + index
        h.ServeHTTP(w, newR)
    })
}
