package main

import (
	"os"

	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

func main() {
	connectDb()
	defer closeDb()
	dt := struct {
		Test string
	}{dbName()}

	// port := flag.String("http", ":8080", "Listen address")
	// flag.Parse()
	port := os.Getenv("PORT")
	if port == "" {
		// it's docker inner port (second port in from docker-compose: ports)
		port = "5000"
	}

	http.Handle("/", root(dt))
	log.Println("Starting web server on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func root(data interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// tmpl is the HTML template that drives the user interface.
		var tmpl = template.Must(template.New("tmpl").Parse(`
<!DOCTYPE html><html><body><center>
	<h2>Hello, {{.Test}}, from docker! </h2>
</center></body></html>
`))
		tmpl.Execute(w, data)
	})
}
