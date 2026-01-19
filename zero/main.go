package zero

import (
	_ "embed"
)

//go:embed embed/home.html
var homeHtml string

//go:embed embed/slides.html
var slidesHtml string

//go:embed embed/text.html
var textHtml string

//go:embed embed/keyboard.html
var keyboardHtml string

//go:embed embed/layouts.json
var layoutsJson []byte

//go:embed embed/focus.json
var focus []byte

type Zero struct {
	SlidesTemplate string
	TextTemplate   string
	HomeTemplate   string
	Keyboard       string
	Layouts        []byte
	Focus          []byte
}

func NewZero() *Zero {
	return &Zero{
		SlidesTemplate: slidesHtml,
		TextTemplate:   textHtml,
		HomeTemplate:   homeHtml,
		Layouts:        layoutsJson,
		Keyboard:       keyboardHtml,
		Focus:          focus,
	}
}
