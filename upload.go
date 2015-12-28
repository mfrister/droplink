package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"path"
	"regexp"
	"text/template"

	"github.com/tv42/base58"
)

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		uploadGet(w, r)
	} else if r.Method == "POST" {
		uploadPost(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func uploadGet(w http.ResponseWriter, r *http.Request) {
	var uploadTemplate = template.Must(template.ParseFiles("upload.html"))

	uploadTemplate.Execute(w, nil)
}

func uploadPost(w http.ResponseWriter, r *http.Request) {
	mr, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "Not a multipart request", http.StatusBadRequest)
		return
	}
	p, err := mr.NextPart()
	if err != nil {
		http.Error(w, fmt.Sprintf("NextPart failed: %v", err), http.StatusBadRequest)
		return
	}

	key, filename, err := processUpload(p, p.FileName())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to process upload: %v", err),
			http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s/%s", key, filename)
}

func processUpload(r io.ReadCloser, origName string) (string, string, error) {
	key, err := randomKey()
	if err != nil {
		return "", "", err
	}
	dirPath, err := createRandomDir(key)
	if err != nil {
		return "", "", err
	}
	sanitizedName := sanitizeFilename(origName)
	filePath := path.Join(dirPath, sanitizedName)

	if err := copyToFile(filePath, r); err != nil {
		return "", "", err
	}
	return key, sanitizedName, nil
}

func createRandomDir(key string) (string, error) {
	dir := path.Join(DATA_DIR, key)
	return dir, os.MkdirAll(dir, 0755)
}

func copyToFile(filePath string, r io.ReadCloser) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if n, err := io.Copy(file, r); err != nil {
		return err
	} else {
		log.Printf("Copied %d bytes", n)
	}
	if err := r.Close(); err != nil {
		return err
	}
	return file.Close()
}

// 6 bytes
var MAX_KEY = new(big.Int).Lsh(big.NewInt(1), (6 * 8))

func randomKey() (string, error) {
	randBig, err := rand.Int(rand.Reader, MAX_KEY)
	if err != nil {
		return "", err
	}
	dirName := base58.EncodeBig(make([]byte, 0), randBig)

	return string(dirName), nil
}

var FILENAME_REPLACE_REGEXP = regexp.MustCompile("[^a-zA-Z0-9\\-\\.]+")
var FILENAME_REPLACE_WITH = "_"

func sanitizeFilename(filename string) string {
	// TODO use first x chars of filename + ext to limit length
	return FILENAME_REPLACE_REGEXP.ReplaceAllString(
		filename, FILENAME_REPLACE_WITH)

}

// var uploadTemplate = template.Must(template.ParseFiles("upload.html"))
