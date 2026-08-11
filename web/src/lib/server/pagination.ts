export interface Pagination {
  currentPage: number;
  totalPages: number;
  hasNextPage: boolean;
  hasPrevPage: boolean;
}

export interface PaginatedResult<T> extends Pagination {
  items: T[];
}

export function totalPagesFor(itemCount: number, perPage: number): number {
  return Math.ceil(itemCount / perPage);
}

export function paginate<T>(items: T[], page: number, perPage: number): PaginatedResult<T> {
  const totalPages = totalPagesFor(items.length, perPage);
  const startIndex = (page - 1) * perPage;
  const endIndex = startIndex + perPage;

  return {
    items: items.slice(startIndex, endIndex),
    currentPage: page,
    totalPages,
    hasNextPage: page < totalPages,
    hasPrevPage: page > 1
  };
}
