import type { Book, GoogleBooksResponse } from "../types/book";

export const generateDescriptionFromText = (body: string) => {
  const hasHTMLTags = /<[a-z][\s\S]*>/i.test(body);

  const plainText: string = hasHTMLTags ? body.replace(/<[^>]+>/g, "") : body;

  const normalizedText = plainText.replace(/\s+/g, " ").trim();
  let description = normalizedText.slice(0, 100);

  if (normalizedText.length > 100) {
    const lastChar = description.slice(-1);
    if (!/[。、．，!！?？]/.test(lastChar)) {
      description += "…";
    }
  }

  return description;
};

export async function getBook(isbn: string, apiKey: string): Promise<Book> {
  let url = `https://www.googleapis.com/books/v1/volumes?q=isbn:${isbn}&key=${apiKey}`;

  const response = await fetch(url);

  if (!response.ok) {
    throw new Error(`API request failed with status: ${response.status}`);
  }

  const data: GoogleBooksResponse = await response.json();

  const volumeInfo = data.items[0].volumeInfo;

  const book = {
    thumbnailUrl: volumeInfo.imageLinks?.thumbnail,
    infoLink: volumeInfo.infoLink,
    description: volumeInfo.description
  };
  return book;
}
