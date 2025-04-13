export type Book = {
  title: string;
  thumbnailUrl: string;
}

export type GoogleBooksResponse = {
  totalItems: number;
  items: Array<{
    volumeInfo: {
      title: string;
      imageLinks: {
        thumbnail: string;
      }
    }
  }>;
}
