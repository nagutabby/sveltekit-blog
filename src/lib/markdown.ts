import { marked } from 'marked';
import { Renderer } from 'marked';
import { gfmHeadingId } from "marked-gfm-heading-id";
import path from 'path';

export const  convertMarkdownToHtml = (content: string) => {
  const renderer = new Renderer();

  renderer.image = ({ href, title, text }) => {
    if (href && href.startsWith('images/')) {
      href = `/content/articles/images/${href.substring(7)}`;
    }

    const titleAttr = title ? ` title="${title}"` : '';

    return `<img src="${href}" alt="${text}"${titleAttr}>`;
  };

  marked.use({ renderer });
  marked.use(gfmHeadingId());

  return marked(content);
}

export const transformImagePath = (imagePath: string): string => {
  if (imagePath && typeof imagePath === 'string' && imagePath.startsWith('images/')) {
    const fileName = path.basename(imagePath);
    return `/content/articles/images/${fileName}`;
  }
  return imagePath;
};
