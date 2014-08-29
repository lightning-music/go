package lightning

import (
	"net/http"
)

type Server interface {
	Listen(addr string) error
}

type serverImpl struct {
	patterns map[string]Pattern
	engine Engine
}

func apiHandler(writer http.ResponseWriter, req *http.Request) {
}

func (this *serverImpl) Listen(addr string) error {
	return http.ListenAndServe(addr, nil)
}

func NewServer(webroot string) Server {
	http.Handle("/", http.FileServer(http.Dir(webroot)))
	http.HandleFunc("/api", apiHandler)
	return &serverImpl{
		make(map[string]Pattern),
		NewEngine(),
	}
}
