package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"
)

var keyPath = regexp.MustCompile(`\A/([a-zA-Z0-9]+)\z`)
var fullPath = regexp.MustCompile(`\A\/([a-zA-Z0-9]+)\/([a-zA-Z0-9\-\._]+)\z`)

func download(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// TODO check path length, depending on filename sanitation during upload
	if !fullPath.MatchString(r.URL.Path) {
		if !keyPath.MatchString(r.URL.Path) {
			log.Printf("Invalid path: %v", r.URL.Path)
			http.NotFound(w, r)
			return
		} else {
			dirPath := path.Join(config.DataDir, r.URL.Path)
			files, err := ioutil.ReadDir(dirPath)
			if err != nil || len(files) != 1 {
				log.Printf("Error: Download: ReadDir failed (%v) or != 1 file found (%d)",
					err, len(files))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			key := strings.TrimPrefix(r.URL.Path, "/")
			target := fmt.Sprintf("%s%s/%s", config.URLPrefix, key, files[0].Name())

			log.Printf("Redirecting '%s' to '%s'", r.URL.Path, target)
			http.Redirect(w, r, target, http.StatusSeeOther)
			return
		}
	}

	filePath := path.Join(config.DataDir, r.URL.Path)
	log.Printf("Serving %s", filePath)
	http.ServeFile(w, r, filePath)
}
