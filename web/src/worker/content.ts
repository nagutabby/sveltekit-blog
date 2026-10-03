import articles from './articles.generated.json';

export interface Article {
  id: string;
  title: string;
  publishedAt: Date;
}

const publishedArticles: Article[] = articles.map((article) => ({
  ...article,
  publishedAt: new Date(article.publishedAt)
}));

export function listArticles(): Article[] {
  return publishedArticles;
}

export function getArticle(id: string): Article | undefined {
  return publishedArticles.find((article) => article.id === id);
}
