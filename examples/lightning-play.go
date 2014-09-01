package main

import (
	"github.com/lightning/lightning"
	"time"
)

func main() {
	engine := lightning.NewEngine()
	samp := "/home/brian/Audio/freesound/48223__slothrop__trumpetc2.wav"
	note := lightning.NewNote(samp, 60, 120)
	engine.PlayNote(note)
	time.Sleep(5 * time.Second)
}
