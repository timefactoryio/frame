package fx

import (
	"bytes"
	"html/template"
	"strings"
)

// func (fx *Fx) Logo(path string) template.HTML {
// 	if strings.ToLower(filepath.Ext(path)) == ".svg" {
// 		if b := fx.ToBytes(path); b != nil {
// 			return template.HTML(b)
// 		}
// 		return ""
// 	}

// 	if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
// 		fx.Read(path, "")
// 		name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
// 		path = fx.APIURL + "/" + name
// 	}

// 	alt := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
// 	return template.HTML(fmt.Sprintf(`<img src="%s" alt="%s">`,
// 		html.EscapeString(path),
// 		html.EscapeString(alt),
// 	))
// }

func (fx *Fx) Home(logo, heading string) {
	logoEmbed := fx.ToBytes(logo)
	if logoEmbed == nil {
		return
	}

	tmpl := template.Must(template.New("home").Parse(fx.HomeTemplate))

	var buf strings.Builder
	if err := tmpl.Execute(&buf, map[string]template.HTML{
		"LOGO":    template.HTML(string(logoEmbed)),
		"HEADING": template.HTML(heading),
	}); err != nil {
		return
	}

	o := One(template.HTML(buf.String()))
	fx.Build("", &o)
}

func (fx *Fx) Text(path string) {
	content := fx.ToBytes(path)
	if content == nil {
		return
	}

	var buf bytes.Buffer
	if err := (*fx.Markdown()).Convert(content, &buf); err != nil {
		return
	}

	html := buf.String()
	html = strings.ReplaceAll(html, "<p><img", "<img")
	html = strings.ReplaceAll(html, "\"></p>", "\">")
	html = strings.ReplaceAll(html, "\" /></p>", "\" />")
	html = strings.ReplaceAll(html, "\"/></p>", "\"/>")

	markdown := One(template.HTML(html))
	template := One(template.HTML(fx.TextTemplate))
	fx.Build("text", &markdown, &template)
}

func (fx *Fx) Slides(dir string) {
	prefix := fx.Reader(dir)
	tmpl, err := template.New("slides").Parse(fx.SlidesTemplate)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	data := map[string]string{"PREFIX": prefix}
	if err := tmpl.Execute(&buf, data); err != nil {
		return
	}

	html := One(template.HTML(buf.String()))
	fx.Build("slides", &html)
}
