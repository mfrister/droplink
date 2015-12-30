package servefs

import (
	"net/http"
	"os"
)

// from https://groups.google.com/d/msg/golang-nuts/bStLPdIVM6w/hidTJgDZpHcJ

type JustFiles struct {
	Fs http.FileSystem
}

func (fs JustFiles) Open(name string) (http.File, error) {
	f, err := fs.Fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

type neuteredReaddirFile struct {
	http.File
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}
