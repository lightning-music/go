package lightning

// Engine is the interface to the liblightning C bindings.
type Engine interface {
	Connect(ch1 string, ch2 string) error
	AddDir(file string) int
	PlaySample(file string, pitch float64, gain float64) error
	PlayNote(note Note) error
	ExportStart(file string) int
	ExportStop() int
}
