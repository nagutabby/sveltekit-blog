import { convertMarkdownToHtml } from '$lib/markdown';
import { error } from '@sveltejs/kit';
import { Code, ConnectError } from '@connectrpc/connect';
import type { Article, Review } from '$lib/types/blog';
import memoize from 'lodash.memoize';
import { dev } from '$app/environment';
import { createBackendClient } from '$lib/server/backendClient';
import { ContentService } from '$lib/gen/blog/content/v1/content_pb';
import type { Article as ArticlePb, Review as ReviewPb } from '$lib/gen/blog/content/v1/content_pb';

const contentClient = createBackendClient(ContentService);

const toDate = (value: string) => (value ? new Date(value) : new Date(0));

const toArticle = (pb: ArticlePb): Article => ({
  id: pb.id,
  body: pb.body,
  title: pb.title,
  image: pb.image,
  publishedAt: toDate(pb.publishedAt),
  updatedAt: toDate(pb.updatedAt)
});

const toReview = (pb: ReviewPb): Review => ({
  id: pb.id,
  body: pb.body,
  title: pb.title,
  description: pb.description,
  jp_e_code: pb.jpECode,
  image: pb.image,
  rating: pb.rating,
  publishedAt: toDate(pb.publishedAt),
  updatedAt: toDate(pb.updatedAt)
});

const getAllRawDataImpl = async (type: "articles" | "reviews"): Promise<(Article | Review)[]> => {
  try {
    if (type === "articles") {
      const response = await contentClient.listArticles({});
      return response.articles.map(toArticle);
    }
    const response = await contentClient.listReviews({});
    return response.reviews.map(toReview);
  } catch {
    throw error(500, 'コンテンツの取得に失敗しました');
  }
};

export const getAllRawData = dev ? getAllRawDataImpl : memoize(getAllRawDataImpl);

export const getAllHTMLData = async (type: "articles" | "reviews") => {
  const allData = await getAllRawData(type);
  await Promise.all(
    allData.map(async data => {
      data.body = await convertMarkdownToHtml(data.body);
    })
  );

  return allData;
};

export const getHTMLData = async (id: string, type: "articles" | "reviews"): Promise<Article | Review> => {
  try {
    if (type === "articles") {
      const response = await contentClient.getArticle({ id });
      const article = toArticle(response.article!);
      article.body = await convertMarkdownToHtml(article.body);
      return article;
    }

    const response = await contentClient.getReview({ id });
    const review = toReview(response.review!);
    review.body = await convertMarkdownToHtml(review.body);
    return review;
  } catch (err) {
    if (err instanceof ConnectError && err.code === Code.NotFound) {
      throw error(404, `記事が見つかりません: ${id}`);
    }
    throw error(500, 'コンテンツの取得に失敗しました');
  }
};
