package lightning

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/hypebeast/go-osc/osc"
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
	oscServer *osc.OscServer
	router *mux.Router
}

func (this *serverImpl) Listen(addr string) error {
	return http.ListenAndServe(addr, this.router)
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
	srv := &serverImpl{
		audioRoot,
		NewEngine(),
		osc.NewOscServer("127.0.0.1", 4800),
		mux.NewRouter(),
	}
	// api handler
	api, ea := NewApi(audioRoot)
	if ea != nil {
		return nil, ea
	}
	// osc comm
	srv.oscServer.AddMsgHandler("/sample/play", func(msg *osc.OscMessage) {
		osc.PrintOscMessage(msg)
	})
	go srv.oscServer.ListenAndDispatch();
	// setup handlers under default ServeMux
	srv.router.Handle("/", http.FileServer(http.Dir(webRoot)))
	srv.router.Handle("/sample/play", srv)
	srv.router.Handle("/api", api)
	return srv, nil
}
