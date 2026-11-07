package backend

import (
	"fmt"
	"testing"
)

func testImage(name, relPath, modTime string, size int64, width, height int) ImageFile {
	return ImageFile{
		Name:    name,
		RelPath: relPath,
		ModTime: modTime,
		Size:    size,
		Width:   width,
		Height:  height,
	}
}

func imageRelPaths(images []ImageFile) []string {
	result := make([]string, 0, len(images))
	for _, image := range images {
		result = append(result, image.RelPath)
	}
	return result
}

func assertRelPaths(t *testing.T, images []ImageFile, want []string) {
	t.Helper()
	got := imageRelPaths(images)
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("unexpected order: got %v, want %v", got, want)
	}
}

func TestSortImageFiles(t *testing.T) {
	t.Run("time asc and desc", func(t *testing.T) {
		images := []ImageFile{
			testImage("b.png", "a/b.png", "2026-01-02T00:00:00Z", 1, 10, 10),
			testImage("a.png", "a/a.png", "2026-01-01T00:00:00Z", 2, 10, 10),
			testImage("c.png", "a/c.png", "2026-01-03T00:00:00Z", 3, 10, 10),
		}

		sortImageFiles(images, "time", "asc")
		assertRelPaths(t, images, []string{"a/a.png", "a/b.png", "a/c.png"})

		sortImageFiles(images, "time", "desc")
		assertRelPaths(t, images, []string{"a/c.png", "a/b.png", "a/a.png"})
	})

	t.Run("equal timestamps use relpath tie-breaker", func(t *testing.T) {
		images := []ImageFile{
			testImage("b.png", "a/b.png", "2026-01-01T00:00:00Z", 1, 10, 10),
			testImage("a.png", "a/a.png", "2026-01-01T00:00:00Z", 2, 10, 10),
		}

		sortImageFiles(images, "time", "asc")
		assertRelPaths(t, images, []string{"a/a.png", "a/b.png"})

		sortImageFiles(images, "time", "desc")
		assertRelPaths(t, images, []string{"a/b.png", "a/a.png"})
	})

	t.Run("size and dimensions", func(t *testing.T) {
		images := []ImageFile{
			testImage("small.png", "a/small.png", "2026-01-01T00:00:00Z", 100, 10, 10),
			testImage("large.png", "a/large.png", "2026-01-02T00:00:00Z", 300, 20, 20),
			testImage("medium.png", "a/medium.png", "2026-01-03T00:00:00Z", 200, 10, 20),
		}

		sortImageFiles(images, "size", "desc")
		assertRelPaths(t, images, []string{"a/large.png", "a/medium.png", "a/small.png"})

		sortImageFiles(images, "dimensions", "desc")
		assertRelPaths(t, images, []string{"a/large.png", "a/medium.png", "a/small.png"})
	})

	t.Run("name", func(t *testing.T) {
		images := []ImageFile{
			testImage("c.png", "b/c.png", "2026-01-01T00:00:00Z", 1, 1, 1),
			testImage("a.png", "b/a.png", "2026-01-02T00:00:00Z", 1, 1, 1),
			testImage("b.png", "b/b.png", "2026-01-03T00:00:00Z", 1, 1, 1),
		}

		sortImageFiles(images, "name", "asc")
		assertRelPaths(t, images, []string{"b/a.png", "b/b.png", "b/c.png"})
	})
}

func TestMatchesScopeRelPath(t *testing.T) {
	cases := []struct {
		name    string
		relPath string
		scope   string
		want    bool
	}{
		{name: "empty scope", relPath: "a/b.png", scope: "", want: true},
		{name: "exact path", relPath: "a/b", scope: "a/b", want: true},
		{name: "child path", relPath: "a/b/c.png", scope: "a/b", want: true},
		{name: "prefix boundary", relPath: "a/bc.png", scope: "a/b", want: false},
		{name: "sibling", relPath: "a/c.png", scope: "a/b", want: false},
		{name: "windows separator in source", relPath: `a\b\c.png`, scope: "a/b", want: true},
		{name: "leading slash in scope", relPath: "a/b/c.png", scope: "/a/b/", want: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := matchesScopeRelPath(tc.relPath, tc.scope); got != tc.want {
				t.Fatalf("matchesScopeRelPath(%q, %q) = %v, want %v", tc.relPath, tc.scope, got, tc.want)
			}
		})
	}
}

func TestBuildPagedResultFromItems(t *testing.T) {
	items := []ImageFile{
		testImage("1.png", "a/1.png", "2026-01-01T00:00:00Z", 1, 1, 1),
		testImage("2.png", "a/2.png", "2026-01-02T00:00:00Z", 1, 1, 1),
		testImage("3.png", "a/3.png", "2026-01-03T00:00:00Z", 1, 1, 1),
		testImage("4.png", "a/4.png", "2026-01-04T00:00:00Z", 1, 1, 1),
		testImage("5.png", "a/5.png", "2026-01-05T00:00:00Z", 1, 1, 1),
	}

	app := &App{}

	t.Run("middle page", func(t *testing.T) {
		result := app.buildPagedResultFromItems(items, GetImagesPageQuery{Page: 2, PageSize: 2}, "performance", "reason")
		if result.Total != 5 || result.TotalPages != 3 || !result.HasMore {
			t.Fatalf("unexpected pagination meta: %+v", result)
		}
		assertRelPaths(t, result.Items, []string{"a/3.png", "a/4.png"})
	})

	t.Run("last page", func(t *testing.T) {
		result := app.buildPagedResultFromItems(items, GetImagesPageQuery{Page: 3, PageSize: 2}, "performance", "reason")
		if result.HasMore {
			t.Fatalf("last page should not have more items")
		}
		assertRelPaths(t, result.Items, []string{"a/5.png"})
	})

	t.Run("page overflow clamps", func(t *testing.T) {
		result := app.buildPagedResultFromItems(items, GetImagesPageQuery{Page: 99, PageSize: 2}, "performance", "reason")
		if result.Page != 3 {
			t.Fatalf("page should clamp to last page, got %d", result.Page)
		}
		assertRelPaths(t, result.Items, []string{"a/5.png"})
	})

	t.Run("empty result", func(t *testing.T) {
		result := app.buildPagedResultFromItems(nil, GetImagesPageQuery{Page: 5, PageSize: 50}, "performance", "reason")
		if result.Total != 0 || result.TotalPages != 0 || result.Page != 1 || result.HasMore {
			t.Fatalf("unexpected empty pagination meta: %+v", result)
		}
		if len(result.Items) != 0 {
			t.Fatalf("expected no items, got %d", len(result.Items))
		}
	})
}
