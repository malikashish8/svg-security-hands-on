package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/img/", image)

	log.Print("Listening on localhost:8000")
	http.ListenAndServe("localhost:8000", mux)
}

func image(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %s", r.URL.Path)
	// get file name from url
	fileName := r.URL.Path[5:]
	log.Println(fileName)
	// read filename
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Error reading file")
	} else {
		w.Header().Set("Content-Type", "image/svg+xml")
		// if query sets sandbox, set sandbox
		if r.URL.Query().Get("sandbox") == "true" {
			w.Header().Set("Content-Security-Policy", "sandbox")
		}
		if r.URL.Query().Get("csp-src-none") == "true" {
			w.Header().Set("Content-Security-Policy", "script-src 'none'")
		}
		if r.URL.Query().Get("content-disposition") == "true" {
			w.Header().Set("Content-Disposition", `attachment;filename="`+fileName+`"`)
		}
		w.Write(fileBytes)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fileBytes, err := ioutil.ReadFile("index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Error reading file")
	} else {
		w.Header().Set("Content-Type", "text/html")
		w.Write(fileBytes)
	}
}
