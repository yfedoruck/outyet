package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
)

func main() {
	httpAddr   := flag.String("http", ":8080", "Listen address")
	flag.Parse()

	http.HandleFunc("/", root)
	log.Println("Starting web server on", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}

// tmpl is the HTML template that drives the user interface.
var tmpl = template.Must(template.New("tmpl").Parse(`
<!DOCTYPE html><html><body><center>
	<h2>Hello, Go! </h2>
</center></body></html>
`))

func root(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}