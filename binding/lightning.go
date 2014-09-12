// go binding for the lightning audio engine
package binding

// #cgo CFLAGS: -Wall -O2
// #cgo LDFLAGS: -L. -llightning -lm -ljack -lsndfile -lpthread -lsamplerate
// #include <lightning/lightning.h>
// #include <lightning/types.h>
import "C"

import (
	"errors"
	"math"
)

type Pitch float64
type Gain  float64

type impl struct {
	handle C.Lightning
}

func (this *impl) AddDir(file string) int {
	return int(C.Lightning_add_dir(this.handle, C.CString(file)))
}

func (this *impl) PlaySample(file string, pitch Pitch, gain Gain) error {
	err := C.Lightning_play_sample(
		this.handle, C.CString(file), C.pitch_t(pitch), C.gain_t(gain),
	)
	if err != 0 {
		return errors.New("could not play sample")
	} else {
		return nil
	}
}

func getPitch(note Note) Pitch {
	return Pitch(math.Pow(2.0, (float64(note.Number) - 60.0) / 12.0))
}

func (this *impl) PlayNote(note Note) error {
	pitch := getPitch(note)
	gain := Gain(Gain(note.Velocity) / 127.0)
	return this.PlaySample(note.Sample, pitch, gain)
}

func (this *impl) ExportStart(file string) int {
	return int(C.Lightning_export_start(
		this.handle, C.CString(file),
	))
}

func (this *impl) ExportStop() int {
	return int(C.Lightning_export_stop(this.handle))
}

func NewEngine() Engine {
	instance := new(impl)
	instance.handle = C.Lightning_init()
	return instance
}
