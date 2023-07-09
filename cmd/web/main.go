package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/csn2002/Snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	errorlog      *log.Logger
	session       *sessions.Session
	infolog       *log.Logger
	snippets      *mysql.SnippetModel
	users         *mysql.Usermodel
	templateCatch map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn", "web1:123456789@/snippetbox?parseTime=true", "MySQL Database")
	secret := flag.String("secret", "mynameisanthonygolandservice@123", "secret key")
	flag.Parse()
	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorlog.Fatal(err)
	}
	defer db.Close()
	m := &mysql.SnippetModel{
		DB: db,
	}
	u := &mysql.Usermodel{
		DB: db,
	}
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorlog.Fatal(err)
	}
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	app := &application{
		errorlog:      errorlog,
		infolog:       infolog,
		session:       session,
		snippets:      m,
		templateCatch: templateCache,
		users:         u,
	}
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	infolog.Println("Starting server on %s", *addr)
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorlog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorlog.Fatal(err)
}
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
