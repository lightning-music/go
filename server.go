package lightning

import (
	"encoding/json"
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

type PatternEdit struct {
	Pos  Pos  `json:"pos"`
	Note Note `json:"note"`
}

type Response struct {
	Status  string `json:"string"`
	Message string `json:"message"`
}

type Server interface {
	Connect(ch1 string, ch2 string) error
	Listen(addr string) error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type simp struct {
	audioRoot string
	engine    Engine
	oscServer *osc.OscServer
	metro     *Metro
	pattern   *Pattern
}

func (this *simp) Listen(addr string) error {
	return http.ListenAndServe(addr, nil)
}

func (this *simp) SetNote(pos Pos, note Note) error {
	return this.pattern.Set(pos, note)
}

func (this *simp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("could not upgrade http conn to ws: " + err.Error())
		return
	}

	// log.Println("connected to websocket endpoint")

	for {
		var res Response
		msgType, bs, err := conn.ReadMessage()

		if err != nil {
			log.Println("could not read ws message: " + err.Error())
			return
		}

		note, enp := ParseNote(bs)
		if enp != nil {
			log.Println("could not parse note: " + enp.Error())
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
	}
}

// generate the MetroFunc that wires the metro to
// the pattern and the audio engine
func genMetroFunc(s *simp) MetroFunc {
	return func(pos Pos) {
		notes := s.pattern.NotesAt(pos % PATTERN_LENGTH)
		for _, note := range notes {
			s.engine.PlayNote(note)
		}
	}
}

// generate endpoint for starting pattern
func patternPlay(s *simp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.metro.Start()
	}
}

// generate endpoint for stopping pattern
func patternStop(s *simp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.metro.Stop()
	}
}

// generate endpoint for editing pattern
func patternEdit(s *simp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("could not upgrade http conn to ws: " + err.Error())
			return
		}

		for {
			var res Response
			msgType, bs, err := conn.ReadMessage()

			if err != nil {
				log.Println("could not read ws message: " + err.Error())
				return
			}

			pe := new(PatternEdit)
			ed := json.Unmarshal(bs, pe)
			if ed != nil {
				log.Println("could not decode ws message: " + ed.Error())
				return
			}

			err = s.SetNote(pe.Pos, pe.Note)
			if err != nil {
				log.Println("could not set note: " + err.Error())
				return
			}

			res = Response{"ok", "note added"}
			resb, ee := json.Marshal(res)
			if ee != nil {
				log.Println("could not encode response: " + ee.Error())
			}
			conn.WriteMessage(msgType, resb)
		}
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
	metro := NewMetro(Bpm(120), PATTERN_DIV)
	pat := NewPattern(PATTERN_LENGTH)
	srv := &simp{
		audioRoot,
		NewEngine(),
		osc.NewOscServer(OSC_ADDR, OSC_PORT),
		metro,
		&pat,
	}
	// api handler
	api, ea := NewApi(audioRoot)
	if ea != nil {
		log.Println("could not create api: " + ea.Error())
		return nil, ea
	}
	// osc comm
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

	http.Handle("/", fh)
	http.Handle("/sample/play", srv)
	http.HandleFunc("/samples", api.ListSamples())
	http.HandleFunc("/pattern", patternEdit(srv))
	http.HandleFunc("/pattern/play", patternPlay(srv))
	http.HandleFunc("/pattern/stop", patternStop(srv))

	if 0 != srv.engine.AddDir(audioRoot) {
		log.Fatal("could not add dir " + audioRoot)
	}
	return srv, nil
}
