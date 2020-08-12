package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/rmouram/snippetbox/pkg/models"
	"github.com/rmouram/snippetbox/pkg/models/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
)

type application struct {
	infoLog *log.Logger
	errorLog *log.Logger
	snippets *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func openDB(dns string) (*sql.DB, error){
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dns := flag.String("dns", "root:@tcp(127.0.0.1:3308)/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile )

	db, err := openDB(*dns)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("../../ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}


	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
