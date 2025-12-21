package fx

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"strings"

	"github.com/timefactoryio/frame/zero"
)

//go:embed html/slides.html
var slidesHtml string

//go:embed html/text.html
var textHtml string

type Fx interface {
	Text(path string)
	Slides(dir string)
	Keyboard() *zero.One
	ToBytes(input string) []byte
	ToString(input string) string
	Pathless() string
	Api() string
}

type fx struct {
	*zero.Zero
	keyboard *zero.One
	pathless string
	api      string
}

func (f *fx) Pathless() string {
	return f.pathless
}

func (f *fx) Api() string {
	return f.api
}

func NewFx(pathless, apiUrl string) Fx {
	if pathless == "" {
		pathless = "http://localhost:1000"
	}
	if apiUrl == "" {
		apiUrl = "http://localhost:1001"
	}

	f := &fx{
		Zero:     zero.NewZero(),
		pathless: pathless,
		api:      apiUrl,
	}
	f.keyboard = f.BuildFromFile("html/keyboard.html", "")
	return f
}

func (f *fx) Keyboard() *zero.One {
	return f.keyboard
}

func (f *fx) Text(path string) {
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
	template := zero.One(template.HTML(textHtml))
	f.Build("text", &markdown, &template)
}

func (f *fx) Slides(dir string) {
	prefix := f.AddPath(dir)

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
	f.Build("slides", &html)
}

func (f *fx) Home(heading, github, x string) {
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

func (f *fx) buildFooter(github, x string) *zero.One {
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

func (f *fx) GithubLink(username string) *zero.One {
	if username == "" {
		return nil
	}
	logo := fmt.Sprintf("%s/img/gh", f.Api())
	href := fmt.Sprintf("https://github.com/%s", username)
	return f.LinkedIcon(href, logo, "GitHub")
}

func (f *fx) XLink(username string) *zero.One {
	if username == "" {
		return nil
	}
	logo := fmt.Sprintf("%s/img/x", f.Api())
	href := fmt.Sprintf("https://x.com/%s", username)
	return f.LinkedIcon(href, logo, "X")
}
