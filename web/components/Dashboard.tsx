"use client";

import { useEffect, useState } from "react";
import { fetchNews, mockNews } from "@/lib/api";
import { t, topicsList } from "@/lib/i18n";
import type { Lang, NewsItem } from "@/lib/types";
import { LanguageSwitcher } from "./LanguageSwitcher";
import { NewsCard } from "./NewsCard";
import { ThemeProvider, useTheme } from "./ThemeProvider";
import { ThemeToggle } from "./ThemeToggle";

function DashboardInner() {
  const [lang, setLang] = useState<Lang>("pt");
  const [items, setItems] = useState<NewsItem[]>(mockNews);
  const [selectedTopics, setSelectedTopics] = useState<string[]>([]);
  const [selectedSources, setSelectedSources] = useState<string[]>([]);
  const [query, setQuery] = useState("");
  const [sort, setSort] = useState<"date" | "source">("date");
  const [order, setOrder] = useState<"desc" | "asc">("desc");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(false);

  const sources = Array.from(new Set(items.map((i) => i.source))).sort();

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setError(false);
    fetchNews({
      topics: selectedTopics,
      sources: selectedSources,
      q: query,
      sort,
      order,
    })
      .then((res) => {
        if (!cancelled) setItems(res.items);
      })
      .catch(() => {
        // Keep mock data visible when the API isn't reachable yet.
        if (!cancelled) setError(true);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [selectedTopics, selectedSources, query, sort, order]);

  const toggleTopic = (topic: string) =>
    setSelectedTopics((prev) =>
      prev.includes(topic)
        ? prev.filter((t) => t !== topic)
        : [...prev, topic]
    );

  const toggleSource = (source: string) =>
    setSelectedSources((prev) =>
      prev.includes(source)
        ? prev.filter((s) => s !== source)
        : [...prev, source]
    );

  return (
    <div className="mx-auto max-w-6xl px-4 py-8">
      <header className="flex items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold">{t(lang, "title")}</h1>
          <p className="text-sm text-muted">{t(lang, "subtitle")}</p>
        </div>
        <div className="flex items-center gap-2">
          <LanguageSwitcher value={lang} onChange={setLang} />
          <ThemeToggle />
        </div>
      </header>

      <div className="mt-6 flex flex-col gap-4">
        <input
          type="search"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder={t(lang, "searchPlaceholder")}
          className="w-full rounded-md border border-border bg-surface-muted px-3 py-2 text-sm outline-none focus:border-accent"
        />

        <div className="flex flex-wrap gap-2">
          <button
            onClick={() => setSelectedTopics([])}
            className={`rounded-full px-3 py-1 text-sm ${
              selectedTopics.length === 0
                ? "bg-accent text-white"
                : "border border-border text-muted hover:text-content"
            }`}
          >
            {t(lang, "allTopics")}
          </button>
          {topicsList.map((topic) => (
            <button
              key={topic}
              onClick={() => toggleTopic(topic)}
              className={`rounded-full px-3 py-1 text-sm ${
                selectedTopics.includes(topic)
                  ? "bg-accent text-white"
                  : "border border-border text-muted hover:text-content"
              }`}
            >
              {topic}
            </button>
          ))}
        </div>

        {sources.length > 0 && (
          <div className="flex flex-wrap gap-2">
            <select
              value={selectedSources[0] ?? ""}
              onChange={(e) =>
                setSelectedSources(e.target.value ? [e.target.value] : [])
              }
              className="rounded-md border border-border bg-surface-muted px-3 py-1 text-sm"
            >
              <option value="">{t(lang, "allSources")}</option>
              {sources.map((s) => (
                <option key={s} value={s}>
                  {s}
                </option>
              ))}
            </select>
          </div>
        )}

        <div className="flex items-center gap-2">
          <select
            value={sort}
            onChange={(e) => setSort(e.target.value as "date" | "source")}
            className="rounded-md border border-border bg-surface-muted px-3 py-1 text-sm"
          >
            <option value="date">{t(lang, "newest")}</option>
            <option value="source">{t(lang, "sortBy")}</option>
          </select>
          <button
            onClick={() => setOrder((o) => (o === "desc" ? "asc" : "desc"))}
            className="rounded-md border border-border px-3 py-1 text-sm text-muted hover:text-content"
          >
            {order === "desc" ? "↓" : "↑"}
          </button>
        </div>
      </div>

      {error && (
        <p className="mt-6 text-sm text-muted">
          {t(lang, "noResults")}
        </p>
      )}

      <main className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {loading && items.length === 0 ? (
          <p className="text-sm text-muted">{t(lang, "loading")}</p>
        ) : (
          items.map((item) => <NewsCard key={item.id} item={item} />)
        )}
      </main>

      {!loading && items.length === 0 && !error && (
        <p className="mt-6 text-sm text-muted">{t(lang, "noResults")}</p>
      )}
    </div>
  );
}

export function Dashboard() {
  return (
    <ThemeProvider>
      <DashboardInner />
    </ThemeProvider>
  );
}
