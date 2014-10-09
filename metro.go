package lightning

import (
	"errors"
	"fmt"
	"time"
)

// metro position
type Pos uint64

type MetroFunc func(pos Pos)

// metronome
type Metro struct {
	Tempo   Tempo    `json:"tempo"`
	Bardiv  string `json:"div"`
	Channel chan Pos
	F       MetroFunc
	ticker  *time.Ticker
	stop    chan int
	playing bool
}

func (this *Metro) Stop() {
	this.ticker.Stop()
	// signal the count gorouting to exit
	this.stop <- 1
	// wait for it to exit
	<-this.stop
	this.playing = false
}

// change the timing of a metro
func (this *Metro) SetTempo(tempo Tempo, bardiv string) {
	this.Tempo = tempo
	this.Bardiv = bardiv
}

func (this *Metro) SetFunc(f MetroFunc) {
	this.F = f
}

func (this *Metro) Start() error {
	if this.playing {
		return nil
	}
	this.playing = true
	dur, err := duration(this.Tempo, this.Bardiv)
	if err != nil {
		return err
	}
	this.ticker = time.NewTicker(dur)
	go count(this)
	return nil
}

func duration(tempo Tempo, bardiv string) (time.Duration, error) {
	nsPerBar := 1000000000 * (240 / tempo)
	div, err := ParseDivisor(bardiv)
	if err != nil {
		return 0, err
	}
	dur := nsPerBar / Tempo(div)
	return time.Duration(dur), nil
}

// Create a new metro and start it
func NewMetro(tempo Tempo, bardiv string) *Metro {
	// bar div scalar
	return &Metro{
		tempo,
		bardiv,
		make(chan Pos, 1),
		nil,
		nil,
		make(chan int),
		false,
	}
}

func count(metro *Metro) {
	var pos Pos = 0
mainloop:
	for {
		select {
		case <-metro.ticker.C:
			metro.Channel <- pos
			if metro.F != nil {
				metro.F(pos)
			}
			pos++
		case <-metro.stop:
			// break out of mainloop and signal done
			break mainloop
		}
	}
	metro.stop <- 1
}

// meter should be of the form "1/DIV"
// where DIV can be any of
// 1, 2, 3, 4, 6, 8,
// 12, 16, 24, 32, 64, 128
func ParseDivisor(meter string) (int, error) {
	var numerator, mult int
	scanned, err := fmt.Sscanf(meter, "%d/%d", &numerator, &mult)
	if err != nil {
		return 0, err
	}
	if scanned != 2 || numerator != 1 || mult == 0 {
		return 0, errors.New("invalid meter")
	}
	// acceptable clock divisors for slave syncing
	divisors := []int{1, 2, 3, 4, 6, 8, 12, 16, 24, 32, 64, 128}
	valid := false
	for _, div := range divisors {
		if mult == div {
			valid = true
		}
	}
	if !valid {
		return 0, errors.New("invalid bar divisor")
	}
	return mult, nil
}
