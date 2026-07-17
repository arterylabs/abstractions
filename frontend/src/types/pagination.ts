export interface PaginationParams {
  page?: number;
  perPage?: number;
  query?: string;
}

export interface SortableParams extends PaginationParams {
  sort?: string;
  order?: "asc" | "desc";
}

export interface PageInfo {
  page: number;
  perPage: number;
  totalPages: number;
  hasNextPage: boolean;
  totalRecords: number;
  pageRecords: number;
}

export interface PaginatedResponse<T> {
  results: T[];
  pageInfo: PageInfo;
}
