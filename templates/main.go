package templates

import "github.com/timefactoryio/frame/zero"

type Templates interface {
	Style
	GithubLink(username string) *zero.One
	XLink(username string) *zero.One
	Home(heading, github, x string)
	README(content []byte) *zero.One
	Scroll() *zero.One
	BuildSlides(dir string) *zero.One
	BuildVideo(dir string) *zero.One
	Keyboard(asFrame bool) *zero.One
}

type templates struct {
	Style
	zero.Zero
}

func NewTemplates(zero zero.Zero) Templates {
	return &templates{
		Style: NewStyle(),
		Zero:  zero,
	}
}
