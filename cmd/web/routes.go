package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardmiddleware := alice.New(app.recoverpanic, app.logRequest, secureHeaders)
	mux := pat.New()
	//mux := http.NewServeMux()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.createsnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.createsnippet))
	mux.Get("/snippet/:id", http.HandlerFunc(app.showsnippet)) //
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardmiddleware.Then(mux)
}
