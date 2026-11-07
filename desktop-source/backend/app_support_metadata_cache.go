package backend

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

func (a *App) loadImageMetaCache() (ImageMetaCache, error) {
	var cache ImageMetaCache
	data, err := os.ReadFile(a.imageMetaCacheFile())
	if err != nil {
		if os.IsNotExist(err) {
			return ImageMetaCache{}, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, &cache); err != nil {
		return ImageMetaCache{}, nil
	}
	if cache == nil {
		cache = ImageMetaCache{}
	}
	return cache, nil
}

func (a *App) saveImageMetaCache(cache ImageMetaCache) error {
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	a.imageMetaSaveMu.Lock()
	defer a.imageMetaSaveMu.Unlock()

	path := a.imageMetaCacheFile()
	tempPath := path + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		_ = os.Remove(tempPath)
		return err
	}
	if err := os.Rename(tempPath, path); err != nil {
		_ = os.Remove(tempPath)
		return err
	}
	return nil
}

func (a *App) scheduleImageMetaCacheSave() {
	a.imageMetaSaveMu.Lock()
	defer a.imageMetaSaveMu.Unlock()

	if a.imageMetaSaveTimer != nil {
		a.imageMetaSaveTimer.Stop()
	}
	a.imageMetaSaveTimer = time.AfterFunc(600*time.Millisecond, func() {
		if err := a.flushImageMetaCacheSave(); err != nil {
			log.Printf("failed to save image metadata cache: %v", err)
		}
	})
}

func (a *App) flushImageMetaCacheSave() error {
	a.imageMetaSaveMu.Lock()
	if a.imageMetaSaveTimer != nil {
		a.imageMetaSaveTimer.Stop()
		a.imageMetaSaveTimer = nil
	}
	a.imageMetaSaveMu.Unlock()

	return a.saveImageMetaCache(a.snapshotImageMetaCache())
}

func measureDirectoryUsage(root string) (int, int, int64, error) {
	files := 0
	dirs := 0
	var bytesFreed int64

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}
		if d.IsDir() {
			dirs++
			return nil
		}
		info, infoErr := d.Info()
		if infoErr != nil {
			return infoErr
		}
		files++
		bytesFreed += info.Size()
		return nil
	})

	return files, dirs, bytesFreed, err
}

func (a *App) ClearPreviewCache() (CacheClearResult, error) {
	result := CacheClearResult{}
	var firstErr error

	if _, err := os.Stat(a.imageVariantsDir()); err == nil {
		files, dirs, bytesFreed, measureErr := measureDirectoryUsage(a.imageVariantsDir())
		if measureErr != nil {
			firstErr = measureErr
		} else {
			result.DeletedFiles += files
			result.DeletedDirs += dirs
			result.BytesFreed += bytesFreed
		}

		if err := os.RemoveAll(a.imageVariantsDir()); err != nil && firstErr == nil {
			firstErr = err
		}
	} else if err != nil && !os.IsNotExist(err) {
		firstErr = err
	}

	if err := os.MkdirAll(a.previewVariantsDir(), 0755); err != nil && firstErr == nil {
		firstErr = err
	}
	if err := os.MkdirAll(a.thumbVariantsDir(), 0755); err != nil && firstErr == nil {
		firstErr = err
	}

	return result, firstErr
}

func (a *App) ensureImageMetaCacheLoaded() {
	a.imageMetaMu.Lock()
	defer a.imageMetaMu.Unlock()

	if a.imageMetaLoaded {
		return
	}

	cache, err := a.loadImageMetaCache()
	if err != nil {
		log.Printf("failed to load image metadata cache: %v", err)
		cache = ImageMetaCache{}
	}

	a.imageMetaCache = cache
	a.imageMetaLoaded = true
}

func (a *App) snapshotImageMetaCache() ImageMetaCache {
	a.imageMetaMu.RLock()
	defer a.imageMetaMu.RUnlock()

	cache := make(ImageMetaCache, len(a.imageMetaCache))
	for key, value := range a.imageMetaCache {
		cache[key] = value
	}
	return cache
}

func (a *App) getImageMetaCacheEntry(relPath string) (ImageMetaCacheEntry, bool) {
	a.ensureImageMetaCacheLoaded()

	a.imageMetaMu.RLock()
	defer a.imageMetaMu.RUnlock()

	entry, ok := a.imageMetaCache[relPath]
	return entry, ok
}

func (a *App) replaceImageMetaCache(cache ImageMetaCache) {
	a.imageMetaMu.Lock()
	defer a.imageMetaMu.Unlock()

	// 扫描生成的新条目可能尚未扫描元数据；若旧缓存同版本且已扫描，保留旧值，
	// 避免后台预热结果被前台扫描覆盖、同一批文件被反复解析。
	for key, entry := range cache {
		existing, ok := a.imageMetaCache[key]
		if !ok || !existing.MetadataScanned || entry.MetadataScanned {
			continue
		}
		if existing.ModTime == entry.ModTime && existing.Size == entry.Size {
			cache[key] = existing
		}
	}

	a.imageMetaCache = cache
	a.imageMetaLoaded = true
}
