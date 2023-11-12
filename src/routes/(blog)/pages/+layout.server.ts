import type { LayoutServerLoad } from './$types'

export const load: LayoutServerLoad = async ({ parent }) => {
  const { titles } = await parent();
  const blogData: Blog = {
    titles: titles?.concat(["ページの一覧"])
  }
  return blogData;
}
