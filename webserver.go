package gossh

import (
	"encoding/json"
	"github.com/bmizerany/pat"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func launchBrowser() {

	mux := pat.New()
	mux.Get("/", http.HandlerFunc(serveHome))
	mux.Get("/data", http.HandlerFunc(serveData))
	http.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(serveStatic)))

	http.Handle("/", mux)
	l := log.New(os.Stdout, "[gossh] ", 0)
	l.Printf("listening on %s", "4000")
	l.Fatal(http.ListenAndServe(":4000", nil))

}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	data, err := Asset("static/" + r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if filepath.Ext(r.URL.Path) == ".css" {
		w.Header().Set("Content-Type", "text/css")
	} else {
		w.Header().Set("Content-Type", "application/javascript")
	}
	w.Write(data)
}

//serve the homepage
func serveHome(w http.ResponseWriter, r *http.Request) {
	data, _ := Asset("static/main.html")
	w.Write(data)
}

//ajax request for JSON data representing SSH
func serveData(w http.ResponseWriter, r *http.Request) {
	jsonString, err := json.Marshal(applicationResults)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}
