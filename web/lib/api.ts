import type { NewsItem, PaginatedNews } from "./types";

const API_URL = process.env.API_URL ?? "http://localhost:8080";

export interface NewsQuery {
  topics?: string[];
  sources?: string[];
  q?: string;
  sort?: string;
  order?: string;
  page?: number;
  pageSize?: number;
}

export async function fetchNews(q: NewsQuery): Promise<PaginatedNews> {
  const params = new URLSearchParams();
  if (q.topics?.length) params.set("topics", q.topics.join(","));
  if (q.sources?.length) params.set("sources", q.sources.join(","));
  if (q.q) params.set("q", q.q);
  if (q.sort) params.set("sort", q.sort);
  if (q.order) params.set("order", q.order);

  const pageSize = q.pageSize ?? 20;
  const page = q.page ?? 1;
  params.set("limit", String(pageSize));
  params.set("offset", String((page - 1) * pageSize));

  const res = await fetch(`${API_URL}/api/v1/news?${params.toString()}`, {
    cache: "no-store",
  });
  if (!res.ok) {
    throw new Error(`API error ${res.status}`);
  }
  return res.json();
}

// Mock fallback so the dashboard renders even before the API/ingestion is up.
export const mockNews: NewsItem[] = [
  {
    id: 1,
    url: "https://example.com/1",
    title: "LLMs are changing how we write software",
    summary:
      "Large language models are increasingly used to generate, review, and maintain code across the industry.",
    source: "Hacker News",
    topics: ["AI/ML", "Software Engineering"],
    published_at: new Date().toISOString(),
    language: "en",
  },
  {
    id: 2,
    url: "https://example.com/2",
    title: "A practical guide to zero-downtime migrations",
    summary:
      "Rolling out schema changes without locking your tables or taking your service offline.",
    source: "InfoQ",
    topics: ["DevOps/SRE"],
    published_at: new Date().toISOString(),
    language: "en",
  },
  {
    id: 3,
    url: "https://example.com/3",
    title: "Supply-chain attacks are on the rise",
    summary:
      "Attackers are shifting focus to the software supply chain, targeting build systems and dependencies.",
    source: "Krebs on Security",
    topics: ["Security"],
    published_at: new Date().toISOString(),
    language: "en",
  },
];
