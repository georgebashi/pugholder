
package main

import (
	"github.com/georgebashi/pugholder/image"
	"net/http"
	"path/filepath"
	"log"
	"sort"
	"hash/fnv"
	"code.google.com/p/gorilla/mux"
	"strconv"
	"fmt"
	"time"
)

type handler struct {
	files []string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var width, height int

	if size, ok := vars["size"]; ok {
		width, _ = strconv.Atoi(size)
		height = width
	} else {
		width, _ = strconv.Atoi(vars["width"])
		height, _ = strconv.Atoi(vars["height"])
	}

	hash := sum([]byte(fmt.Sprintf("%d/%d", width, height)))
	file := h.files[hash % uint32(len(h.files))]

	img := image.Open(file)
	defer img.Close()
	img.Strip()
	img.Resize(width, height)

	if vars["g"] != "" {
		img.Grayscale()
	}

	out := img.GetBytes()

	hours_in_month, _ := time.ParseDuration("730h")
	expire := time.Now().Add(hours_in_month)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Expires", expire.Format(time.RFC1123))
	w.Header().Set("ETag", fmt.Sprintf("%x", hash))
	w.Write(out)
}

func sum(input []byte) uint32 {
	hash := fnv.New32a()
	hash.Write(input)
	return hash.Sum32()
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

	r.Handle("/{g:[g/]*}{width:[1-9][0-9]*}/{height:[1-9][0-9]*}", h)
	r.Handle("/{g:[g/]*}{width:[1-9][0-9]*}x{height:[1-9][0-9]*}", h)
	r.Handle("/{g:[g/]*}{size:[1-9][0-9]*}", h)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	http.ListenAndServe(":9090", r)
}
