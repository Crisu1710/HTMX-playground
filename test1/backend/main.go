package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "GET" {
			fmt.Fprintf(w, "<p>GET</p>")
		}
		if r.Method == "POST" {
			fmt.Fprintf(w, "<p>POST</p>")
		}
	})
	http.ListenAndServe(":8080", nil)
}
