// go binding for the lightning audio engine
package lightning

// #cgo CFLAGS: -Wall -O2
// #cgo LDFLAGS: -L. -llightning -lm -ljack -lsndfile -lpthread -lsamplerate
// #include <lightning/lightning.h>
// #include <lightning/types.h>
import "C"

type Pitch float32
type Gain  float32

type impl struct {
	handle C.Lightning
}

func (this *impl) AddDir(file string) int {
	return int(C.Lightning_add_dir(this.handle, C.CString(file)))
}

func (this *impl) PlaySample(file string, pitch Pitch, gain Gain) int {
	return int(C.Lightning_play_sample(
		this.handle, C.CString(file), C.pitch_t(pitch), C.gain_t(gain),
	))
}

func (this *impl) ExportStart(file string) int {
	return int(C.Lightning_export_start(
		this.handle, C.CString(file),
	))
}

func (this *impl) ExportStop() int {
	return int(C.Lightning_export_stop(this.handle))
}

type Engine interface {
	AddDir(file string) int
	PlaySample(file string, pitch Pitch, gain Gain) int
	ExportStart(file string) int
	ExportStop() int
}

func NewEngine() Engine {
	instance := new(impl)
	instance.handle = C.Lightning_init()
	return instance
}
