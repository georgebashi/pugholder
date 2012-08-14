
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
	"bytes"
)

type handler struct {
	files []string
	start_time time.Time
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

	etag := fmt.Sprintf("\"%x\"", sum([]byte(fmt.Sprintf("%d/%d/%s/%d", width, height, file, h.start_time.Unix()))))
	none_match := r.Header.Get("If-None-Match")
	if none_match == etag || fmt.Sprintf("\"%s\"", none_match) == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	img := image.Open(file)
	defer img.Close()
	img.Strip()
	img.Resize(width, height)

	if vars["g"] != "" {
		img.Grayscale()
	}

	out := img.GetBytes()

	hours_in_month, _ := time.ParseDuration("730h")
	expire := time.Now().Add(hours_in_month).UTC()
	w.Header().Set("Expires", expire.Format(http.TimeFormat))
	w.Header().Set("ETag", etag)
	http.ServeContent(w, r, "img.jpg", h.start_time, bytes.NewReader(out))
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
	h.start_time = time.Now()

	r.Handle("/{g:[g/]*}{width:[1-9][0-9]*}/{height:[1-9][0-9]*}", h)
	r.Handle("/{g:[g/]*}{width:[1-9][0-9]*}x{height:[1-9][0-9]*}", h)
	r.Handle("/{g:[g/]*}{size:[1-9][0-9]*}", h)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	http.ListenAndServe(":9090", r)
}
