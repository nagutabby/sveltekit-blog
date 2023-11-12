import { createClient, type MicroCMSQueries, type MicroCMSImage } from "microcms-js-sdk";
import { MICROCMS_SERVICE_DOMAIN, MICROCMS_API_KEY } from '$env/static/private';

const client = createClient({
  serviceDomain: MICROCMS_SERVICE_DOMAIN,
  apiKey: MICROCMS_API_KEY,
});

export type Blog = {
  id: string;
  createdAt: string;
  updatedAt: string;
  publishedAt: string;
  revisedAt: string;
  image: MicroCMSImage;
  title: string;
  description: string;
  body: string;
  attribs: {
    id: string;
  }
  tags: string;
};
export type ArticleResponse = {
  totalCount: number;
  offset: number;
  limit: number;
  contents: Blog[];
};

export const getList = async (queries?: MicroCMSQueries) => {
  return await client.get<ArticleResponse>({ endpoint: "blog", queries });
};
export const getDetail = async (
  contentId: string,
  queries?: MicroCMSQueries
) => {
  return await client.getListDetail<Blog>({
    endpoint: "blog",
    contentId,
    queries,
  });
};

export const getAllContents = async (limit = 10, offset = 0) => {
  const data = await client.get<ArticleResponse>({ endpoint: "blog", queries: { limit, offset } });
  let tmpOffset = 0;
  while (tmpOffset < data.totalCount) {
    tmpOffset += limit;
    const tmpData: ArticleResponse = await getList({ limit: limit, offset: tmpOffset });
    data.contents = data.contents.concat(tmpData.contents);
  }
  return data;
}
