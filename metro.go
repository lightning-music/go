package lightning

import (
	"time"
)

// tempo in bpm
type Tempo uint32

type Pos uint64

type Metro interface {
	Ticker() chan Pos
	Stop()
}

type metroImpl struct {
	t *time.Ticker
	c chan Pos
}

func (m *metroImpl) Ticker() chan Pos {
	return m.c
}

func (m *metroImpl) Stop() {
	m.t.Stop()
}

func count(m *metroImpl) {
	var i Pos = 0
	for _ = range m.t.C {
		i++
		m.c <- i
	}
}


// Create a new metro
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
