package lightning

import (
	"errors"
	"log"
	"net/http"
	"net/url"
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

// get the path component of a url
func GetPath(u *url.URL) (string, error) {
	// re := regexp.MustCompile("^http://[^/]*/([^/]*)\\??([\d%_]=[\d%_])*$")
	// matches := re.FindStringSubmatch(u.String())
	// if len(matches) < 2 {
	// 	msg := "could not get path component of url (" + u.String() + ")"
	// 	return "", errors.New(msg)
	// }
	// return matches[1], nil
	s := u.String()
	slidx := strings.LastIndex(s, "/")
	if slidx == -1 {
		return "", errors.New(s + " -- bad url")
	}
	rest := s[slidx:]
	qidx := strings.Index(s, "?")
	if qidx == -1 {
		return rest, nil
	} else {
		return s[slidx:qidx], nil
	}
}
