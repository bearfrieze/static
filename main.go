package main

import (
	"encoding/json"
	"fmt"
	"github.com/coocood/freecache"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

const (
	cacheSize    = 800 * 1024 * 1024
	gcPercent    = 20
	logFrequency = 60 * 15
)

var nimbus string
var cache *freecache.Cache
var maxAgeRex = regexp.MustCompile(`max-age:(\d+)`)

func getChannel(url string) (*json.RawMessage, error) {

	key := []byte(url)

	var data []byte
	data, err := cache.Get(key)
	if err != nil {
		log.Printf("Fetching %s\n", url)
		r, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch %s: %s", url, err)
		}
		data, _ = ioutil.ReadAll(r.Body)
		maxAge, err := getMaxAge(&r.Header)
		if err == nil {
			log.Printf("Storing %s in cache, expires in %d seconds\n", url, *maxAge)
			err = cache.Set(key, data, *maxAge)
			if err != nil {
				log.Printf("Failed to store %s in cache: %s\n", err)
			}
		}
	}

	rm := json.RawMessage(data)
	return &rm, nil
}

func getMaxAge(h *http.Header) (*int, error) {

	cc := h.Get("Cache-Control")
	if cc == "" {
		return nil, fmt.Errorf("Cache-Control header not set")
	}

	matches := maxAgeRex.FindStringSubmatch(cc)
	if len(matches) == 0 {
		return nil, fmt.Errorf("max-age not specified")
	}

	maxAge, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, err
	}

	return &maxAge, nil
}

func handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

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

	cache = freecache.NewCache(cacheSize)
	debug.SetGCPercent(gcPercent)

	// Start logging cache stats
	go func() {
		for _ = range time.Tick(logFrequency * time.Second) {
			log.Printf("EntryCount: %d\n", cache.EntryCount())
			log.Printf("EvacuateCount: %d\n", cache.EvacuateCount())
			log.Printf("AverageAccessTime: %d\n", time.Now().Unix()-cache.AverageAccessTime())
			log.Printf("HitCount: %d\n", cache.HitCount())
			log.Printf("LookupCount: %d\n", cache.LookupCount())
			log.Printf("HitRate: %f\n\n", cache.HitRate())
		}
	}()

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	nimbus = os.Getenv("NIMBUS")
	log.Printf("Listening on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
