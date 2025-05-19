import { marked } from 'marked';
import { Renderer } from 'marked';

export function convertMarkdownToHtml(content: string) {
  const renderer = new Renderer();

  renderer.image = ({ href, title, text }) => {
    if (href && href.startsWith('images/')) {
      href = `/content/articles/images/${href.substring(7)}`;
    }

    const titleAttr = title ? ` title="${title}"` : '';

    return `<img src="${href}" alt="${text}"${titleAttr}>`;
  };

  const markedOptions = {
    renderer,
    headerIds: true,
    gfm: true,
    breaks: false,
    pedantic: false,
    smartLists: true,
    smartypants: true
  };

  return marked(content, markedOptions);
}
