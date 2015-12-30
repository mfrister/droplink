package main

import (
	"log"
	"net/http"

	"frister.net/go/droplink/servefs"

	"github.com/carbocation/interpose"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Address      string `envconfig:"ADDRESS"`
	DataDir      string `envconfig:"DATA_DIR"`
	PathPrefix   string `envconfig:"PATH_PREFIX"`
	URLPrefix    string `envconfig:"URL_PREFIX"`
	UploadSecret string `envconfig:"UPLOAD_SECRET"`
}

var config *Config

func loadConfig() *Config {
	conf := Config{}
	err := envconfig.Process("droplink", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	if conf.Address == "" {
		conf.Address = "localhost:8000"
	}
	if conf.DataDir == "" {
		conf.DataDir = "./data"
	}
	if conf.URLPrefix == "" {
		conf.URLPrefix = "http://localhost:8000/"
	}
	if conf.UploadSecret == "" {
		log.Printf("WARNING: No upload secret set, using default")
		conf.UploadSecret = "insecure"
	}
	return &conf
}

func main() {
	config = loadConfig()
	middle := interpose.New()

	router := http.NewServeMux()
	middle.UseHandler(router)

	router.Handle("/upload/", http.HandlerFunc(upload))
	router.Handle("/media/",
		http.StripPrefix("/media/",
			http.FileServer(servefs.JustFiles{http.Dir("media")})))
	router.Handle("/", http.HandlerFunc(download))

	log.Printf("Listening on %s", config.Address)
	err := http.ListenAndServe(config.Address,
		http.StripPrefix(config.PathPrefix, middle))
	if err != nil {
		log.Fatal(err)
	}
}
