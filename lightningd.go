package main

import (
	"encoding/json"
	"flag"
	"github.com/lightning/go/api"
	"github.com/lightning/go/seq"
	"github.com/lightning/go/types"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func main() {
	defaultRoot := path.Join(os.Getenv("HOME"), "lightning", "www")
	bind := flag.String("bind", "localhost:3428", "bind address")
	root := flag.String("root", defaultRoot, "web root")
	flag.Parse()
	server, err := api.NewServer(*root, "/home/brian/Audio/freesound")
	if err != nil {
		panic(err.Error())
	}
	log.Printf("serving static content from %s\n", *root)
	log.Printf("binding to %s\n", *bind)
	server.Listen(*bind)

	/* setup a pattern from a chunk of json */
	pat := seq.NewPattern(0)
	content, err := ioutil.ReadFile("pat.json")
	if err != nil {
		panic("could not read pat.json")
	}
	err = json.Unmarshal(content, &pat)
	if err != nil {
		panic("could not parse pat.json")
	}

	/* setup clocks */
	// metro, err := lightning.NewMetro(120, "1/16")
	// if err != nil {
	// 	panic("Could not create slave")
	// }

	// engine := lightning.NewEngine()
	// notes := make(chan lightning.Note, 64)
	// go playNotes(engine, notes)

	// for pos := range metro.Channel {
	// 	ns := pat.NotesAt(lightning.Pos(int(pos) % pat.Length))
	// 	for _, note := range ns {
	// 		notes <- note
	// 	}
	// }
}

func playNotes(engine types.Engine, notes chan types.Note) {
	for note := range notes {
		engine.PlayNote(note)
	}
}
