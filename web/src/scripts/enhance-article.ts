import { codeToHtml } from 'shiki';

const languages = new Set([
  'python', 'ruby', 'bash', 'html', 'xml', 'css', 'javascript', 'typescript',
  'json', 'dockerfile', 'yaml', 'java', 'tsx'
]);

async function enhanceArticle() {
  const content = document.getElementById('content');
  const toc = document.querySelector<HTMLElement>('.toc');
  if (!content || !toc) return;

  toc.replaceChildren();
  const headings = Array.from(content.querySelectorAll<HTMLElement>('h1, h2, h3'));
  if (headings.length) {
    const list = document.createElement('ul');
    list.className = 'article-toc';
    for (const heading of headings) {
      if (!heading.id) continue;
      const item = document.createElement('li');
      item.style.marginInlineStart = `${(Number(heading.tagName.slice(1)) - 1) * 0.75}rem`;
      const link = document.createElement('a');
      link.href = `#${heading.id}`;
      link.className = 'link !link-secondary';
      link.textContent = heading.textContent ?? '';
      item.appendChild(link);
      list.appendChild(item);
    }
    toc.appendChild(list);

    let activeLink: HTMLAnchorElement | null = null;
    const observer = new IntersectionObserver((entries) => {
      for (const entry of entries) {
        if (!entry.isIntersecting) continue;
        activeLink?.classList.remove('active');
        activeLink = toc.querySelector<HTMLAnchorElement>(`a[href="#${entry.target.id}"]`);
        activeLink?.classList.add('active');
      }
    }, { rootMargin: '0% 0% -100% 0%', threshold: 0 });
    headings.forEach((heading) => observer.observe(heading));
  }

  for (const code of content.querySelectorAll<HTMLElement>('pre code')) {
    const match = code.className.match(/language-(\w+)/);
    const language = match?.[1] ?? 'text';
    if (match && languages.has(language)) {
      try {
        const highlighted = await codeToHtml(code.textContent ?? '', {
          lang: language,
          theme: 'dracula'
        });
        const temporary = document.createElement('div');
        temporary.innerHTML = highlighted;
        const result = temporary.querySelector('pre code');
        if (result) {
          code.innerHTML = result.innerHTML;
          code.parentElement?.classList.add('shiki');
        }
      } catch (error) {
        console.error(`Failed to highlight ${language} code:`, error);
      }
    }

    const pre = code.closest('pre');
    if (!pre || pre.parentElement?.classList.contains('code-wrapper')) continue;
    const wrapper = document.createElement('div');
    wrapper.className = 'code-wrapper';
    pre.parentElement?.insertBefore(wrapper, pre);
    wrapper.appendChild(pre);
    const button = document.createElement('button');
    button.type = 'button';
    button.className = 'copy-button';
    button.setAttribute('aria-label', 'コードをコピー');
    button.textContent = '⧉';
    button.addEventListener('click', async () => {
      await navigator.clipboard.writeText(code.textContent ?? '');
      button.textContent = '✓';
      window.setTimeout(() => { button.textContent = '⧉'; }, 3000);
    });
    wrapper.appendChild(button);
  }
}

document.addEventListener('astro:page-load', () => { void enhanceArticle(); });
void enhanceArticle();
