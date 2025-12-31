package frame

import (
	"github.com/timefactoryio/frame/zero"
)

type Frame interface {
	zero.Zero
	Home(heading, github, x string)
	Text(path string)
	Slides(dir string)
}

type frame struct {
	zero.Zero
}

func NewFrame(pathlessUrl, apiUrl string) Frame {
	if pathlessUrl == "" {
		pathlessUrl = "http://localhost:1000"
	}
	if apiUrl == "" {
		apiUrl = "http://localhost:1001"
	}

	f := &frame{
		Zero: zero.NewZero(pathlessUrl, apiUrl),
	}
	return f
}
