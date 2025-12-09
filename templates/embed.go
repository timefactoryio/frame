package templates

import _ "embed"

//go:embed css/home.css
var homeCSS string

//go:embed css/slides.css
var slidesCSS string

//go:embed css/footer.css
var footerCSS string

//go:embed css/text.css
var textCSS string

//go:embed css/keyboard.css
var keyboardCSS string

//go:embed css/video.css
var videoCSS string

type Style interface {
	HomeCSS() string
	SlidesCSS() string
	FooterCSS() string
	TextCSS() string
	KeyboardCSS() string
	VideoCSS() string
}

type style struct{}

func NewStyle() Style {
	return &style{}
}

func (s *style) VideoCSS() string {
	return videoCSS
}

func (s *style) HomeCSS() string {
	return homeCSS
}

func (s *style) SlidesCSS() string {
	return slidesCSS
}

func (s *style) FooterCSS() string {
	return footerCSS
}

func (s *style) TextCSS() string {
	return textCSS
}

func (s *style) KeyboardCSS() string {
	return keyboardCSS
}
