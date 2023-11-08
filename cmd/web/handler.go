package main

import (
	"fmt"
	"github.com/csn2002/Snippetbox/pkg/forms"
	"github.com/csn2002/Snippetbox/pkg/models"
	"net/http"
	"strconv"
)

//all handler function is now a method against application
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if app.authenticateduser(r) != 0 {
		s, err := app.snippets.Latest(app.authenticateduser(r))
		if err != nil {
			app.serverError(w, err)
			return
		}
		data := &templateData{Snippets: s}
		app.render(w, r, "home.page.tmpl", data)
	} else {
		//app.session.Put(r, "flash", "Welcome to Snippetbox")
		app.session.Put(r, "flash", "Welcome to Snippetbox \n Login to see your Snippets")
		app.render(w, r, "login.page.tmpl", &templateData{
			Form: forms.New(nil),
		})
	}

}
func (app *application) showsnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.snippets.Get(id, app.authenticateduser(r))
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	//CHANGE STARTS HERE
	app.session.Put(r, "snippet_id", id)
	//CHANGE TILL HERE
	data := &templateData{
		Snippet: s,
	}
	app.render(w, r, "show.page.tmpl", data)
}
func (app *application) createsnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})

}
func (app *application) createsnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"), app.authenticateduser(r))
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Snippet Created Successfully")
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
func (app *application) signupform(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.Minlength("password", 10)
	form.MatchPattern("email", forms.EmailRX)
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) loginuserform(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) loginuser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredential {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "userID", id)
	//app.session.Put(r, "flash", "Login successful")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) logoutuser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userID")
	//CHANGES START HERE
	app.session.Remove(r, "snippet_id")
	//CHANGE TILL HERE
	app.session.Put(r, "flash", "You have been LoggedOut Successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// CHANGE STARTS HERE
func (app *application) shareuserform(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "share.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) shareuser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("email")
	form.MatchPattern("email", forms.EmailRX)
	if !form.Valid() {
		app.render(w, r, "share.page.tmpl", &templateData{
			Form: form,
		})
		return
	}
	id, err := app.share.Authenticate(form.Get("email"))
	if err == models.ErrInvalidCredential {
		form.Errors.Add("generic", "User not found")
		app.session.Put(r, "flash", "User not found")
		app.render(w, r, "share.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.share.Insert(app.authenticateduser(r), id, app.session.GetInt(r, "snippet_id"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Your Snippet Sharing is successful")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) sharedsnippets(w http.ResponseWriter, r *http.Request) {
	if app.authenticateduser(r) != 0 {
		s, err := app.share.Sharedsnippets(app.authenticateduser(r))
		if err != nil {
			app.serverError(w, err)
			return
		}
		data := &templateData{Snippets: s}
		app.render(w, r, "home.page.tmpl", data)
	} else {
		//app.session.Put(r, "flash", "Welcome to Snippetbox")
		app.session.Put(r, "flash", "Welcome to Snippetbox \n Login to see your Snippets")
		app.render(w, r, "login.page.tmpl", &templateData{
			Form: forms.New(nil),
		})
	}

}
func (app *application) sharedwithyou(w http.ResponseWriter, r *http.Request) {
	if app.authenticateduser(r) != 0 {
		s, err := app.share.LatestSharedWithYou(app.authenticateduser(r))
		if err != nil {
			app.serverError(w, err)
			return
		}
		data := &templateData{Snippets: s}
		app.render(w, r, "home.page.tmpl", data)
	} else {
		//app.session.Put(r, "flash", "Welcome to Snippetbox")
		app.session.Put(r, "flash", "Welcome to Snippetbox \n Login to see your Snippets")
		app.render(w, r, "login.page.tmpl", &templateData{
			Form: forms.New(nil),
		})
	}

}

// TILL HERE
