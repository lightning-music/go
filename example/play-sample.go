package main

import (
	"github.com/lightning/go"
	"time"
)

func main() {
	engine := lightning.NewEngine()
	engine.Connect("system:playback_1", "system:playback_2")
	note := lightning.NewNote("/home/brian/lightning/kits/meow.wav", 60, 72)
	engine.PlayNote(note)
	// engine.PlaySample("/home/brian/lightning/kits/meow.wav", 1.0, 1.0)
	time.Sleep(3000 * time.Millisecond)
}
