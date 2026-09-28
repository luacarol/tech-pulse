"use client";

import type { Lang } from "@/lib/types";

const options: { value: Lang; label: string }[] = [
  { value: "pt", label: "PT" },
  { value: "en", label: "EN" },
  { value: "es", label: "ES" },
];

export function LanguageSwitcher({
  value,
  onChange,
}: {
  value: Lang;
  onChange: (lang: Lang) => void;
}) {
  return (
    <div className="flex rounded-md border border-border">
      {options.map((o) => (
        <button
          key={o.value}
          onClick={() => onChange(o.value)}
          className={`px-2 py-1 text-sm ${
            value === o.value
              ? "bg-accent text-white"
              : "text-muted hover:text-content"
          }`}
        >
          {o.label}
        </button>
      ))}
    </div>
  );
}
