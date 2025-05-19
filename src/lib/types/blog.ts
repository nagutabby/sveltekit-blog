import type { Blog } from "$lib/microcms";
import type { MicroCMSImage } from "microcms-js-sdk";

export type ArticleFrontMatter = {
  title: string;
  image: string;
  publishedAt: Date;
  updatedAt: Date;
};

export type Article = {
  id: string;
  body: string;
  title: string;
  image: string;
  publishedAt: Date;
  updatedAt: Date;
};

export type ArticleInputData = Pick<Blog, "title" | "image" | "body">;

export type ReviewInputData = {
  image: MicroCMSImage;
  title: string;
  body: string;
};
