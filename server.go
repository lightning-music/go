package lightning

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Status string            `json:"string"`
}

type Server interface {
	Listen(addr string) error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
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

func (this *serverImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:      1024,
		WriteBufferSize:     1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		var res Response
		msgType, p, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
			return
		}
		note := new(Note)
		ed := json.Unmarshal(p, note)
		if ed != nil {
			log.Fatal(ed)
			return
		}
		ep := this.engine.PlayNote(*note)
		if ep != nil {
			log.Fatal(ep)
			return
		}
		res = Response{ "ok" }
		resb, em := json.Marshal(res)
		if em != nil {
			log.Fatal(em)
			return
		}
		ew := conn.WriteMessage(msgType, resb)
		if ew != nil {
			log.Fatal(ew)
		}
	}
}

func NewServer(webroot string) Server {
	server := &serverImpl{
		make(map[string]Pattern),
		NewEngine(),
	}
	http.Handle("/", http.FileServer(http.Dir(webroot)))
	http.Handle("/sample/play", server)
	http.HandleFunc("/api", apiHandler)
	return server
}
