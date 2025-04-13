export type Book = {
  title: string;
  authors: string[];
  publishedDate: string;
  thumbnailUrl: string;
}

export type GoogleBooksResponse = {
  totalItems: number;
  items: Array<{
    volumeInfo: {
      title: string;
      authors: string[];
      publishedDate: string;
      imageLinks: {
        thumbnail: string;
      }
    }
  }>;
}
