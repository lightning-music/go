package api

import (
	"encoding/json"
	"errors"
	// "fmt"
	// "github.com/lightning/lightning/api/handler"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Api interface {
	ListSamples() http.HandlerFunc
}

type api struct {
	Samps []string
}

func (this *api) Samples() []string {
	return this.Samps
}

func (this *api) ListSamples() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sl, me := json.Marshal(this.Samps)
		if me != nil {
			log.Println("could not marshal sample list: " + me.Error())
		}
		w.Write(sl)
	}
}

func NewApi(audioRoot string) (Api, error) {
	supportedExtensions := []string{
		".wav", ".flac", ".aif", ".aiff",
	}

	fh, eo := os.Open(audioRoot)
	if eo != nil {
		log.Println("could not open " + audioRoot + ": " + eo.Error())
		return nil, eo
	}
	// determine if it is a directory
	info, es := fh.Stat()
	if es != nil {
		log.Println("could not stat " + audioRoot + ": " + es.Error())
		return nil, es
	}
	if !info.IsDir() {
		log.Println(audioRoot + " is not a directory")
		return nil, errors.New(audioRoot + " is not a directory")
	}
	fs, er := fh.Readdir(1024)
	if er == io.EOF {
		log.Println("no samples in " + audioRoot)
		return nil, errors.New("no samples in " + audioRoot)
	}
	var files []string
	for _, f := range fs {
		if isSupported(f, supportedExtensions) {
			files = append(files, f.Name())
		}
	}

	return &api{
		files,
	}, nil
}

// determine if a file has a supported extension
func isSupported(f os.FileInfo, exts []string) bool {
	is := false
	for _, ext := range exts {
		if strings.HasSuffix(f.Name(), ext) {
			is = true
		}
	}
	return is
}
