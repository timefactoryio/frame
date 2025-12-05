package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/timefactoryio/frame/zero"
)

func (t *templates) Landing(heading, github, x string) {
	logo := t.Api() + "/img/logo"
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
	js := `
(function(){
  const frame = pathless.frame();
  const key = 'scroll';
  
  frame.scrollTop = pathless.state()[key] || 0;
  frame.addEventListener('scroll', () => {
    pathless.update(key, frame.scrollTop);
  });
  
  let speed = 0;
  let isScrolling = false;
  
  const scroll = () => {
    if (speed === 0) {
      isScrolling = false;
      return;
    }
    frame.scrollBy({ top: speed });
    requestAnimationFrame(scroll);
  };
  
  const speeds = { w: -20, s: 20, a: -40, d: 40 };
  pathless.keybind((k) => {
    if (speeds[k]) {
      speed = speeds[k];
      if (!isScrolling) {
        isScrolling = true;
        scroll();
      }
    }
  });
  
  document.addEventListener('keyup', (e) => {
    if (speeds[e.key]) speed = 0;
  });
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
    const frame = pathless.frame();
    let slides = [];
    let index = pathless.state().nav || 0;

    async function show(i) {
        if (!slides.length) return;
        index = ((i %% slides.length) + slides.length) %% slides.length;
        pathless.update("nav", index);

        const imgEl = frame.querySelector('img');
        if (!imgEl) return;

        const slide = slides[index];
        const fetchKey = '%s.' + slide;
        try {
            const { data } = await pathless.fetch(apiUrl + '/%s/' + slide, { key: fetchKey });
            imgEl.src = data;
            imgEl.alt = slide;
        } catch (e) {
            imgEl.alt = "Failed to load image";
        }
    }

    pathless.fetch(apiUrl + '/%s/order', { key: '%s.order' })
        .then(({ data }) => {
            slides = data || [];
            if (slides.length) show(index);
        });

    pathless.keybind((k) => {
        k = k.toLowerCase();
        if (k === 'a') show(index - 1);
        else if (k === 'd') show(index + 1);
    });
})();
    `, prefix, prefix, prefix, prefix))
	return t.Build("slides", true, img, &css, &js)
}

func (t *templates) BuildVideo(dir string) *zero.One {
	prefix := t.AddPath(dir)
	video := t.Video("")
	css := t.CSS(t.VideoCSS())
	js := t.JS(fmt.Sprintf(`
(function() {
    const frame = pathless.frame();
    const videoEl = frame.querySelector('video');
    if (!videoEl) return;

    let videos = [];
    let index = pathless.state().nav || 0;

    async function show(i) {
        if (!videos.length) return;
        index = ((i %% videos.length) + videos.length) %% videos.length;
        pathless.update("nav", index);

        const video = videos[index];
        videoEl.src = apiUrl + '/%s/' + video;
        videoEl.load();
    }

    pathless.fetch(apiUrl + '/%s/order', { key: '%s.order' })
        .then(({ data }) => {
            videos = data || [];
            if (videos.length) show(index);
        });

    document.addEventListener('keydown', (e) => {
        if (e.key === ' ') {
            e.preventDefault();
            if (videoEl.paused) {
                videoEl.play().catch(() => {});
            } else {
                videoEl.pause();
            }
        }
    });

    pathless.keybind((k) => {
        k = k.toLowerCase();
        if (k === 'a') {
            show(index - 1);
        } else if (k === 'd') {
            show(index + 1);
        } else if (k === 'w') {
            videoEl.playbackRate = Math.min(videoEl.playbackRate + 0.25, 4);
        } else if (k === 's') {
            videoEl.playbackRate = Math.max(videoEl.playbackRate - 0.25, 0.25);
        } else if (k === 'x') {
            videoEl.volume = Math.min(videoEl.volume + 0.1, 1);
        } else if (k === 'c') {
            videoEl.volume = Math.max(videoEl.volume - 0.1, 0);
        }
    });
})();
    `, prefix, prefix, prefix))
	return t.Build("video", true, video, &css, &js)
}
