package handler

import (
	"log"
	"net/http"
	"os"
	"strings"
)

type Api interface {
	Samples() []string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type samples struct {
	samples []string
}


func (this *samples) Samples(w http.ResponseWriter, r *http.Request) {
	return this.samples
}

func (this *samples) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

// read audio files in a directory and expose
// some api endpoints
func SamplesHandler(root string) http.Handler {
	supportedExtensions := []string{
		".wav", ".flac", ".aif", ".aiff",
	}

	fh, eo := os.Open(audioRoot)
	if eo != nil {
		return nil, eo
	}
	// determine if it is a directory
	info, es := fh.Stat()
	if es != nil {
		return nil, es
	}
	if !info.IsDir() {
		log.Fatal(audioRoot + " is not a directory")
	}
	fs, er := fh.Readdir(1024)
	if er != nil {
		return nil, er
	}
	var files []string
	for _, f := range fs {
		if isSupported(f, supportedExtensions) {
			files = append(files, f.Name())
		}
	}

	return &api{
		files,
		router,
	}, nil
}

func isSupported(f os.FileInfo, exts []string) bool {
	is := false
	for _, ext := range exts {
		if strings.HasSuffix(f.Name(), ext) {
			is = true
		}
	}
	return is
}
