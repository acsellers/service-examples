package main

import (
	"log"
	"net/http"
	"os"

	"git.andrewcsellers.com/acsellers/web2015/store"
	"github.com/acsellers/dr/migrate"
	"github.com/acsellers/platform/router"
)

var Conn *store.Conn

func init() {
	// setup database
	var err error
	Conn, err = store.Open("postgres", "user=demo password=demo dbname=suggest")
	if err != nil {
		log.Fatal("Open Conn:", err)
	}
	m := migrate.Database{
		DB:         Conn.DB,
		Schema:     store.Schema,
		Translator: store.NewAppConfig("postgres"),
		Log:        log.New(os.Stdout, "Migration: ", 0),
	}
	err = m.Migrate()
	if err != nil {
		log.Fatal("Migrate Conn:", err)
	}
}

func main() {
	r := router.NewRouter()
	r.Many(NewSuggestCtrl())
	log.Fatal(http.ListenAndServe(":8080", r))
}
