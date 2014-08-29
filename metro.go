package lightning

import (
	"errors"
	"fmt"
	"time"
)

// tempo in bpm
type Tempo uint64

// metro position
type Pos uint64

// bar divider for master metro creation
const divider int = 128 * 3

// metronome
type Metro struct {
	// channel that emits position
	Channel chan Pos
	ticker *time.Ticker
	stop chan int
}

// master
type Master struct {
	Metro
	slaves []*Slave
}

// slave
type Slave struct {
	// slave channel
	Channel chan Pos
	// clock divisor
	divisor int
}

func (master *Master) addSlave(slave *Slave) {
	master.slaves = append(master.slaves, slave)
}

// Create a new metro that is slaved to a master
// meter should be of the form "1/DIV"
// where DIV can be any of
// 1, 2, 3, 4, 5, 6, 7, 8,
// 12, 16, 24, 32, 64, 128
func (master *Master) NewSlave(meter string) (*Slave, error) {
	mult, err := ParseDivisor(meter)
	if err != nil || mult == 0 {
		return nil, err
	}
	slave := Slave{
		make(chan Pos, 1),
		mult,
	}
	master.addSlave(&slave)
	go sync(&slave, master)
	return &slave, nil
}

func (metro *Metro) Stop() {
	metro.ticker.Stop()
	metro.stop <- 1
}

// change the timing of a master clock
func (master *Master) SetTempo(tempo Tempo) error {
	// how to switch out the current ticker
	// for one that uses the new tempo?
	// the old one is looping through the ticker with range,
	// so we should probably stop the old one first
	master.Stop()
	// wait for it to signal that it is done
	<-master.Metro.stop
	master.ticker = time.NewTicker(duration(tempo))	
	go count(&master.Metro)
	return nil
}

func duration(tempo Tempo) time.Duration {
	nsPerBar := 1000000000 * (240 / tempo)
	dur := nsPerBar / Tempo(divider)
	return time.Duration(dur)
}

// Create a new metro and start it
// Tempo is in bpm and metro will tick at the rate of bar/div
func NewMaster(tempo Tempo) *Master {
	// bar div scalar
	master := Master{
		Metro{
			make(chan Pos, 1),
			time.NewTicker(duration(tempo)),
			make(chan int),
		},
		make([]*Slave, 4),
	}
	go count(&master.Metro)
	return &master
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

func sync(slave *Slave, master *Master) {
	var pos, rel Pos = 0, 0
	for _ = range master.Metro.Channel {
		// send on master channel
		master.Metro.Channel <- pos
		// send to slaves
		for _, slave := range master.slaves {
			if rel == 0 {
				slave.Channel <- pos
			}
		}
		// increment position and relative position
		pos++
		rel = Pos( int(rel + 1) % slave.divisor )
	}
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
	return divider / mult, nil
}
