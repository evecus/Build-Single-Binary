package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed all:public
var static embed.FS

func main() {
	port := flag.String("port", "", "port to listen on (default 3030, or PORT env)")
	flag.Parse()

	if *port == "" {
		*port = os.Getenv("PORT")
	}
	if *port == "" {
		*port = "3030"
	}

	dist, err := fs.Sub(static, "public")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(dist))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}
		_, statErr := fs.Stat(dist, path)
		if statErr != nil {
			// SPA fallback: return index.html for unknown routes
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})

	log.Printf("metacubexd running at http://0.0.0.0:%s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
