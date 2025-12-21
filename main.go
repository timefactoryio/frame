package frame

import (
	"github.com/timefactoryio/frame/one"
)

type Frame struct {
	*one.One
}

func NewFrame(pathlessUrl, apiURL string) *Frame {
	f := &Frame{
		One: one.NewOne(pathlessUrl, apiURL),
	}
	return f
}
