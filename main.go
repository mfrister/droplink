package main

import (
	"log"
	"net/http"

	"github.com/carbocation/interpose"
)

const DATA_DIR = "./data"

func main() {
	middle := interpose.New()

	router := http.NewServeMux()
	middle.UseHandler(router)

	router.Handle("/upload", http.HandlerFunc(upload))
	router.Handle("/", http.HandlerFunc(download))

	log.Println("Listening....")
	err := http.ListenAndServe("127.0.0.1:8000", middle)
	if err != nil {
		log.Fatal(err)
	}
}
