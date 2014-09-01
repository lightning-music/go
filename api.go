package lightning

import (
	"net/http"
	"log"
	"os"
	"strings"
)

type Api interface {
	Samples() []string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type api struct {
	samples []string
}

func (this *api) Samples() []string {
	return this.samples
}

func (this *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
	}
}

func NewApi(audioRoot string) (Api, error) {
	supportedExtensions := []string{ ".wav", ".flac", ".opus" }

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
