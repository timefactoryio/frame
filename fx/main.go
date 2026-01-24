package fx

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/timefactoryio/frame/zero"
)

type Fx struct {
	*zero.Zero
	Forge
	Element
	Circuit
	Hello []byte
}

type universe struct {
	Frames   []string        `json:"frames"`
	Layouts  json.RawMessage `json:"layouts"`
	Keyboard string          `json:"keyboard"`
	Tab      string          `json:"tab,omitempty"`
}

func NewFx() *Fx {
	return &Fx{
		Forge:   NewForge().(*forge),
		Element: NewElement().(*element),
		Circuit: NewCircuit().(*circuit),
		Zero:    zero.NewZero(),
	}
}

func (fx *Fx) BuildHello() {
	frames := make([]string, 0, len(fx.Frames()))
	for _, frame := range fx.Frames() {
		if frame != nil {
			frames = append(frames, string(*frame))
		}
	}

	u := &universe{
		Frames:   frames,
		Layouts:  json.RawMessage(fx.Layouts),
		Keyboard: fx.Keyboard,
	}

	if jsonData, err := json.Marshal(u); err == nil {
		fx.Hello = fx.Compress(jsonData)
	}
}

func (fx *Fx) BuildHelloTest(keyboardPath, tabPath string) error {
	frames := make([]string, 0, len(fx.Frames()))
	for _, frame := range fx.Frames() {
		if frame != nil {
			frames = append(frames, string(*frame))
		}
	}

	kb, err := os.ReadFile(keyboardPath)
	if err != nil {
		return err
	}

	u := &universe{
		Frames:   frames,
		Layouts:  json.RawMessage(fx.Layouts),
		Keyboard: string(kb),
		Tab:      tabPath,
	}

	if jsonData, err := json.Marshal(u); err == nil {
		fx.Hello = fx.Compress(jsonData)
	}

	return nil
}

func (fx *Fx) HandleHello(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Encoding", "gzip")
	w.Write(fx.Hello)
}
