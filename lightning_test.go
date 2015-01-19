package lightning

import (
	"testing"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine()
	if engine == nil {
		t.Fail()
	}
}
