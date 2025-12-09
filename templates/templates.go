package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/timefactoryio/frame/zero"
)

func (t *templates) Home(heading, github, x string) {
	logo := t.Api() + "/img/logo"
	img := t.Img(logo, "logo")
	h1 := t.H1(heading)
	css := t.CSS(t.HomeCSS())
	footer := t.buildFooter(github, x)
	t.Build("home", true, img, h1, footer, &css)
}

func (t *templates) buildFooter(github, x string) *zero.One {
	if github == "" && x == "" {
		return nil
	}

	footerCSS := t.CSS(t.FooterCSS())
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

func (t *templates) README(content []byte) *zero.One {
	var buf bytes.Buffer
	if err := (*t.Markdown()).Convert(content, &buf); err != nil {
		empty := zero.One("")
		return &empty
	}

	html := buf.String()
	html = strings.ReplaceAll(html, "<p><img", "<img")
	html = strings.ReplaceAll(html, "\"></p>", "\">")
	html = strings.ReplaceAll(html, "\" /></p>", "\" />")
	html = strings.ReplaceAll(html, "\"/></p>", "\"/>")

	markdown := zero.One(template.HTML(html))
	scroll := t.Scroll()
	css := t.CSS(t.TextCSS())
	result := t.Build("text", true, &markdown, scroll, &css)
	return result
}

func (t *templates) Scroll() *zero.One {
	js := t.JS(`
(function(){
  const frame = pathless.frame();
  const key = 'scroll';

  frame.scrollTop = pathless.state()[key] || 0;
  frame.addEventListener('scroll', () => pathless.update(key, frame.scrollTop));

  let speed = 0;
  const scroll = () => {
    if (!speed) return;
    frame.scrollBy({ top: speed });
    requestAnimationFrame(scroll);
  };

  pathless.onKey('w', () => { speed = -20; scroll(); }, () => { speed = 0; });
  pathless.onKey('s', () => { speed = 20; scroll(); }, () => { speed = 0; });
  pathless.onKey('a', () => { speed = -40; scroll(); }, () => { speed = 0; });
  pathless.onKey('d', () => { speed = 40; scroll(); }, () => { speed = 0; });
})();
`)
	return &js
}

func (t *templates) BuildSlides(dir string) *zero.One {
	prefix := t.AddPath(dir)
	css := t.CSS(t.SlidesCSS())
	js := t.JS(fmt.Sprintf(`
(function() {
  const frame = pathless.frame();
  let slides = [];
  let index = pathless.state()[%q] || 0;

  async function show(i) {
    if (!slides.length) return;
    index = ((i %% slides.length) + slides.length) %% slides.length;
    pathless.update(%q, index);

    const imgEl = frame.querySelector('img');
    if (!imgEl) return;

    const slide = slides[index];
    try {
      const { data } = await pathless.fetch(apiUrl + '/%s/' + slide);
      imgEl.src = data;
      imgEl.alt = slide;
    } catch {
      imgEl.alt = "Failed to load image";
    }
  }

  pathless.fetch(apiUrl + '/%s/order')
    .then(({ data }) => {
      slides = data || [];
      if (slides.length) show(index);
    });

  pathless.onKey('a', () => show(index - 1));
  pathless.onKey('d', () => show(index + 1));
})();
`, prefix, prefix, prefix, prefix))
	return t.Build("slides", true, t.Img("", ""), &css, &js)
}

func (t *templates) BuildVideo(filePath string) *zero.One {
	t.AddFile(filePath, "video")
	name := filepath.Base(filePath)
	name = name[:len(name)-len(filepath.Ext(name))]

	video := t.Video("")
	css := t.CSS(t.VideoCSS())
	js := t.JS(fmt.Sprintf(`
(function() {
    const el = pathless.frame().querySelector('video');
    if (!el) return;

    const state = pathless.state();
    el.volume = 1;
    el.src = apiUrl + '/video/%s#t=' + (state.t || 0);
    el.load();

    if (!state.paused) el.play().catch(() => {});

    pathless.keybind((k) => {
        if (k === ' ') {
            if (el.paused) {
                el.play().catch(() => {});
                pathless.update('paused', false);
            } else {
                el.pause();
                pathless.update('paused', true);
            }
        }
    });

    pathless.cleanup(() => {
        pathless.update('t', el.currentTime || 0);
        el.pause();
    });
})();
    `, name))
	return t.Build("video", true, video, &css, &js)
}
