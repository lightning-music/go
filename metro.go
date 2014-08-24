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
	if err != nil {
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

// Create a new metro and start it
// Tempo is in bpm and metro will tick at the rate of bar/div
func NewMaster(tempo Tempo) *Master {
	// bar div scalar
	nsPerBar := 1000000000 * (240 / tempo)
	dur := nsPerBar / Tempo(divider)
	master := Master{
		Metro{
			make(chan Pos, 1),
			time.NewTicker(time.Duration(dur)),
		},
		make([]*Slave, 4),
	}
	go count(&master.Metro)
	return &master
}

func count(metro *Metro) {
	var pos Pos = 0
	for _ = range metro.ticker.C {
		pos++
		metro.Channel <- pos
	}
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
	if scanned != 2 || numerator != 1 {
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
