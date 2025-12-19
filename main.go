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
	Hello map[string][]byte `json:"map"`
}

func NewFrame(pathlessUrl, apiURL string) *Frame {
	f := &Frame{
		Zero:  zero.NewZero(pathlessUrl, apiURL),
		Hello: make(map[string][]byte),
	}
	f.Templates = templates.NewTemplates(f.Zero)
	f.Router().HandleFunc("/frame", f.HandleFrame)
	f.Router().HandleFunc("/hello", f.HandleHello)

	return f
}

func (f *Frame) HandleFrame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Encoding", "gzip")
	w.Write(f.Frames())
}

func (f *Frame) HandleHello(w http.ResponseWriter, r *http.Request) {
	f.Hello["frames"] = f.Frames()
	f.Hello["keyboard"] = f.KeyboardBytes()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Encoding", "gzip")
	json.NewEncoder(w).Encode(f.Hello)
}
