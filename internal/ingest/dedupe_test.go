package ingest

import "testing"

func TestHashURL_StripsUTMParams(t *testing.T) {
	a := HashURL("https://example.com/post?utm_source=x&id=1")
	b := HashURL("https://example.com/post?id=1")
	if a != b {
		t.Fatalf("expected utm-stripped URLs to hash equal:\n  %s\n  %s", a, b)
	}
}

func TestHashURL_DifferentPathsDiffer(t *testing.T) {
	if HashURL("https://example.com/a") == HashURL("https://example.com/b") {
		t.Fatal("expected different paths to hash differently")
	}
}

func TestStripHTML(t *testing.T) {
	got := stripHTML("<p>Hello <b>world</b></p>")
	want := "Hello world"
	if got != want {
		t.Fatalf("stripHTML = %q, want %q", got, want)
	}
}

func TestTruncate(t *testing.T) {
	got := truncate("abcdefghij", 5)
	want := "abcde…"
	if got != want {
		t.Fatalf("truncate = %q, want %q", got, want)
	}
}
