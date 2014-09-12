package seq

import (
	"errors"
	"fmt"
	"time"
)

// tempo in bpm
type Bpm float64

// metro position
type Pos uint64

// metronome
type Metro struct {
	Tempo Bpm                   `json:"tempo"`
	// bar divisor (number of ticks per bar)
	Bardiv string               `json:"div"`
	// channel that emits position
	Channel chan Pos
	// the underlying ticker driving the metro
	ticker *time.Ticker
	// send any int on this channel to tell the
	// metro to stop
	stop chan int
}

func (metro *Metro) Stop() {
	metro.ticker.Stop()
	metro.stop <- 1
}

// change the timing of a metro
func (metro *Metro) SetTempo(tempo Bpm, bardiv string) error {
	// how to switch out the current ticker
	// for one that uses the new tempo?
	// the old one is looping through the ticker with range,
	// so we should probably stop the old one first
	metro.Stop()
	// wait for it to signal that it is done
	<-metro.stop
	dur, err := duration(tempo, bardiv)
	if err != nil {
		return err
	}
	metro.ticker = time.NewTicker(dur)
	go count(metro)
	return nil
}

func duration(tempo Bpm, bardiv string) (time.Duration, error) {
	nsPerBar := 1000000000 * (240 / tempo)
	div, err := ParseDivisor(bardiv)
	if err != nil {
		return 0, err
	}
	dur := nsPerBar / Bpm(div)
	return time.Duration(dur), nil
}

// Create a new metro and start it
func NewMetro(tempo Bpm, bardiv string) (*Metro, error) {
	dur, err := duration(tempo, bardiv)
	if err != nil {
		return nil, err
	}
	// bar div scalar
	metro := Metro{
		tempo,
		bardiv,
		// Channel
		make(chan Pos, 1),
		// ticker
		time.NewTicker(dur),
		// stop
		make(chan int),
	}
	go count(&metro)
	return &metro, nil
}

func count(metro *Metro) {
	var pos Pos = 0
mainloop:
	for {
		select {
		case <-metro.ticker.C:
			metro.Channel <- pos
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
	divisors := []int{ 1, 2, 3, 4, 6, 8, 12, 16, 24, 32, 64, 128 }
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
