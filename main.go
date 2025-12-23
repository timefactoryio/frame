package frame

import (
	"github.com/timefactoryio/frame/fx"
	"github.com/timefactoryio/frame/zero"
)

type Frame struct {
	zero.Zero
	fx.Fx
}

func NewFrame(pathlessUrl, apiUrl string) *Frame {
	if pathlessUrl == "" {
		pathlessUrl = "http://localhost:1000"
	}
	if apiUrl == "" {
		apiUrl = "http://localhost:1001"
	}

	f := &Frame{
		Fx: fx.NewFx(),
	}
	return f
}
