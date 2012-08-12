
package main

import (
	"net/http"
	"path/filepath"
	"log"
	"sort"
	"hash/fnv"
	"code.google.com/p/gorilla/mux"
	"strconv"
)

type handler struct {
	files []string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	width, _ := strconv.Atoi(vars["width"])
	height, _ := strconv.Atoi(vars["height"])

	hash := fnv.New32a()
	hash.Write([]byte(r.RequestURI))

	file := h.files[hash.Sum32() % uint32(len(h.files))]
	w.Write(resize(file, width, height))
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
