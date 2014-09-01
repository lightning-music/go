package lightning

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Status string            `json:"string"`
	Message string           `json:"message"`
}

type Server interface {
	Listen(addr string) error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type serverImpl struct {
	audioRoot string
	engine Engine
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
		res = Response{ "ok", "played " + note.Sample, }
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

func NewServer(webRoot string, audioRoot string) (Server, error) {
	server := &serverImpl{
		audioRoot,
		NewEngine(),
	}
	api, ea := NewApi(audioRoot)
	if ea != nil {
		return nil, ea
	}
	http.Handle("/", http.FileServer(http.Dir(webRoot)))
	http.Handle("/sample/play", server)
	http.Handle("/api", api)
	return server, nil
}
