package lightning

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestNewMaster(t *testing.T) {
	metro := NewMaster(120)
	var pos Pos = 0
	for ; pos < 3; {
		pos = <-metro.Channel
	}
	assert.Equal(t, int(pos), 3)
}

func TestParseDivisor(t *testing.T) {
	_, err := ParseDivisor("1/2")
	assert.Equal(t, err, nil)
	_, err = ParseDivisor("2/3")
	assert.NotEqual(t, err, nil)
	_, err = ParseDivisor("2_3")
	assert.NotEqual(t, err, nil)
	_, err = ParseDivisor("1/5")
	assert.NotEqual(t, err, nil)
}

func TestNewSlave(t *testing.T) {
	metro := NewMaster(120)
	
	slave, err := metro.NewSlave("3/abc")
	// not sure why assert.Equal fails here
	if slave != nil {
		t.Fail()
	}
	assert.NotEqual(t, err, nil)

	slave, err = metro.NewSlave("1/16")
	assert.NotEqual(t, slave, nil)
	assert.Equal(t, err, nil)

	if slave == nil || slave.Channel == nil {
		t.Fail()
	}

	var pos Pos = 0
	for ; pos < 4; {
		pos = <-slave.Channel
	}
	assert.Equal(t, int(pos), 4)
}
