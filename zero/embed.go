package zero

import (
	_ "embed"
)

//go:embed embed/slides.html
var slidesHtml string

//go:embed embed/text.html
var textHtml string

//go:embed embed/layouts.json
var layoutsJson []byte

//go:embed embed/home.html
var homeHtml string

//go:embed embed/timefactory.svg
var timefactory []byte

type Embed interface {
	TextTemplate() string
	SlidesTemplate() string
	HomeTemplate() string
	Layouts() []byte
	Timefactory() []byte
}

func NewEmbed() Embed {
	return &embed{}
}

type embed struct{}

func (e *embed) Layouts() []byte {
	return layoutsJson
}

func (e *embed) HomeTemplate() string {
	return homeHtml
}

func (e *embed) TextTemplate() string {
	return textHtml
}

func (e *embed) SlidesTemplate() string {
	return slidesHtml
}

func (e *embed) Timefactory() []byte {
	return timefactory
}
