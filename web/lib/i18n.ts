import type { Lang } from "./types";

type Dict = Record<string, string>;

export const topicsList = [
  "AI/ML",
  "Cloud",
  "Software Engineering",
  "DevOps/SRE",
  "Security",
  "Career/Market",
] as const;

const dictionaries: Record<Lang, Dict> = {
  en: {
    title: "Tech Pulse",
    subtitle: "Your daily tech news digest",
    searchPlaceholder: "Search news…",
    allTopics: "All topics",
    allSources: "All sources",
    sortBy: "Sort",
    newest: "Newest",
    oldest: "Oldest",
    loading: "Loading…",
    noResults: "No news found. Try adjusting your filters.",
    language: "Language",
    theme: "Theme",
  },
  pt: {
    title: "Tech Pulse",
    subtitle: "Seu resumo diário de notícias de tecnologia",
    searchPlaceholder: "Buscar notícias…",
    allTopics: "Todos os tópicos",
    allSources: "Todas as fontes",
    sortBy: "Ordenar",
    newest: "Mais recentes",
    oldest: "Mais antigas",
    loading: "Carregando…",
    noResults: "Nenhuma notícia encontrada. Ajuste seus filtros.",
    language: "Idioma",
    theme: "Tema",
  },
  es: {
    title: "Tech Pulse",
    subtitle: "Tu resumen diario de noticias de tecnología",
    searchPlaceholder: "Buscar noticias…",
    allTopics: "Todos los temas",
    allSources: "Todas las fuentes",
    sortBy: "Ordenar",
    newest: "Más recientes",
    oldest: "Más antiguas",
    loading: "Cargando…",
    noResults: "No se encontraron noticias. Ajusta tus filtros.",
    language: "Idioma",
    theme: "Tema",
  },
};

export function t(lang: Lang, key: string): string {
  return dictionaries[lang]?.[key] ?? dictionaries.en[key] ?? key;
}
