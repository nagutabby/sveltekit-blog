import type { MicroCMSQueries } from "microcms-js-sdk";
import type { PageServerLoad } from "./$types";
import type { ReviewInputData } from "$lib/types/blog";
import { getReviewList, type Review } from "$lib/microcms";
import { GOOGLE_CLOUD_API_KEY } from "$env/static/private";
import { getBook } from "$lib/utils";

export const load: PageServerLoad = async ({ url }) => {
  const page = Number(url.searchParams.get("page")) || 1;
  const limit = 10;
  const startIndex = limit * (page - 1);
  let pageQueries: MicroCMSQueries = {
    limit: limit,
    offset: startIndex
  };
  const reviewData = await getReviewList(pageQueries);

  const reviewHomeData: ReviewInputData = {
    image: {
      url: "https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/695170d8996a4099b0c13e031f8dd918/Microsoft-Fluentui-Emoji-3d-Open-Book-3d.1024.png"
    },
    title: "本のレビュー",
    body: "本を読んだ感想を掲載しています"
  };

  const bookDataPromises = reviewData.contents.map(async (review) => {
    const isbn = review.isbn_13;
    const apiKey = GOOGLE_CLOUD_API_KEY;
    return await getBook(isbn, apiKey);
  });

  let bookData = await Promise.all(bookDataPromises);
  bookData = bookData.map(book => {
    if (book.thumbnailUrl.startsWith('http:')) {
      book.thumbnailUrl = book.thumbnailUrl.replace('http:', 'https:');
    }
    return book;
  });

  const mergedReviewData = reviewData.contents.map((review: Review, index: number) => ({
    review: review,
    book: bookData[index]
  }));

  const { contents, ...reviewMetaData } = reviewData;
  const data = {
    contents: mergedReviewData,
    ...reviewHomeData,
    ...reviewMetaData
  };
  return data;
};
