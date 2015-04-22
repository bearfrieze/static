package main

import (
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	"log"
	"net/http"
	"os"
	// "time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "GET" {
		http.Error(w, fmt.Sprintf("Unsupported method '%s'\n", r.Method), 501)
		return
	}

	fmt.Fprintln(w, "Hello, World!")
}

func main() {

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	log.Printf("Listening on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
