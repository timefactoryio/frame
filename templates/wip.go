package templates

import (
	"html/template"

	"github.com/timefactoryio/frame/zero"
)

func (t *templates) Keyboard() {
	css := t.CSS(t.KeyboardCSS())
	js := t.JS(`
(function(){
  const keys = [
    ['Tab', '', ''],
    ['1', '2', '3'],
    ['q', 'w', 'e'],
    ['a', 's', 'd']
  ];

  const grid = pathless.space.querySelector('.grid');
  if (!grid) return;

  keys.flat().forEach((k) => {
    const entry = pathless.keyboard().find(x => x.key === k);
    const keyEl = document.createElement('div');
    keyEl.className = 'key';
    keyEl.dataset.key = k;
    keyEl.textContent = k.toUpperCase();
    if (entry && entry.style) keyEl.style.cssText = entry.style;
    if (entry && entry.pressed) keyEl.classList.add('pressed');
    grid.appendChild(keyEl);
  });

  document.addEventListener('keydown', (e) => {
    const entry = pathless.keyboard().find(x => x.key === e.key);
    if (entry) {
      const keyEl = grid.querySelector('[data-key="' + e.key + '"]');
      if (keyEl) keyEl.classList.add('pressed');
    }
  });

  document.addEventListener('keyup', (e) => {
    const entry = pathless.keyboard().find(x => x.key === e.key);
    if (entry) {
      const keyEl = grid.querySelector('[data-key="' + e.key + '"]');
      if (keyEl) keyEl.classList.remove('pressed');
    }
  });
})();
`)
	html := zero.One(template.HTML(`<div class="grid"></div>`))
	t.Build("keyboard", true, &html, &css, &js)
}
