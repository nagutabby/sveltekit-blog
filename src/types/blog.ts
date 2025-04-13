import type { Blog } from "$lib/microcms";
import type { MicroCMSImage } from "microcms-js-sdk";

export type ArticleInputData = Pick<Blog, "title" | "image" | "body">
export type ReviewInputData = {
  image: MicroCMSImage;
  title: string;
  body: string;
}
