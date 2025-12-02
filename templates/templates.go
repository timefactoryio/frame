package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/timefactoryio/frame/zero"
)

func (t *templates) Landing(heading, github, x string) {
	logo := t.ApiUrl() + "/img/logo"
	img := t.Img(logo, "logo")
	h1 := t.H1(heading)
	css := t.CSS(t.ZeroCSS())
	footer := t.buildFooter(github, x)
	t.Build("zero", true, &css, img, h1, footer)
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
	logo := fmt.Sprintf("%s/img/gh", t.ApiUrl())
	href := fmt.Sprintf("https://github.com/%s", username)
	return t.LinkedIcon(href, logo, "GitHub")
}

func (t *templates) XLink(username string) *zero.One {
	if username == "" {
		return nil
	}
	logo := fmt.Sprintf("%s/img/x", t.ApiUrl())
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
	js := `
(function(){
  const { frame, state } = pathless.ctx();
  
  frame.scrollTop = state.scroll || 0;
  frame.addEventListener('scroll', () => pathless.update('scroll', frame.scrollTop));
  
  const speeds = { w: -20, s: 20, a: -40, d: 40 };
  let speed = 0, raf = 0;
  
  const tick = () => {
    if (!speed) return raf = 0;
    frame.scrollBy({ top: speed });
    raf = requestAnimationFrame(tick);
  };
  
  pathless.onKey(k => {
    if (speeds[k]) {
      speed = speeds[k];
      if (!raf) tick();
    }
  });
  
  document.addEventListener('keyup', e => speeds[e.key] && (speed = 0));
})();
`
	result := zero.One(template.HTML(fmt.Sprintf(`<script>%s</script>`, js)))
	return &result
}

func (t *templates) BuildSlides(dir string) *zero.One {
	prefix := t.AddPath(dir)
	img := t.Img("", "")
	css := t.CSS(t.SlidesCSS())
	js := t.JS(fmt.Sprintf(`
(function() {
    const { frame, state } = pathless.ctx();
    let slides = [], index = state.nav || 0;

    const show = async (i) => {
        if (!slides.length) return;
        index = ((i %% slides.length) + slides.length) %% slides.length;
        pathless.update('nav', index);
        
        const img = frame.querySelector('img');
        if (!img) return;
        
        try {
            const { data } = await pathless.fetch(apiUrl + '/%s/' + slides[index], '%s.' + slides[index]);
            img.src = data;
            img.alt = slides[index];
        } catch { img.alt = 'Failed'; }
    };

    pathless.fetch(apiUrl + '/%s/order', '%s.order')
        .then(({ data }) => { slides = data || []; show(index); });

    pathless.onKey(k => {
        if (k === 'a') show(index - 1);
        else if (k === 'd') show(index + 1);
    });
})();
    `, prefix, prefix, prefix, prefix))

	return t.Build("slides", true, img, &css, &js)
}
