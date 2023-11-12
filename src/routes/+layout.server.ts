import type { LayoutServerLoad } from './$types'

export const load: LayoutServerLoad = async () => {
  const blogData: Blog = {
    titles: ["Home"]
  }
  return blogData;
}
