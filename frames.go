package frame

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"strings"

	"github.com/timefactoryio/frame/zero"
)

func (f *frame) Text(path string) {
	content := f.ToBytes(path)
	if content == nil {
		return
	}

	var buf bytes.Buffer
	if err := (*f.Markdown()).Convert(content, &buf); err != nil {
		return
	}

	html := buf.String()
	html = strings.ReplaceAll(html, "<p><img", "<img")
	html = strings.ReplaceAll(html, "\"></p>", "\">")
	html = strings.ReplaceAll(html, "\" /></p>", "\" />")
	html = strings.ReplaceAll(html, "\"/></p>", "\"/>")

	markdown := zero.One(template.HTML(html))
	template := zero.One(template.HTML(f.TextTemplate()))
	f.Build("text", &markdown, &template)
}

func (f *frame) Slides(dir string) {
	prefix := f.Reader(dir)

	tmpl, err := template.New("slides").Parse(f.SlidesTemplate())
	if err != nil {
		return
	}

	var buf bytes.Buffer
	data := map[string]string{"PREFIX": prefix}
	if err := tmpl.Execute(&buf, data); err != nil {
		return
	}

	html := zero.One(template.HTML(buf.String()))
	f.Build("slides", &html)
}

func (f *frame) Home(heading, github, x string) {
	logo := f.Api() + "/img/logo"
	img := f.Img(logo, "logo")
	h1 := f.H1(heading)
	css := f.CSS(`
  .home {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	height: 100%;
	width: 100%;
	text-align: center;
	box-sizing: border-box;
	overflow: hidden;
}
.home img {
	max-width: 95%;
	max-height: 30vh;
	width: auto;
	height: auto;
	display: block;
	object-fit: contain;
}
.home h1 {
	color: inherit;
	width: 100%;
	white-space: nowrap;
	overflow: hidden;
	font-size: clamp(2rem, 3vw, 3rem);
	margin: 0;
}
`)
	footer := f.buildFooter(github, x)
	f.Build("home", img, h1, footer, &css)
}

func (f *frame) buildFooter(github, x string) *zero.One {
	if github == "" && x == "" {
		return nil
	}
	footerCSS := f.CSS(`
.footer {
    display: flex;
    justify-content: center;
    gap: 1.5em;
    margin-top: 1.5em;
}
.footer img.icon {
    width: 2em;
    height: 2em;
    object-fit: contain;
}
`)
	elements := []*zero.One{&footerCSS}

	if github != "" {
		elements = append(elements, f.GithubLink(github))
	}
	if x != "" {
		elements = append(elements, f.XLink(x))
	}
	return f.Builder("footer", elements...)
}

func (f *frame) GithubLink(username string) *zero.One {
	if username == "" {
		return nil
	}
	logo := fmt.Sprintf("%s/img/gh", f.Api())
	href := fmt.Sprintf("https://github.com/%s", username)
	return f.LinkedImg(href, logo, "GitHub")
}

func (f *frame) XLink(username string) *zero.One {
	if username == "" {
		return nil
	}
	logo := fmt.Sprintf("%s/img/x", f.Api())
	href := fmt.Sprintf("https://x.com/%s", username)
	return f.LinkedImg(href, logo, "X")
}
