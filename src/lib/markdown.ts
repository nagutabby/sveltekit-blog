import { marked } from 'marked';
import { Renderer } from 'marked';
import { gfmHeadingId } from "marked-gfm-heading-id";
import path from 'path';
import { getWebpPath } from '$lib/utils';
import memoize from 'lodash.memoize';

const renderer = new Renderer();

renderer.image = ({ href, title, text }) => {
  if (href && href.startsWith('images/')) {
    href = `/content/articles/images/${href.substring(7)}`;
  }

  const titleAttr = title ? ` title="${title}"` : '';

  return `<img src="${getWebpPath(href)}" alt="${text}"${titleAttr}>`;
};

marked.use({ renderer });
marked.use(gfmHeadingId());

export const convertMarkdownToHtml = memoize((content: string) => marked(content));

export const transformImagePath = (imagePath: string, type: "articles" | "reviews"): string => {
  if (imagePath && typeof imagePath === 'string' && imagePath.startsWith('images/')) {
    const fileName = path.basename(imagePath);
    return `/content/${type}/images/${fileName}`;
  }
  return imagePath;
};
