package main

import (
	"log"
	"net/http"
)

// Added some more functionalities as handler
func home(w http.ResponseWriter, r *http.Request) {
	//added a check to prevent "/" this to act as subtree pattern
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Snippetbox"))
}
func showsnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet"))
}
func createsnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST") //header() acts as a map
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
	}

	w.Write([]byte("create a new snippet"))
}
func main() {
	mux := http.NewServeMux()
	//we can also use http.HandleFunc without defining mux which use its own DefaultServemux
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showsnippet)
	mux.HandleFunc("/snippet/create", createsnippet)

	//listen and serve is used to start a web server with two parameter "TCP network address", servemux
	//log is used to log error if it returns
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

//Home(handler)->servemux->listenandserve->4000
