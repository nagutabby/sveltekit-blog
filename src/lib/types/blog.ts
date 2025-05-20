import type { Blog } from "$lib/microcms";

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

export type ReviewFrontMatter = {
  title: string;
  isbn_13: string;
  rating: number;
  publishedAt: Date;
  updatedAt: Date;
};

export type Review = {
  id: string;
  body: string;
  title: string;
  isbn_13: string;
  rating: number;
  publishedAt: Date;
  updatedAt: Date;
};
