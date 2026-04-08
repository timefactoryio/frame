package fx

import (
	"encoding/json"
	"net/http"

	"github.com/timefactoryio/frame/zero"
)

type Fx struct {
	*zero.Zero
	Forge
	Element
	Circuit
	Hello []byte
}

type hello struct {
	Frames   []string        `json:"frames"`
	Keyboard string          `json:"keyboard"`
	Layouts  json.RawMessage `json:"layouts"`
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

	u := &hello{
		Frames:   frames,
		Keyboard: fx.Keyboard,
		Layouts:  json.RawMessage(fx.Layouts),
	}

	if jsonData, err := json.Marshal(u); err == nil {
		fx.Hello = fx.Compress(jsonData)
	}
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
