package fx

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/timefactoryio/frame/zero"
)

type Fx struct {
	*zero.Zero
	Forge
	Element
	Circuit
	// APIURL string
	Hello []byte
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
	var values []*Value
	kb := []byte(fx.Keyboard)
	values = append(values, &Value{Name: "keyboard", Type: "text/html", Size: len(kb), Data: kb})

	for i, frame := range fx.Frames() {
		if frame != nil {
			data := []byte(string(*frame))
			values = append(values, &Value{Name: strconv.Itoa(i), Type: "text/html", Size: len(data), Data: data})
		}
	}

	manifestJSON, _ := json.Marshal(values)
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, uint32(len(manifestJSON)))
	buf.Write(manifestJSON)
	for _, v := range values {
		buf.Write(v.Data)
	}

	fx.Hello = fx.Compress(buf.Bytes())
}

func (fx *Fx) HandleHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Encoding", "gzip")
	w.Write(fx.Hello)
}
