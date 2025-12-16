package templates

import (
	"fmt"

	"github.com/timefactoryio/frame/zero"
)

func (t *templates) Home(heading, github, x string) {
	logo := t.Api() + "/img/logo"
	img := t.Img(logo, "logo")
	h1 := t.H1(heading)
	css := t.CSS(`
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
	footer := t.buildFooter(github, x)
	t.Build("home", true, img, h1, footer, &css)
}

func (t *templates) buildFooter(github, x string) *zero.One {
	if github == "" && x == "" {
		return nil
	}
	footerCSS := t.CSS(`
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
		elements = append(elements, t.GithubLink(github))
	}
	if x != "" {
		elements = append(elements, t.XLink(x))
	}
	return t.Build("footer", false, elements...)
}

func (t *templates) GithubLink(username string) *zero.One {
	if username == "" {
		return nil
	}
	logo := fmt.Sprintf("%s/img/gh", t.Api())
	href := fmt.Sprintf("https://github.com/%s", username)
	return t.LinkedIcon(href, logo, "GitHub")
}

func (t *templates) XLink(username string) *zero.One {
	if username == "" {
		return nil
	}
	logo := fmt.Sprintf("%s/img/x", t.Api())
	href := fmt.Sprintf("https://x.com/%s", username)
	return t.LinkedIcon(href, logo, "X")
}
