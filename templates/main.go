package templates

import (
	"bytes"
	_ "embed"
	"html/template"
	"strings"

	"github.com/timefactoryio/frame/zero"
)

//go:embed html/slides.html
var slidesHtml string

//go:embed html/text.html
var textHtml string

type Templates interface {
	Home(heading, github, x string)
	Text(content []byte) *zero.One
	Slides(dir string) *zero.One
	BuildFromFile(html, class string, asFrame bool) *zero.One
}

type templates struct {
	zero.Zero
}

func NewTemplates(zero zero.Zero) Templates {
	return &templates{
		Zero: zero,
	}
}

func (t *templates) BuildFromFile(html, class string, asFrame bool) *zero.One {
	file := t.HTML(t.ToString(html))
	final := t.Build(class, asFrame, file)
	return final
}

func (t *templates) Text(content []byte) *zero.One {
	var buf bytes.Buffer
	if err := (*t.Markdown()).Convert(content, &buf); err != nil {
		empty := zero.One("")
		return &empty
	}

	html := buf.String()
	html = strings.ReplaceAll(html, "<p><img", "<img")
	html = strings.ReplaceAll(html, "\"></p>", "\">")
	html = strings.ReplaceAll(html, "\" /></p>", "\" />")
	html = strings.ReplaceAll(html, "\"/></p>", "\"/>")

	markdown := zero.One(template.HTML(html))
	template := zero.One(template.HTML(textHtml))
	result := t.Build("text", true, &markdown, &template)
	return result
}

func (t *templates) Slides(dir string) *zero.One {
	prefix := t.AddPath(dir)

	tmpl, err := template.New("slides").Parse(slidesHtml)
	if err != nil {
		empty := zero.One("")
		return &empty
	}

	var buf bytes.Buffer
	data := map[string]string{"PREFIX": prefix}
	if err := tmpl.Execute(&buf, data); err != nil {
		empty := zero.One("")
		return &empty
	}

	html := zero.One(template.HTML(buf.String()))
	return t.Build("slides", true, &html)
}
