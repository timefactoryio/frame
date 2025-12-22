package one

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"strings"

	"github.com/timefactoryio/frame/fx"
	"github.com/timefactoryio/frame/zero"
)

//go:embed html/slides.html
var slidesHtml string

//go:embed html/text.html
var textHtml string

type One interface {
	Text(path string)
	Slides(dir string)
	Keyboard() *zero.One
	ToBytes(input string) []byte
	ToString(input string) string
	Pathless() string
	Api() string
}

type one struct {
	*zero.Zero
	fx.Fx
	keyboard *zero.One
	pathless string
	api      string
}

func (o *one) Pathless() string {
	return o.pathless
}

func (o *one) Api() string {
	return o.api
}

func NewOne(pathless, apiUrl string) One {
	if pathless == "" {
		pathless = "http://localhost:1000"
	}
	if apiUrl == "" {
		apiUrl = "http://localhost:1001"
	}

	f := &one{
		Zero:     zero.NewZero(),
		pathless: pathless,
		api:      apiUrl,
	}
	return f
}

func (o *one) Keyboard() *zero.One {
	return o.keyboard
}

func (o *one) Text(path string) {
	content := o.ToBytes(path)
	if content == nil {
		return
	}

	var buf bytes.Buffer
	if err := (*o.Markdown()).Convert(content, &buf); err != nil {
		return
	}

	html := buf.String()
	html = strings.ReplaceAll(html, "<p><img", "<img")
	html = strings.ReplaceAll(html, "\"></p>", "\">")
	html = strings.ReplaceAll(html, "\" /></p>", "\" />")
	html = strings.ReplaceAll(html, "\"/></p>", "\"/>")

	markdown := zero.One(template.HTML(html))
	template := zero.One(template.HTML(textHtml))
	o.Build("text", &markdown, &template)
}

func (o *one) Slides(dir string) {
	prefix := o.AddPath(dir)

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
	o.Build("slides", &html)
}

func (o *one) Home(heading, github, x string) {
	logo := o.Api() + "/img/logo"
	img := o.Img(logo, "logo")
	h1 := o.H1(heading)
	css := o.CSS(`
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
	footer := o.buildFooter(github, x)
	o.Build("home", img, h1, footer, &css)
}

func (o *one) buildFooter(github, x string) *zero.One {
	if github == "" && x == "" {
		return nil
	}
	footerCSS := o.CSS(`
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
		elements = append(elements, o.GithubLink(github))
	}
	if x != "" {
		elements = append(elements, o.XLink(x))
	}
	return o.Builder("footer", elements...)
}

func (o *one) GithubLink(username string) *zero.One {
	if username == "" {
		return nil
	}
	logo := fmt.Sprintf("%s/img/gh", o.Api())
	href := fmt.Sprintf("https://github.com/%s", username)
	return o.LinkedIcon(href, logo, "GitHub")
}

func (o *one) XLink(username string) *zero.One {
	if username == "" {
		return nil
	}
	logo := fmt.Sprintf("%s/img/x", o.Api())
	href := fmt.Sprintf("https://x.com/%s", username)
	return o.LinkedIcon(href, logo, "X")
}
