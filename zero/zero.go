package zero

import (
	"encoding/json"
	"net/http"
)

type Zero struct {
	Fx
	Forge
	Element
	cachedJSON []byte
}

func NewZero(pathlessUrl, apiUrl string) *Zero {
	z := &Zero{
		Fx:      NewFx(pathlessUrl, apiUrl).(*fx),
		Forge:   NewForge().(*forge),
		Element: NewElement().(*element),
	}
	return z
}

// Encode marshals the input data to JSON and caches it for subsequent Out calls
func (z *Zero) Encode(input any) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	z.cachedJSON = data
	return nil
}

// Out writes the cached JSON response for http requests or returns an HTTP 500 error if no data is cached
func (z *Zero) Out(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	if z.cachedJSON == nil {
		http.Error(w, "No data encoded", http.StatusInternalServerError)
		return
	}
	w.Write(z.cachedJSON)
}
