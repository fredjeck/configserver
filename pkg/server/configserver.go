package server

import (
	"embed"
	b64 "encoding/base64"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/fredjeck/configserver/pkg/config"
	"github.com/fredjeck/configserver/pkg/encrypt"
)

type ConfigServer struct {
	Configuration config.Config
	Key           *[32]byte
}

//go:embed resources
var content embed.FS

func (server ConfigServer) encryptValue(w http.ResponseWriter, req *http.Request) {
	value, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ciphered, error := encrypt.Encrypt(value, server.Key)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	base := b64.StdEncoding.EncodeToString(ciphered[:])
	w.Write([]byte(base))

}

func (server ConfigServer) Start() {

	router := http.NewServeMux()
	middleware := gitRepoMiddleWare()
	routerWithMiddleware := middleware(router)

	serverRoot, err := fs.Sub(content, "resources")
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/api/encrypt", server.encryptValue)
	router.Handle("/", http.FileServer(http.FS(serverRoot)))

	err = http.ListenAndServe(":8090", routerWithMiddleware)
	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
	}
}

func gitRepoMiddleWare() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.RequestURI, "/git") {
				// handle

				return
			}

			// call next handler
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
