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
	Build(class string, elements ...*One)
	Builder(class string, elements ...*One) *One
	Frames(frame ...*One) []*One
}

func (f *forge) Build(class string, elements ...*One) {
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

	f.Frames(&result)
}

func (f *forge) Builder(class string, elements ...*One) *One {
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

func (f *forge) Frames(frame ...*One) []*One {
	if len(frame) > 0 && frame[0] != nil {
		f.frames = append(f.frames, frame[0])
	}
	return f.frames
}
