package main

import (
	"log"
	"net/http"
)

// Added some more functionalities as handler
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}
func showsnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet"))
}
func createsnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create a new snippet"))
}
func main() {
	mux := http.NewServeMux() //new serve mux used to map the url pattern ("/") to the handler home
	mux.HandleFunc("/", home)
	//"/snippet","/snippet/create" are fixed path patterns while "/" followed a subtree patterns(it's a catch-all)
	mux.HandleFunc("/snippet", showsnippet)
	mux.HandleFunc("/snippet/create", createsnippet)

	//listen and serve is used to start a web server with two parameter "TCP network address", servemux
	//log is used to log error if it returns
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

//Home(handler)->servemux->listenandserve->4000
