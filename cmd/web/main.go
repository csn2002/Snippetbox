package main

import (
	"database/sql"
	"flag"
	"github.com/csn2002/Snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorlog      *log.Logger
	infolog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCatch map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn", "web1:123456789@/snippetbox?parseTime=true", "MySQL Database")
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
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorlog.Fatal(err)
	}
	app := &application{
		errorlog:      errorlog,
		infolog:       infolog,
		snippets:      m,
		templateCatch: templateCache,
	}

	infolog.Println("Starting server on %s", *addr)
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorlog,
		Handler:  app.routes(),
	}
	err = srv.ListenAndServe()
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
