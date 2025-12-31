package zero

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type One template.HTML

func NewForge() Forge {
	f := &forge{
		frames:  []*One{},
		Element: NewElement().(*element),
	}
	// f.Keyboard(keyboardHtml)
	return f
}

type forge struct {
	frames     []*One
	framesJson []byte
	forgeMap   []byte
	Element
}

type Forge interface {
	Build(class string, elements ...*One)
	Builder(class string, elements ...*One) *One
	Frames(frame ...*One) []*One
	ToBytes(input string) []byte
	ToString(input string) string
	ToJSON()
	FrameJson() []byte
	Element
}

// func (f *forge) Keyboard(html string) {
// 	keyboard := f.HTML(f.ToString(html))
// 	f.Build("", keyboard)
// }

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

func (f *forge) ToBytes(input string) []byte {
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		resp, err := http.Get(input)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil
		}
		return b
	}
	b, err := os.ReadFile(input)
	if err != nil {
		return nil
	}
	return b
}

func (f *forge) ToString(input string) string {
	b := f.ToBytes(input)
	if b == nil {
		return ""
	}
	return string(b)
}

func (f *forge) ToJSON() {
	frameStrings := make([]string, 0, len(f.frames))
	for _, frame := range f.frames {
		if frame != nil {
			frameStrings = append(frameStrings, string(*frame))
		}
	}

	framesData, _ := json.Marshal(frameStrings)
	layoutsData := f.ToBytes(f.Layouts())

	helloData := map[string]json.RawMessage{
		"frames":  framesData,
		"layouts": layoutsData,
	}

	combinedData, _ := json.Marshal(helloData)

	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write(combinedData)
	w.Close()
	f.framesJson = buf.Bytes()
}

func (f *forge) FrameJson() []byte {
	return f.framesJson
}
