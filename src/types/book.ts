export type Book = {
  thumbnailUrl: string;
}

export type GoogleBooksResponse = {
  totalItems: number;
  items: Array<{
    volumeInfo: {
      imageLinks: {
        thumbnail: string;
      }
    }
  }>;
}
