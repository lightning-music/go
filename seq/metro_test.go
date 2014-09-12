package lightning

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestNewMetro(t *testing.T) {
	metro, err := NewMetro(Bpm(120), "1/16")
	assert.Equal(t, err, nil)
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

func TestSetTempo(t *testing.T) {
	metro, err := NewMetro(Bpm(120), "1/16")
	assert.Equal(t, err, nil)
	err = metro.SetTempo(150, "1/16")
	assert.Equal(t, err, nil)
}