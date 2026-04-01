package web

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path"
)

//go:embed all:dist
var frontendFS embed.FS

func SPAHandler() http.Handler {
	dist, _ := fs.Sub(frontendFS, "dist")
	fileServer := http.FileServer(http.FS(dist))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := path.Clean(r.URL.Path)
		if p == "/" {
			p = "index.html"
		} else {
			p = p[1:]
		}

		if _, err := fs.Stat(dist, p); os.IsNotExist(err) {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
}
