package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	_ "github.com/rmouram/snippetbox/pkg/models"
	"github.com/rmouram/snippetbox/pkg/models/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type contextKey string

var contextKeyUser = contextKey("user")

type application struct {
	infoLog *log.Logger
	errorLog *log.Logger
	session *sessions.Session
	snippets *mysql.SnippetModel
	templateCache map[string]*template.Template
	user *mysql.UserModel
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
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
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

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		session: session,
		snippets: &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		user: &mysql.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
		TLSConfig: tlsConfig,
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("../../tls/cert.pem", "../../tls/key.pem")
	errorLog.Fatal(err)
}
