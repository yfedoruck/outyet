package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

const (
	DbUser     = "postgres"
	DbPassword = "1"
	DbName     = "postgres"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func closeDb(db *sql.DB) {
	err := db.Close()
	check(err)
}

func dbName(db *sql.DB) string {
	rows, err := db.Query("SELECT datname FROM pg_database WHERE datistemplate = false limit 1;")
	check(err)
	if rows.Next() == false {
		log.Fatal("no database")
		return ""
	} else {
		var datname string
		err = rows.Scan(&datname)
		check(err)
		return datname
	}
}

func dbSelect() interface{} {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DbUser, DbPassword, DbName)
	db, err := sql.Open("postgres", dbInfo)
	check(err)
	defer closeDb(db)

	dbname := dbName(db)
	return struct {
		Test string
	}{dbname}
}

func main() {
	dt := dbSelect()

	httpAddr := flag.String("http", ":8080", "Listen address")
	flag.Parse()

	http.Handle("/", root(dt))
	log.Println("Starting web server on", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}

func root(data interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// tmpl is the HTML template that drives the user interface.
		var tmpl = template.Must(template.New("tmpl").Parse(`
<!DOCTYPE html><html><body><center>
	<h2>Hello, Go! {{.Test}} </h2>
</center></body></html>
`))
		tmpl.Execute(w, data)
	})
}
