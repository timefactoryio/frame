package frame

import (
	"net/http"

	"github.com/timefactoryio/frame/fx"
)

type Frame struct {
	*fx.Fx
}

func NewFrame() *Frame {
	return &Frame{
		Fx: fx.NewFx(),
	}
}

func (f *Frame) Start(pathless string) {
	if pathless == "" {
		pathless = "http://localhost:1000"
	}

	f.BuildHello()
	f.Router().HandleFunc("/", f.HandleHello)
	go f.serve(pathless)
}

func (f *Frame) serve(pathless string) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", pathless)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		f.Router().ServeHTTP(w, r)
	})
	http.ListenAndServe(":1001", handler)
}
