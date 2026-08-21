package main

import (
	"embed"
	"flag"
	"log"
	"net/http"

	"nernst-cell/internal/web"
)

//go:embed web
var webFS embed.FS

//go:embed example/cu-conc.json
var cuConcJSON []byte

func main() {
	httpAddr := flag.String("http", ":8080", "serve the web console on this address (e.g. :8080)")
	flag.Parse()
	handler := web.NewServer(web.Assets{
		WebFS: webFS,
		Examples: map[string][]byte{
			"cu-conc": cuConcJSON,
		},
	})
	log.Printf("nernst-cell web console on http://localhost%s", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, handler); err != nil {
		log.Fatal(err)
	}
}
