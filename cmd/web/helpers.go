package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintln("%s\n%s", err.Error(), debug.Stack())
	app.errorlog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
func (app *application) adddefaultdata(r *http.Request, td *templateData) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.Currentyear = time.Now().Year()
	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCatch[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}
	buf := new(bytes.Buffer)
	err := ts.Execute(buf, app.adddefaultdata(r, td))
	if err != nil {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w)

}
