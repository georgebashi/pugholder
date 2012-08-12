
package main

import (
	"regexp"
	"net/http"
	"path/filepath"
	"log"
	"sort"
	"hash/fnv"
	"code.google.com/p/gorilla/mux"
)

var request_regex = regexp.MustCompile("/([0-9]+)/([0-9]+)")

func handler(w http.ResponseWriter, r *http.Request) {
	dims := request_regex.FindStringSubmatch(r.URL.RequestURI())
	if dims == nil {
		return
	}
	width := dims[1]
	height := dims[2]
	log.Printf("w %s h %s", width, height)

	hash := fnv.New32a()
	hash.Write([]byte(width + "/" + height))
	log.Printf("%x", hash.Sum32())
}

func main() {
	image_paths, err := filepath.Glob("img/*.jpg")
	if err != nil {
		log.Fatal(err)
		return
	}

	if len(image_paths) == 0 {
		log.Fatal("no images found!")
		return
	}

	sort.Strings(image_paths)

	r := mux.NewRouter()
	r.HandleFunc("/{width:[1-9][0-9]*}/{height:[1-9][0-9]*}", handler)
	http.ListenAndServe(":9090", r)
}
