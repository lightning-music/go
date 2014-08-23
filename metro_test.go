package lightning

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestNewMetro(t *testing.T) {
	metro := NewMetro(120, 128)
	var pos Pos = 0
	ticker := metro.Ticker()
	for ; pos < 3; {
		pos = <-ticker
	}
	assert.Equal(t, int(pos), 3)
}
