package main

import (
	"log"
	"net/http"
)

// hame handler for executing home page logic
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}
func main() {
	mux := http.NewServeMux() //new serve mux used to map the url pattern ("/") to the handler home
	mux.HandleFunc("/", home)
	//listen and serve is used to start a web server with two parameter "TCP network address", servemux
	//log is used to log error if it returns
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

//Home(handler)->servemux->listenandserve->4000
