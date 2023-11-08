package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardmiddleware := alice.New(app.recoverpanic, app.logRequest, secureHeaders)
	dynamicmiidleware := alice.New(app.session.Enable)
	mux := pat.New()
	//mux := http.NewServeMux()
	mux.Get("/", dynamicmiidleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicmiidleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createsnippetForm))
	mux.Post("/snippet/create", dynamicmiidleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createsnippet))
	mux.Get("/snippet/:id", dynamicmiidleware.ThenFunc(app.showsnippet)) //
	mux.Get("/user/signup", dynamicmiidleware.ThenFunc(app.signupform))
	mux.Post("/user/signup", dynamicmiidleware.ThenFunc(app.signup))
	mux.Get("/user/login", dynamicmiidleware.ThenFunc(app.loginuserform))
	mux.Post("/user/login", dynamicmiidleware.ThenFunc(app.loginuser))
	mux.Post("/user/logout", dynamicmiidleware.Append(app.requireAuthenticatedUser).ThenFunc(app.logoutuser))
	// CHANGE START HERE
	mux.Get("/user/share", dynamicmiidleware.Append(app.requireAuthenticatedUser).ThenFunc(app.shareuserform))
	mux.Post("/user/share", dynamicmiidleware.Append(app.requireAuthenticatedUser).ThenFunc(app.shareuser))
	mux.Get("/user/sharedsnippets", dynamicmiidleware.Append(app.requireAuthenticatedUser).ThenFunc(app.sharedsnippets))
	mux.Get("/user/sharedwithyou", dynamicmiidleware.Append(app.requireAuthenticatedUser).ThenFunc(app.sharedwithyou))
	//TILL HERE
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardmiddleware.Then(mux)
}
