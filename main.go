package frame

import (
	"html/template"

	"github.com/timefactoryio/frame/templates"
	"github.com/timefactoryio/frame/zero"
)

type One template.HTML

type Frame struct {
	templates.Templates
	*zero.Zero
	Hello map[string]any `json:"hello"`
}

func NewFrame(pathlessUrl, apiURL string) *Frame {
	f := &Frame{
		Zero:  zero.NewZero(pathlessUrl, apiURL),
		Hello: make(map[string]any),
	}
	f.Templates = templates.NewTemplates(f.Zero)
	return f
}
