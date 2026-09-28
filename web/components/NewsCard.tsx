import type { NewsItem } from "@/lib/types";

function formatDate(iso: string): string {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return "";
  return d.toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

export function NewsCard({ item }: { item: NewsItem }) {
  return (
    <a
      href={item.url}
      target="_blank"
      rel="noopener noreferrer"
      className="flex flex-col gap-2 rounded-lg border border-border bg-surface-muted p-4 transition hover:border-accent"
    >
      <div className="flex items-center gap-2 text-xs text-muted">
        <span className="font-medium text-content">{item.source}</span>
        <span aria-hidden>·</span>
        <time dateTime={item.published_at}>{formatDate(item.published_at)}</time>
      </div>

      <h3 className="line-clamp-2 font-semibold leading-snug">{item.title}</h3>

      <p className="line-clamp-3 text-sm text-muted">
        {item.summary || "Summary coming soon."}
      </p>

      <div className="mt-auto flex flex-wrap gap-1 pt-1">
        {item.topics.map((topic) => (
          <span
            key={topic}
            className="rounded-full bg-accent/10 px-2 py-0.5 text-xs text-accent"
          >
            {topic}
          </span>
        ))}
      </div>
    </a>
  );
}
