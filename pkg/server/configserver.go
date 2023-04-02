package server

import (
	"embed"
	"net/http"
)

//go:embed resources/index.html
var content embed.FS

func Start() {
	fileServer := http.FileServer(http.FS(content))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	http.Handle("/", fileServer)
	http.ListenAndServe(":8090", nil)
}
