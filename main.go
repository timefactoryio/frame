package frame

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/timefactoryio/frame/templates"
	"github.com/timefactoryio/frame/zero"
)

type One template.HTML

type Frame struct {
	templates.Templates
	zero.Zero
	Hello map[string]json.RawMessage `json:"hello"`
}

func NewFrame(pathlessUrl, apiURL string) *Frame {
	f := &Frame{
		Zero:  zero.NewZero(pathlessUrl, apiURL),
		Hello: make(map[string]json.RawMessage),
	}
	f.Templates = templates.NewTemplates(f.Zero)
	f.Router().HandleFunc("/frame", f.HandleFrame)
	f.Router().HandleFunc("/hello", f.HandleHello)
	return f
}

func (f *Frame) HandleFrame(w http.ResponseWriter, r *http.Request) {
	framesJSON, _ := json.Marshal(f.Frames())
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(framesJSON)
}

func (f *Frame) HandleHello(w http.ResponseWriter, r *http.Request) {
	framesJSON, _ := json.Marshal(f.Frames())
	keyboardJSON, _ := json.Marshal(f.Keyboard())

	f.Hello["frames"] = framesJSON
	f.Hello["keyboard"] = keyboardJSON

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(f.Hello)
}
