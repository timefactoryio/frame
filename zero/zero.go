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

func (z *Zero) SVGToEmbed(svg []byte) string {
	if svg == nil {
		return ""
	}
	return string(svg)
}

func (z *Zero) Home(logo, heading string) {
	logoEmbed := z.ToString(logo)
	if logoEmbed == "" {
		return
	}

	tmpl := template.Must(template.New("home").Parse(homeHtml))

	var buf strings.Builder
	if err := tmpl.Execute(&buf, map[string]string{
		"LOGO":    logoEmbed,
		"HEADING": heading,
	}); err != nil {
		return
	}

	o := One(template.HTML(buf.String()))
	z.Build("", &o)
}

func (z *Zero) HomeWithFooter(logo, heading, link, icon, alt string) {
	logoBytes := z.ToString(logo)
	iconBytes := z.ToString(icon)
	tmpl := template.Must(template.New("home").Parse(homeWithFooterHtml))

	var buf strings.Builder
	if err := tmpl.Execute(&buf, map[string]string{
		"LOGO":    logoBytes,
		"HEADING": heading,
		"LINK":    link,
		"ICON":    iconBytes,
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
