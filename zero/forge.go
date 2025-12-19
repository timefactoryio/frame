package zero

import (
	"fmt"
	"html"
	"html/template"
	"regexp"
	"strings"
)

type One template.HTML

func NewForge() Forge {
	f := &forge{
		frames: []*One{},
	}
	return f
}

type forge struct {
	frames []*One
}

type Forge interface {
	Build(class string, updateIndex bool, elements ...*One) *One
	JS(js string) One
	CSS(css string) One
	UpdateIndex(*One)
	Frames() []*One
}

func (f *forge) Build(class string, updateIndex bool, elements ...*One) *One {
	var b strings.Builder
	for _, el := range elements {
		b.WriteString(string(*el))
	}

	var htmlOut string
	if class == "" {
		htmlOut = b.String()
	} else {
		consolidatedContent := b.String()
		htmlOut = fmt.Sprintf(`<div class="%s">%s</div>`, html.EscapeString(class), consolidatedContent)
	}
	cleaned := f.consolidateAssets(htmlOut)
	result := One(template.HTML(cleaned))

	if updateIndex {
		f.UpdateIndex(&result)
	}
	return &result
}

func (f *forge) consolidateAssets(html string) string {
	styleRe := regexp.MustCompile(`(?s)<style>(.*?)</style>`)
	styleMatches := styleRe.FindAllStringSubmatch(html, -1)
	var styleBlock string
	if len(styleMatches) > 1 {
		for _, m := range styleMatches {
			styleBlock += m[1] + "\n"
		}
		html = styleRe.ReplaceAllString(html, "")
		if styleBlock != "" {
			html = fmt.Sprintf("<style>%s</style>%s", styleBlock, html)
		}
	}
	scriptRe := regexp.MustCompile(`(?s)<script>(.*?)</script>`)
	scriptMatches := scriptRe.FindAllStringSubmatch(html, -1)
	var scriptBlock string
	if len(scriptMatches) > 1 {
		for _, m := range scriptMatches {
			scriptBlock += m[1] + "\n"
		}
		html = scriptRe.ReplaceAllString(html, "")
		if scriptBlock != "" {
			html = fmt.Sprintf("%s<script>%s</script>", html, scriptBlock)
		}
	}
	return html
}

func (f *forge) JS(js string) One {
	var b strings.Builder
	b.WriteString(`<script>`)
	b.WriteString(js)
	b.WriteString(`</script>`)
	return One(template.HTML(b.String()))
}

func (f *forge) CSS(css string) One {
	var b strings.Builder
	b.WriteString(`<style>`)
	b.WriteString(css)
	b.WriteString(`</style>`)
	return One(template.HTML(b.String()))
}

func (f *forge) UpdateIndex(frame *One) {
	if frame != nil {
		f.frames = append(f.frames, frame)
	}
}

func (f *forge) Frames() []*One {
	return f.frames
}
