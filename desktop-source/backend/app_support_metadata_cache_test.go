package backend

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestSaveImageMetaCacheWritesCompactAtomicJSON(t *testing.T) {
	app := &App{dataDir: t.TempDir()}
	cache := ImageMetaCache{
		"a/b.png": {
			Name:    "b.png",
			RelPath: "a/b.png",
			ModTime: "2026-01-01T00:00:00Z",
			Size:    1,
		},
	}

	if err := app.saveImageMetaCache(cache); err != nil {
		t.Fatalf("saveImageMetaCache failed: %v", err)
	}

	data, err := os.ReadFile(app.imageMetaCacheFile())
	if err != nil {
		t.Fatalf("failed to read cache file: %v", err)
	}
	if strings.Contains(string(data), "\n  ") {
		t.Fatalf("expected compact JSON without indentation")
	}

	var decoded ImageMetaCache
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("cache file is not valid JSON: %v", err)
	}
	if decoded["a/b.png"].Size != 1 {
		t.Fatalf("unexpected decoded cache: %+v", decoded)
	}

	if _, err := os.Stat(app.imageMetaCacheFile() + ".tmp"); !os.IsNotExist(err) {
		t.Fatalf("temporary cache file should be removed after rename")
	}

	cache["a/b.png"] = ImageMetaCacheEntry{Name: "b.png", RelPath: "a/b.png", ModTime: "2026-01-02T00:00:00Z", Size: 2}
	if err := app.saveImageMetaCache(cache); err != nil {
		t.Fatalf("second save should overwrite atomically: %v", err)
	}

	data, err = os.ReadFile(app.imageMetaCacheFile())
	if err != nil {
		t.Fatalf("failed to read cache file after overwrite: %v", err)
	}
	decoded = ImageMetaCache{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("overwritten cache file is not valid JSON: %v", err)
	}
	if decoded["a/b.png"].Size != 2 {
		t.Fatalf("overwrite did not take effect: %+v", decoded)
	}
}
