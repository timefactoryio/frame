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

// //go:embed html/keyboard.html
// var keyboardHtml string

type Templates interface {
	Home(heading, github, x string)
	Text(path string)
	Slides(dir string)
	BuildFromFile(html, class string) *zero.One
	Keyboard(html string) *zero.One
}

type templates struct {
	*zero.Zero
}

func NewTemplates(zero *zero.Zero) Templates {
	return &templates{
		Zero: zero,
	}
}

func (t *templates) Keyboard(html string) *zero.One {
	keyboard := t.BuildFromFile(html, "")
	return keyboard
}

func (t *templates) BuildFromFile(html, class string) *zero.One {
	file := t.HTML(t.ToString(html))
	final := t.Builder(class, file)
	return final
}

func (t *templates) Text(path string) {
	content := t.ToBytes(path)
	if content == nil {
		return
	}

	var buf bytes.Buffer
	if err := (*t.Markdown()).Convert(content, &buf); err != nil {
		return
	}

	html := buf.String()
	html = strings.ReplaceAll(html, "<p><img", "<img")
	html = strings.ReplaceAll(html, "\"></p>", "\">")
	html = strings.ReplaceAll(html, "\" /></p>", "\" />")
	html = strings.ReplaceAll(html, "\"/></p>", "\"/>")

	markdown := zero.One(template.HTML(html))
	template := zero.One(template.HTML(textHtml))
	t.Build("text", &markdown, &template)
}

func (t *templates) Slides(dir string) {
	prefix := t.AddPath(dir)

	tmpl, err := template.New("slides").Parse(slidesHtml)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	data := map[string]string{"PREFIX": prefix}
	if err := tmpl.Execute(&buf, data); err != nil {
		return
	}

	html := zero.One(template.HTML(buf.String()))
	t.Build("slides", &html)
}
