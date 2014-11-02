package lightning

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hypebeast/go-osc/osc"
	"log"
	"net/http"
)

const (
	PATTERN_LENGTH = 4096
	PATTERN_DIV    = "1/4"
	OSC_ADDR       = "127.0.0.1"
	OSC_PORT       = 4800
)

// function that handles websocket messages
type WebsocketHandler func(conn *websocket.Conn, messageType int, msg []byte)

type PatternEdit struct {
	Pos  Pos  `json:"pos"`
	Note Note `json:"note"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Server interface {
	Connect(ch1 string, ch2 string) error
	Listen(addr string) error
}

type simp struct {
	audioRoot string
	engine    Engine
	oscServer *osc.OscServer
	sequencer *Sequencer
}

func (this *simp) Listen(addr string) error {
	return http.ListenAndServe(addr, nil)
}

func (this *simp) AddTo(pos Pos, note Note) error {
	return this.sequencer.AddTo(pos, note)
}

// generate the MetroFunc that wires the metro to
// the pattern and the audio engine
func genMetroFunc(s *simp) MetroFunc {
	return func(pos Pos) {
		notes := s.sequencer.NotesAt(pos % PATTERN_LENGTH)
		for _, note := range notes {
			s.engine.PlayNote(note)
		}
	}
}

// upgrade an http handler to a websocket handler.
// that is probably not the best way to describe what is
// happening here.
func (s *simp) upgrade(handler WebsocketHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// upgrade http connection
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("could not upgrade http conn to ws: " + err.Error())
			return
		}
		// get messages and call handler
		for {
			msgType, bs, err := conn.ReadMessage()

			if err != nil {
				log.Println("could not read ws message: " + err.Error())
				continue
			}

			handler(conn, msgType, bs)
		}
	}
}

func (this *simp) playSample() http.HandlerFunc {
	return this.upgrade(func(conn *websocket.Conn, msgType int, msg []byte) {
		var res Response
		note, enp := ParseNote(msg)
		if enp != nil {
			fmtstr := "could not parse note from %s: %s\n"
			log.Printf(fmtstr, bytes.NewBuffer(msg).String(), enp.Error())
			return
		}

		// log.Printf("playing sample %v number=%v velocity=%v\n",
		// 	note.Sample(), note.Number(), note.Velocity())

		ep := this.engine.PlayNote(*note)
		if ep != nil {
			log.Println("could not play note: " + ep.Error())
			return
		}
		res = Response{"ok", "played " + note.Sample()}
		resb, em := json.Marshal(res)
		if em != nil {
			log.Println("could not marshal response: " + em.Error())
			return
		}
		ew := conn.WriteMessage(msgType, resb)
		if ew != nil {
			log.Println("could not write ws message: " + ew.Error())
		}
	})
}

// generate endpoint for starting pattern
func (this *simp) patternPlay() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := this.sequencer.Start()
		if err == nil {
			fmt.Fprintf(w, "{\"status\":\"ok\"}")
		} else {
			fmt.Fprintf(w, "{\"error\":\"%s\"}", err.Error())
		}
	}
	// return this.upgrade(func(conn *websocket.Conn, msgType int, msg []byte) {
	// 	this.sequencer.Start()
	// })
}

// generate endpoint for stopping pattern
func (this *simp) patternStop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		this.sequencer.Stop()
	}
	// return this.upgrade(func(conn *websocket.Conn, msgType int, msg []byte) {
	// 	this.sequencer.Stop()
	// })
}

// generate endpoint for editing pattern
func (this *simp) patternEdit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res Response
		pes := make([]PatternEdit, 0)
		dec := json.NewDecoder(r.Body)
		ed := dec.Decode(&pes)
		if ed != nil {
			log.Println("could not decode request body: " + ed.Error())
			return
		}
		for _, pe := range pes {
			err := this.AddTo(pe.Pos, pe.Note)
			if err != nil {
				log.Println("could not set note: " + err.Error())
				return
			}
		}
		res = Response{"ok", "note added"}
		resb, ee := json.Marshal(res)
		if ee != nil {
			log.Println("could not encode response: " + ee.Error())
		}
		buf := bytes.NewBuffer(resb)
		fmt.Fprintf(w, "%s", buf.String())
	}
}

func (this *simp) Connect(ch1 string, ch2 string) error {
	return this.engine.Connect(ch1, ch2)
}

func NewServer(webRoot string, audioRoot string) (Server, error) {
	// our pattern has 16384 sixteenth notes,
	// which means we have 1024 bars available
	// initialize tempo to 120 bpm (a typical
	// starting point for sequencers)
	engine := NewEngine()
	srv := &simp{
		audioRoot,
		engine,
		osc.NewOscServer(OSC_ADDR, OSC_PORT),
		NewSequencer(engine, PATTERN_LENGTH, Tempo(120), PATTERN_DIV),
	}
	// api handler
	api, ea := NewApi(audioRoot)
	if ea != nil {
		log.Println("could not create api: " + ea.Error())
		return nil, ea
	}
	// osc server
	pm := func(msg *osc.OscMessage) {
		if len(msg.Arguments) != 3 {
			log.Fatal("incorrect arguments to /sample/play (expects sii)")
		}
		samp := msg.Arguments[0].(string)
		pitch := msg.Arguments[1].(int32)
		gain := msg.Arguments[2].(int32)
		log.Printf("Note(%s, %d, %d)\n", samp, pitch, gain)
		note := NewNote(samp, pitch, gain)
		srv.engine.PlayNote(note)
	}
	srv.oscServer.AddMsgHandler("/sample/play", pm)
	go srv.oscServer.ListenAndDispatch()
	log.Printf("osc server listening on port %d\n", OSC_PORT)
	// setup handlers under default ServeMux
	fh := http.FileServer(http.Dir(webRoot))
	// static file server
	http.Handle("/", fh)
	// ReST endpoints
	http.HandleFunc("/samples", api.ListSamples())
	// websocket endpoints
	http.Handle("/sample/play", srv.playSample())
	http.HandleFunc("/pattern", srv.patternEdit())
	http.HandleFunc("/pattern/play", srv.patternPlay())
	http.HandleFunc("/pattern/stop", srv.patternStop())
	// add the audio root to the search path
	if 0 != srv.engine.AddDir(audioRoot) {
		log.Fatal("could not add dir " + audioRoot)
	}
	return srv, nil
}
