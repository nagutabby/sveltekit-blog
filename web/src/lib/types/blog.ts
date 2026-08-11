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
  description: string;
  jp_e_code: string;
  image: string;
  rating: number;
  publishedAt: Date;
  updatedAt: Date;
};

export type Review = {
  id: string;
  body: string;
  title: string;
  description: string;
  jp_e_code: string;
  image: string;
  rating: number;
  publishedAt: Date;
  updatedAt: Date;
};
