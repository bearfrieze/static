package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var nimbus string

func getChannel(url string) (*json.RawMessage, error) {

	r, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf(`Failed to get channel "%s": %s`, url, err)
	}
	decoder := json.NewDecoder(r.Body)
	var json *json.RawMessage
	decoder.Decode(&json)
	return json, nil
}

func handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}

	if r.Method != "POST" {
		http.Error(w, fmt.Sprintf("Unsupported method '%s'\n", r.Method), 501)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var urls []string
	err := decoder.Decode(&urls)
	if err != nil {
		log.Printf("Unable to decode request: %s\n", err)
		http.Error(w, err.Error(), 400)
		return
	}

	response := make(map[string]*json.RawMessage)
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string, response map[string]*json.RawMessage) {
			defer wg.Done()
			json, err := getChannel(fmt.Sprintf("%s/?url=%s", nimbus, url))
			if err != nil {
				log.Println(err)
				return
			}
			response[url] = json
		}(url, response)
	}
	wg.Wait()

	json, err := json.Marshal(&response)
	if err != nil {
		log.Printf("Unable to marshal response: %s\n", err)
	}

	w.Write(json)
}

func main() {

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	nimbus = os.Getenv("NIMBUS")
	log.Printf("Listening on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
