package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	db *sql.DB
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

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

func env() string {
	domain := os.Getenv("USERDOMAIN")
	if domain == "home" {
		return "local"
	}
	return "heroku"
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


func connectDb() {
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
