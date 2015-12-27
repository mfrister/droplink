package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"fmt"

	"github.com/carbocation/interpose"
)

func main() {
	middle := interpose.New()

	router := http.NewServeMux()
	middle.UseHandler(router)

	router.Handle("/upload", http.HandlerFunc(upload))

	log.Println("Listening....")
	err := http.ListenAndServe("127.0.0.1:8000", middle)
	if err != nil {
		log.Fatal(err)
	}
}

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
	body, err := ioutil.ReadAll(p)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read part: %v", err),
			http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Length: %d", len(body))
}

var uploadTemplate = template.Must(template.New("uploadForm").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>DropLink</title>
    <link rel="stylesheet" type="text/css" href="./media/dropzone.css">
    <script type="text/javascript" src="./media/dropzone.js"></script>
</head>
<body>
    Dropzone removed during rebase (invalid depdency,
        too lazy to fix old commit)
</body>
</html>
`))
