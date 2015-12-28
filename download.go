package main

import (
	"log"
	"net/http"
	"path"
	"regexp"
)

var validPath = regexp.MustCompile(`\A\/([a-zA-Z0-9]+)\/([a-zA-Z0-9\-\._]+)\z`)

func download(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// TODO check path length, depending on filename sanitation during upload
	if !validPath.MatchString(r.URL.Path) {
		log.Printf("Invalid path: %v", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	filePath := path.Join(DATA_DIR, r.URL.Path)
	log.Printf("Serving %s", filePath)
	http.ServeFile(w, r, filePath)
}
