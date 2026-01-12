package zero

import (
	_ "embed"
)

//go:embed embed/home.html
var homeHtml string

//go:embed embed/layouts.json
var layoutsJson []byte

//go:embed embed/slides.html
var slidesHtml string

//go:embed embed/text.html
var textHtml string

type Zero struct {
	SlidesTemplate string
	TextTemplate   string
	HomeTemplate   string
	Layouts        []byte
}

func NewZero() *Zero {
	return &Zero{
		SlidesTemplate: slidesHtml,
		TextTemplate:   textHtml,
		HomeTemplate:   homeHtml,
		Layouts:        layoutsJson,
	}
}
