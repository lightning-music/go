package lightning

import (
	"time"
)

// tempo in bpm
type Tempo uint32

// metro position
type Pos uint64

// metronome interface
type Metro interface {
	Ticker() chan Pos
}

type metroImpl struct {
	t *time.Ticker
	c chan Pos
}

func (m *metroImpl) Ticker() chan Pos {
	return m.c
}

func count(m *metroImpl) {
	var i Pos = 0
	for _ = range m.t.C {
		i++
		m.c <- i
	}
}


// Create a new metro and start it
// Tempo is in bpm and metro will tick at the rate of bar/div
func NewMetro(t Tempo, div uint32) Metro {
	m := new(metroImpl)
	// bars / sec
	nsPerBar := 1000000000 * (240 / t)
	dur := nsPerBar / Tempo(div)
	m.t = time.NewTicker(time.Duration(dur))
	m.c = make(chan Pos)
	go count(m)
	return m
}
