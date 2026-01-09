package frame

import (
	"encoding/json"
	"net/http"

	"github.com/timefactoryio/frame/zero"
)

type Frame struct {
	*zero.Zero
	Hello []byte
}

func NewFrame(pathless string) *Frame {
	f := &Frame{
		Zero: zero.NewZero(),
	}
	if pathless == "" {
		pathless = "http://localhost:1000"
	}

	go func() {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", pathless)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET")
			f.Router().ServeHTTP(w, r)
		})
		http.ListenAndServe(":1001", handler)
	}()
	f.Router().HandleFunc("/", f.HandleFrame)
	return f
}

func (f *Frame) Start() {
	frames := make([]string, 0, len(f.Frames()))
	for _, frame := range f.Frames() {
		if frame != nil {
			frames = append(frames, string(*frame))
		}
	}

	response := struct {
		Frames  []string        `json:"frames"`
		Layouts json.RawMessage `json:"layouts"`
	}{
		Frames:  frames,
		Layouts: json.RawMessage(f.Layouts()),
	}

	jsonData, _ := json.Marshal(response)
	f.Hello = f.Compress(jsonData)

}

func (f *Frame) HandleFrame(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Encoding", "gzip")
	w.Write(f.Hello)
}
