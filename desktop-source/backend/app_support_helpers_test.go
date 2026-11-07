package backend

import (
	"testing"
	"time"
)

func TestNormalizeRelPath(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{input: "", want: ""},
		{input: "  ", want: ""},
		{input: ".", want: ""},
		{input: "a/b.png", want: "a/b.png"},
		{input: "./a/b.png", want: "a/b.png"},
		{input: "/a/b.png/", want: "a/b.png"},
		{input: `a\b\c.png`, want: "a/b/c.png"},
	}

	for _, tc := range cases {
		if got := normalizeRelPath(tc.input); got != tc.want {
			t.Fatalf("normalizeRelPath(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestExtractDateKeyFromRelPath(t *testing.T) {
	modTime := time.Date(2025, 1, 2, 3, 4, 5, 0, time.Local)

	if got := extractDateKeyFromRelPath("2024/2024-05-06/a.png", modTime); got != "2024-05-06" {
		t.Fatalf("expected folder date, got %q", got)
	}
	if got := extractDateKeyFromRelPath("plain/a.png", modTime); got != "2025-01-02" {
		t.Fatalf("expected modtime fallback, got %q", got)
	}
	if got := extractDateKeyFromRelPath("plain/a.png", time.Time{}); got != "" {
		t.Fatalf("expected empty date for zero modtime, got %q", got)
	}
}

func TestMatchesDatePreset(t *testing.T) {
	now := time.Now().In(time.Local)
	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	tenDaysAgo := now.AddDate(0, 0, -10).Format("2006-01-02")

	cases := []struct {
		name    string
		dateKey string
		preset  string
		start   string
		end     string
		want    bool
	}{
		{name: "all", dateKey: today, preset: "all", want: true},
		{name: "today matches", dateKey: today, preset: "today", want: true},
		{name: "today rejects yesterday", dateKey: yesterday, preset: "today", want: false},
		{name: "yesterday matches", dateKey: yesterday, preset: "yesterday", want: true},
		{name: "last7 includes today", dateKey: today, preset: "last7", want: true},
		{name: "last7 rejects ten days ago", dateKey: tenDaysAgo, preset: "last7", want: false},
		{name: "month includes today", dateKey: today, preset: "month", want: true},
		{name: "custom includes", dateKey: today, preset: "custom", start: today, end: today, want: true},
		{name: "custom rejects before start", dateKey: tenDaysAgo, preset: "custom", start: yesterday, want: false},
		{name: "empty date never matches", dateKey: "", preset: "all", want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := matchesDatePreset(tc.dateKey, tc.preset, tc.start, tc.end); got != tc.want {
				t.Fatalf("matchesDatePreset(%q, %q, %q, %q) = %v, want %v",
					tc.dateKey, tc.preset, tc.start, tc.end, got, tc.want)
			}
		})
	}
}

func TestSanitizeJSONSpecialNumbers(t *testing.T) {
	input := `{"a": NaN, "b": "NaN", "c": Infinity, "d": -Infinity}`
	want := `{"a": null, "b": "NaN", "c": null, "d": null}`

	got, changed := sanitizeJSONSpecialNumbers(input)
	if !changed {
		t.Fatalf("expected changed=true")
	}
	if got != want {
		t.Fatalf("sanitizeJSONSpecialNumbers() = %q, want %q", got, want)
	}

	safe := `{"a": 1, "b": "Infinity"}`
	got, changed = sanitizeJSONSpecialNumbers(safe)
	if changed || got != safe {
		t.Fatalf("safe input should stay unchanged, got %q changed=%v", got, changed)
	}
}

func TestImageVariantFilenameStable(t *testing.T) {
	first := imageVariantFilename("thumb", "a/b.png")
	second := imageVariantFilename("thumb", "a/b.png")
	if first != second {
		t.Fatalf("variant filename should be stable: %q vs %q", first, second)
	}

	if imageVariantFilename("preview", "a/b.png") == first {
		t.Fatalf("different kinds should not share the same filename")
	}
	if imageVariantFilename("thumb", "a/c.png") == first {
		t.Fatalf("different sources should not share the same filename")
	}
	if imageVariantFilename("thumb", `a\b.png`) != first {
		t.Fatalf("windows separators should normalize to the same filename")
	}
}
