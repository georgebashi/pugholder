
package main

import (
	"net/http"
	"path/filepath"
	"log"
	"sort"
	"hash/fnv"
	"code.google.com/p/gorilla/mux"
)

type handler struct {
	files []string
}

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := h.handle(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h handler) handle(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	width := vars["width"]
	height := vars["height"]

	hash := fnv.New32a()
	hash.Write([]byte(width + "/" + height))
	log.Printf("w %s h %s # %x", width, height, hash.Sum32())
	return nil
}

func main() {
	files, err := filepath.Glob("img/*.jpg")
	if err != nil {
		log.Fatal(err)
		return
	}

	if len(files) == 0 {
		log.Fatal("no images found!")
		return
	}

	sort.Strings(files)

	r := mux.NewRouter()
	h := new(handler)
	h.files = files

	r.Handle("/{width:[1-9][0-9]*}/{height:[1-9][0-9]*}", h)
	http.ListenAndServe(":9090", r)
}
