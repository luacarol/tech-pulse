"use client";

import { useTheme } from "./ThemeProvider";

export function ThemeToggle() {
  const { theme, toggle } = useTheme();
  return (
    <button
      onClick={toggle}
      aria-label="Toggle theme"
      className="rounded-md border border-border px-2 py-1 text-sm text-muted hover:text-content"
    >
      {theme === "dark" ? "☀️" : "🌙"}
    </button>
  );
}
