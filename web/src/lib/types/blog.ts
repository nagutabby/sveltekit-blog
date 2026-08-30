export type ArticleFrontMatter = {
  id: string;
  title: string;
  image: string;
};

export type Article = {
  id: string;
  body: string;
  title: string;
  image: string;
  publishedAt: Date;
};

export type ReviewFrontMatter = {
  id: string;
  title: string;
  description: string;
  jp_e_code: string;
  image: string;
  rating: number;
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
};
