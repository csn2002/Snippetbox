package main

import (
	"fmt"
	"github.com/csn2002/Snippetbox/pkg/models"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

//all handler function is now a method against application
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := &templateData{Snippets: s}
	app.render(w, r, "home.page.tmpl", data)

}
func (app *application) showsnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	data := &templateData{
		Snippet: s,
	}
	app.render(w, r, "show.page.tmpl", data)
}
func (app *application) createsnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}
func (app *application) createsnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")
	errors := make(map[string]string)
	if strings.TrimSpace(title) == "" {
		errors["title"] = "This Field Cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "is field is too long (maximum is 100 characters)"
	}
	if strings.TrimSpace(content) == "" {
		errors["content"] = "This Field Cannot be blank"
	}
	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This Field Cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "!" {
		errors["expires"] = "This field is invalid"
	}
	if len(errors) > 0 {
		app.render(w, r, "create.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
