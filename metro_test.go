package lightning

import "testing"

func TestNewMetro(t *testing.T) {
	metro := NewMetro(120, 128)
	metro.Stop()
}
