import type { Blog } from "$lib/microcms";

export type BlogInputData = Pick<Blog, "title" | "image" | "body">

