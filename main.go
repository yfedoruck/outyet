package main

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"runtime"

	// "flag"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	db *sql.DB
)

var baseDir string

func basePath() string {
	if baseDir != "" {
		return baseDir
	}
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("Caller error")
	}

	baseDir = filepath.Dir(b)
	return baseDir
}

type dbConf struct {
	User     string `json:"User"`
	Password string `json:"Password"`
	Name     string `json:"Name"`
	Host     string `json:"Host"`
	Port     string `json:"Port"`
}

func pgConfig() dbConf {
	file, err := os.Open(basePath() + filepath.FromSlash("/config/"+env()+"/postgres.json"))
	check(err)

	dbConf := dbConf{}
	err = json.NewDecoder(file).Decode(&dbConf)
	check(err)

	return dbConf
}

func env() string {
	domain := os.Getenv("USERDOMAIN")
	if domain == "home" {
		return "local"
	}
	return "heroku"
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func closeDb() {
	err := db.Close()
	check(err)
}

func dbName() string {
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

func init() {
	dbConf := pgConfig()
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.Name)
	fmt.Println(dbInfo)
	var err error
	db, err = sql.Open("postgres", dbInfo)
	check(err)

	for i, connected := 0, false; connected == false && i < 4; i++ {
		err = db.Ping()
		if err == nil {
			connected = true
			return
		} else {
			log.Println("Error: Could not establish a connection with the database!", err, " but I still tried to connect...")
			time.Sleep(2 * time.Second)
		}
	}
	panic(err)
}

func main() {
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
