package fx

import (
	"fmt"
	"html"
	"html/template"
	"regexp"
	"strings"
)

type One template.HTML

var (
	style  = regexp.MustCompile(`(?s)<style>(.*?)</style>`)
	script = regexp.MustCompile(`(?s)<script>(.*?)</script>`)
)

func NewForge() Forge {
	return &forge{
		frames: []*One{},
	}
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
	f.Frames(f.Builder(class, elements...))
}

func (f *forge) Builder(class string, elements ...*One) *One {
	var b strings.Builder
	if class != "" {
		fmt.Fprintf(&b, `<div class="%s">`, html.EscapeString(class))
	}
	for _, el := range elements {
		b.WriteString(string(*el))
	}
	if class != "" {
		b.WriteString("</div>")
	}
	cleaned := f.consolidateAssets(b.String())
	result := One(template.HTML(cleaned))
	return &result
}

func (f *forge) Frames(frame ...*One) []*One {
	if len(frame) > 0 && frame[0] != nil {
		f.frames = append(f.frames, frame[0])
	}
	return f.frames
}

func (f *forge) consolidateAssets(s string) string {
	if styleMatches := style.FindAllStringSubmatch(s, -1); len(styleMatches) > 1 {
		var sb strings.Builder
		for _, m := range styleMatches {
			sb.WriteString(m[1])
			sb.WriteByte('\n')
		}
		s = fmt.Sprintf("<style>%s</style>%s", sb.String(), style.ReplaceAllString(s, ""))
	}
	if scriptMatches := script.FindAllStringSubmatch(s, -1); len(scriptMatches) > 1 {
		var sb strings.Builder
		for _, m := range scriptMatches {
			sb.WriteString(m[1])
			sb.WriteByte('\n')
		}
		s = fmt.Sprintf("%s<script>%s</script>", script.ReplaceAllString(s, ""), sb.String())
	}
	return s
}
