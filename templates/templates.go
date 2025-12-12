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
  const frame = pathless.Frame();
  const key = 'scroll';

  frame.scrollTop = pathless.read()[key] || 0;
  frame.addEventListener('scroll', () => write(key, frame.scrollTop));

  let speed = 0;
  let scrolling = false;

  function scroll() {
    if (speed === 0) {
      scrolling = false;
      return;
    }
    frame.scrollBy({ top: speed });
    requestAnimationFrame(scroll);
  }

  function startScroll(newSpeed) {
    speed = newSpeed;
    if (!scrolling) {
      scrolling = true;
      scroll();
    }
  }

  pathless.onKey('w', () => startScroll(-20), () => { speed = 0; });
  pathless.onKey('s', () => startScroll(20), () => { speed = 0; });
  pathless.onKey('a', () => startScroll(-40), () => { speed = 0; });
  pathless.onKey('d', () => startScroll(40), () => { speed = 0; });
})();
`)
	return &js
}

func (t *templates) BuildSlides(dir string) *zero.One {
	prefix := t.AddPath(dir)
	css := t.CSS(t.SlidesCSS())
	img := t.Img("", "")
	js := t.JS(fmt.Sprintf(`
(function() {
  const frame = pathless.Frame();
  let slides = [];
  let index = pathless.read()[%q] || 0;

  async function show(i) {
    if (!slides.length) return;
    index = ((i %% slides.length) + slides.length) %% slides.length;
    write(%q, index);

    const img = frame.querySelector('img');
    if (!img) return;

    const slide = slides[index];
    pathless.fetch(window.apiUrl + '/%s/' + slide)
      .then(({ data }) => {
        img.src = data;
        img.alt = slide;
      })
      .catch(() => {
        img.alt = "Failed to load image";
      });
  }

  pathless.fetch(window.apiUrl + '/%s')
    .then(({ data }) => {
      slides = data || [];
      if (slides.length) show(index);
    });

  pathless.onKey('a', () => show(index - 1));
  pathless.onKey('d', () => show(index + 1));
})();
`, prefix, prefix, prefix, prefix))
	return t.Build("slides", true, img, &css, &js)
}

func (t *templates) Keyboard(asFrame bool) *zero.One {
	css := t.CSS(t.KeyboardCSS())
	js := t.JS(`
(function(){
  const keys = [
    ['tab', '', ''],
    ['1', '2', '3'],
    ['q', 'w', 'e'],
    ['a', 's', 'd']
  ];

  const space = pathless.Space();
  const grid = space.querySelector('.grid');
  if (!grid) return;

  function render() {
    grid.innerHTML = '';
    const keyboard = pathless.keyboard();
    
    keys.flat().forEach((k) => {
      const entry = keyboard.find(x => x.key === k);
      const keyEl = document.createElement('div');
      keyEl.className = 'key';
      keyEl.dataset.key = k;
      keyEl.textContent = k.toUpperCase();
      if (entry?.pressed) keyEl.classList.add('pressed');
      grid.appendChild(keyEl);
    });
  }

  render();
  window.addEventListener('keyboardchange', render);
})();
`)
	html := zero.One(template.HTML(`<div class="grid"></div>`))
	final := t.Build("keyboard", asFrame, &html, &css, &js)
	return final
}

func (t *templates) BuildVideo(filePath string) *zero.One {
	t.AddFile(filePath, "video")
	name := filepath.Base(filePath)
	name = name[:len(name)-len(filepath.Ext(name))]

	video := t.Video("")
	css := t.CSS(t.VideoCSS())
	js := t.JS(fmt.Sprintf(`
(function() {
  const frame = pathless.Frame();
  const el = frame.querySelector('video');
  if (!el) return;

  const state = pathless.state();
  el.volume = 1;
  el.src = apiUrl + '/video/%s#t=' + (state.t || 0);
  el.load();

  if (!state.paused) el.play().catch(() => {});

  pathless.onKey(' ', () => {
    if (el.paused) {
      el.play().catch(() => {});
      pathless.update('paused', false);
    } else {
      el.pause();
      pathless.update('paused', true);
    }
  });

  el.addEventListener('timeupdate', () => {
    pathless.update('t', el.currentTime || 0);
  });

  // Cleanup: save time and pause when frame is unloaded
  window.addEventListener('beforeunload', () => {
    pathless.update('t', el.currentTime || 0);
    el.pause();
  });
  document.addEventListener('visibilitychange', () => {
    if (document.hidden) {
      pathless.update('t', el.currentTime || 0);
      el.pause();
    }
  });
})();
`, name))
	return t.Build("video", true, video, &css, &js)
}
