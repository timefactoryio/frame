package frame

import (
	"encoding/json"
	"net/http"

	"github.com/timefactoryio/frame/zero"
)

type Frame struct {
	*zero.Zero
	response []byte
}

type universe struct {
	Frames  []string        `json:"frames"`
	Layouts json.RawMessage `json:"layouts"`
}

func NewFrame() *Frame {
	return &Frame{
		Zero: zero.NewZero(),
	}
}

func (f *Frame) Start(pathless string) {
	frames := make([]string, 0, len(f.Frames()))
	for _, frame := range f.Frames() {
		if frame != nil {
			frames = append(frames, string(*frame))
		}
	}

	u := &universe{
		Frames:  frames,
		Layouts: json.RawMessage(f.Layouts()),
	}

	if jsonData, err := json.Marshal(u); err == nil {
		f.response = f.Compress(jsonData)
	}

	if pathless == "" {
		pathless = "http://localhost:1000"
	}
	f.Router().HandleFunc("/", f.handleRoot)
	go f.serve(pathless)
}

func (f *Frame) serve(pathless string) {
	handler := f.corsMiddleware(pathless)
	http.ListenAndServe(":1001", handler)
}

func (f *Frame) corsMiddleware(pathless string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", pathless)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		f.Router().ServeHTTP(w, r)
	})
}

func (f *Frame) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Encoding", "gzip")
	w.Write(f.response)
}
