export type Topic =
  | "AI/ML"
  | "Cloud"
  | "Software Engineering"
  | "DevOps/SRE"
  | "Security"
  | "Career/Market";

export interface NewsItem {
  id: number;
  url: string;
  title: string;
  summary: string;
  source: string;
  topics: Topic[];
  published_at: string;
  language: string;
}

export interface PaginatedNews {
  items: NewsItem[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export type Lang = "pt" | "en" | "es";
