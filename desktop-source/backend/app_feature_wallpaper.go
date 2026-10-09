package backend

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 背景壁纸沿用个人头像那套做法：
// 用户选图 → 复制到应用数据目录（data/wallpaper/wallpaper.<ext>）→ 设置里只存 __wallpaper__/ 相对路径，
// 图片本体由 AssetServer 的 ServeImage 通过 __wallpaper__/ 前缀提供。
// 这样原图被移动/删除都不影响壁纸，也不需要把图片塞进 settings.json。

func (a *App) GetWallpaper() (WallpaperConfig, error) {
	settings, err := a.loadSettings()
	if err != nil {
		return defaultWallpaperConfig(), err
	}
	return a.wallpaperForClient(settings.Wallpaper), nil
}

// SaveWallpaper 只保存外观参数；图片仍由 SelectWallpaperImage / ClearWallpaperImage 管理。
func (a *App) SaveWallpaper(config WallpaperConfig) (WallpaperConfig, error) {
	settings, err := a.loadSettings()
	if err != nil {
		return defaultWallpaperConfig(), err
	}

	next := config
	// 客户端传来的版本号属于派生信息，一律以实际文件为准
	next.Version = 0
	next.ImagePath = settings.Wallpaper.ImagePath
	settings.Wallpaper = normalizeWallpaperConfig(next)

	if err := a.saveSettings(settings); err != nil {
		return settings.Wallpaper, err
	}

	return a.wallpaperForClient(settings.Wallpaper), nil
}

func (a *App) SelectWallpaperImage() (WallpaperConfig, error) {
	options := runtime.OpenDialogOptions{
		Title: "选择一张壁纸图片",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Image Files (*.png;*.jpg;*.jpeg;*.webp;*.gif)",
				Pattern:     "*.png;*.jpg;*.jpeg;*.webp;*.gif",
			},
		},
	}

	filePath, err := runtime.OpenFileDialog(a.ctx, options)
	if err != nil {
		return defaultWallpaperConfig(), err
	}

	// 用户取消选择：原样返回，不算错误
	if strings.TrimSpace(filePath) == "" {
		return a.GetWallpaper()
	}

	return a.saveWallpaperImage(filePath)
}

func (a *App) ClearWallpaperImage() (WallpaperConfig, error) {
	settings, err := a.loadSettings()
	if err != nil {
		return defaultWallpaperConfig(), err
	}

	if err := a.removeWallpaperFiles(); err != nil {
		return a.wallpaperForClient(settings.Wallpaper), err
	}

	settings.Wallpaper.ImagePath = ""
	settings.Wallpaper.Enabled = false
	settings.Wallpaper.Version = 0
	settings.Wallpaper = normalizeWallpaperConfig(settings.Wallpaper)

	if err := a.saveSettings(settings); err != nil {
		return settings.Wallpaper, err
	}

	return a.wallpaperForClient(settings.Wallpaper), nil
}

func (a *App) saveWallpaperImage(sourcePath string) (WallpaperConfig, error) {
	settings, err := a.loadSettings()
	if err != nil {
		return defaultWallpaperConfig(), err
	}

	ext := strings.ToLower(filepath.Ext(sourcePath))
	if !isSupportedProfileImageExt(ext) {
		return a.wallpaperForClient(settings.Wallpaper), fmt.Errorf("unsupported image format")
	}

	info, err := os.Stat(sourcePath)
	if err != nil {
		return a.wallpaperForClient(settings.Wallpaper), err
	}
	if info.Size() > maxWallpaperImageBytes {
		return a.wallpaperForClient(settings.Wallpaper),
			fmt.Errorf("壁纸图片过大（上限 %d MB）", maxWallpaperImageBytes/1024/1024)
	}

	if err := os.MkdirAll(a.wallpaperImageDir(), 0755); err != nil {
		return a.wallpaperForClient(settings.Wallpaper), err
	}

	// 只保留一张：先清掉旧的，避免换了格式后残留
	if err := a.removeWallpaperFiles(); err != nil {
		return a.wallpaperForClient(settings.Wallpaper), err
	}

	targetName := wallpaperFileStem + ext
	targetPath := filepath.Join(a.wallpaperImageDir(), targetName)
	if err := copyFile(sourcePath, targetPath); err != nil {
		return a.wallpaperForClient(settings.Wallpaper), err
	}

	settings.Wallpaper.ImagePath = normalizeRelPath(wallpaperAssetPrefix + targetName)
	// 选完图直接启用：用户的意图就是用这张图
	settings.Wallpaper.Enabled = true
	settings.Wallpaper.Version = 0
	settings.Wallpaper = normalizeWallpaperConfig(settings.Wallpaper)

	if err := a.saveSettings(settings); err != nil {
		return settings.Wallpaper, err
	}

	return a.wallpaperForClient(settings.Wallpaper), nil
}

func (a *App) removeWallpaperFiles() error {
	existing, err := filepath.Glob(filepath.Join(a.wallpaperImageDir(), wallpaperFileStem+".*"))
	if err != nil {
		return err
	}
	for _, item := range existing {
		if removeErr := os.Remove(item); removeErr != nil && !os.IsNotExist(removeErr) {
			return removeErr
		}
	}
	return nil
}

// wallpaperForClient 补上文件版本号。
// 壁纸文件名固定（wallpaper.png），换图后 URL 不变会被浏览器缓存住，
// 所以把文件修改时间作为版本号带给前端拼在 URL 后面，强制刷新。
func (a *App) wallpaperForClient(config WallpaperConfig) WallpaperConfig {
	config.Version = 0
	if config.ImagePath == "" {
		return config
	}

	absPath, err := a.resolveWallpaperAssetPath(config.ImagePath)
	if err != nil {
		return config
	}

	info, err := os.Stat(absPath)
	if err != nil {
		// 文件不在了：当作没有壁纸，避免前端一直请求 404
		config.ImagePath = ""
		config.Enabled = false
		return config
	}

	config.Version = info.ModTime().UnixNano()
	return config
}
