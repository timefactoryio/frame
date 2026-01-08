package zero

import (
	"bytes"
	"html/template"
	"strings"
)

type Zero struct {
	Forge
	Element
	Circuit
}

func NewZero(pathless string) *Zero {
	return &Zero{
		Forge:   NewForge().(*forge),
		Element: NewElement().(*element),
		Circuit: NewCircuit(pathless).(*circuit),
	}
}

func (z *Zero) Home(logo, heading string) {
	logoBytes := z.ToBytes(logo)
	if logoBytes == nil {
		return
	}

	logoData := "data:image/svg+xml," + string(logoBytes)
	tmpl := template.Must(template.New("home").Parse(homeHtml))

	var buf strings.Builder
	if err := tmpl.Execute(&buf, map[string]string{
		"LOGO":    logoData,
		"HEADING": heading,
	}); err != nil {
		return
	}

	o := One(template.HTML(buf.String()))
	z.Build("", &o)
}

func (z *Zero) HomeWithFooter(logo, heading, link, icon, alt string) {
	logoBytes := z.ToBytes(logo)
	if logoBytes == nil {
		return
	}

	iconBytes := z.ToBytes(icon)
	if iconBytes == nil {
		return
	}

	logoData := "data:image/svg+xml," + string(logoBytes)
	iconData := "data:image/svg+xml," + string(iconBytes)

	tmpl := template.Must(template.New("home").Parse(homeWithFooterHtml))

	var buf strings.Builder
	if err := tmpl.Execute(&buf, map[string]string{
		"LOGO":    logoData,
		"HEADING": heading,
		"LINK":    link,
		"ICON":    iconData,
		"ALT":     alt,
	}); err != nil {
		return
	}

	o := One(template.HTML(buf.String()))
	z.Build("", &o)
}

func (z *Zero) Text(path string) {
	content := z.ToBytes(path)
	if content == nil {
		return
	}

	var buf bytes.Buffer
	if err := (*z.Markdown()).Convert(content, &buf); err != nil {
		return
	}

	html := buf.String()
	html = strings.ReplaceAll(html, "<p><img", "<img")
	html = strings.ReplaceAll(html, "\"></p>", "\">")
	html = strings.ReplaceAll(html, "\" /></p>", "\" />")
	html = strings.ReplaceAll(html, "\"/></p>", "\"/>")

	markdown := One(template.HTML(html))
	template := One(template.HTML(z.TextTemplate()))
	z.Build("text", &markdown, &template)
}

func (z *Zero) Slides(dir string) {
	prefix := z.Reader(dir)

	tmpl, err := template.New("slides").Parse(z.SlidesTemplate())
	if err != nil {
		return
	}

	var buf bytes.Buffer
	data := map[string]string{"PREFIX": prefix}
	if err := tmpl.Execute(&buf, data); err != nil {
		return
	}

	html := One(template.HTML(buf.String()))
	z.Build("slides", &html)
}
