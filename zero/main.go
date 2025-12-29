package zero

import (
	"net/http"
)

type Zero interface {
	Forge
	Circuit
	Start()
}

type zero struct {
	Forge
	Circuit
}

func NewZero(pathless, apiUrl string) Zero {
	return &zero{
		Circuit: NewCircuit(pathless, apiUrl).(*circuit),
		Forge:   NewForge().(*forge),
	}
}

func (z *zero) Start() {
	z.ToJSON()
	go func() {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", z.Pathless())
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET")
			z.Router().ServeHTTP(w, r)
		})
		http.ListenAndServe(":1001", handler)
	}()
	z.Router().HandleFunc("/", z.HandleFrame)
}

func (z *zero) HandleFrame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Encoding", "gzip")
	w.Write(z.FrameJson())
}
