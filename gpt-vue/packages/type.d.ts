export interface BaseResponse<T> {
  code: number;
  data?: T;
  message?: string;
}

export interface ListResponse<T = Record<string, unknown>> {
  items: T[];
  page: number;
  page_size: number;
  total: number;
  total_page: number
}

