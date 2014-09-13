package api

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/hypebeast/go-osc/osc"
	"github.com/lightning/go/binding"
	"github.com/lightning/go/seq"
	"github.com/lightning/go/types"
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
	engine types.Engine
	oscServer *osc.OscServer
	router *mux.Router
}

func (this *serverImpl) Listen(addr string) error {
	return http.ListenAndServe(addr, this.router)
	// return http.ListenAndServe(addr, nil)
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

		note, enp := seq.ParseNote(p)
		if enp != nil {
			log.Fatal(enp)
			return
		}

		// log.Printf("playing sample %v number=%v velocity=%v\n",
		// 	note.Sample(), note.Number(), note.Velocity())

		ep := this.engine.PlayNote(note)
		if ep != nil {
			log.Fatal(ep)
			return
		}
		res = Response{ "ok", "played " + note.Sample(), }
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
	oscPort := 4800
	rtr := mux.NewRouter()
	srv := &serverImpl{
		audioRoot,
		binding.NewEngine(),
		osc.NewOscServer("127.0.0.1", oscPort),
		rtr,
	}
	// api handler
	api, ea := NewApi(audioRoot)
	if ea != nil {
		return nil, ea
	}
	// osc comm
	pm := func(msg *osc.OscMessage) {
		osc.PrintOscMessage(msg)
	}
	srv.oscServer.AddMsgHandler("/sample/play", pm)
	go srv.oscServer.ListenAndDispatch();
	log.Printf("osc server listening on port %d\n", oscPort)
	// setup handlers under default ServeMux
	fh := http.FileServer(http.Dir(webRoot))
	srv.router.Handle("/", fh)
	srv.router.Handle("/{*.(js|css|png|jpg)}", fh)
	srv.router.Handle("/sample/play", srv)
	srv.router.Handle("/api", api)
	return srv, nil
}
