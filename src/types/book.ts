export type Book = {
  description: string;
  thumbnailUrl: string;
  infoLink: string;
};

export type GoogleBooksResponse = {
  totalItems: number;
  items: Array<{
    volumeInfo: {
      description: string;
      imageLinks: {
        thumbnail: string;
      };
      infoLink: string;
    };
  }>;
};
