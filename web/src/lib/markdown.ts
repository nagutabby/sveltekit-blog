import { marked } from 'marked';
import { Renderer } from 'marked';
import { gfmHeadingId } from "marked-gfm-heading-id";
import { getWebpPath } from '$lib/utils';
import memoize from 'lodash.memoize';
import markedKatex from 'marked-katex-extension';

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

const options = {
  nonStandard: true
};

marked.use(markedKatex(options));

export const convertMarkdownToHtml = memoize((content: string) => marked(content));
