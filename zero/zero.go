package zero

import (
	"bytes"
	"encoding/base64"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

type Zero struct {
	Forge
	Element
	Circuit
}

func NewZero() *Zero {
	return &Zero{
		Forge:   NewForge().(*forge),
		Element: NewElement().(*element),
		Circuit: NewCircuit().(*circuit),
	}
}

func (z *Zero) Home(logo, heading string) {
	logoEmbed := z.ToString(logo)
	if logoEmbed == "" {
		return
	}

	tmpl := template.Must(template.New("home").Parse(homeHtml))

	var buf strings.Builder
	if err := tmpl.Execute(&buf, map[string]template.HTML{
		"LOGO":    template.HTML(logoEmbed),
		"HEADING": template.HTML(heading),
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

func (z *Zero) Text2(path string) {
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

	// Replace image sources with base64 data URIs
	re := regexp.MustCompile(`<img[^>]+src="([^"]+)"`)
	html = re.ReplaceAllStringFunc(html, func(match string) string {
		imgSrc := re.FindStringSubmatch(match)[1]
		imgData := z.ToBytes(imgSrc)
		if imgData == nil {
			return match
		}

		mimeType := http.DetectContentType(imgData)
		encoded := base64.StdEncoding.EncodeToString(imgData)
		dataURI := "data:" + mimeType + ";base64," + encoded
		return strings.Replace(match, imgSrc, dataURI, 1)
	})

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
