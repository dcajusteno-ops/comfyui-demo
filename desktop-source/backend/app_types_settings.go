package backend

type Settings struct {
	TrashRetentionDays             int              `json:"trashRetentionDays"`
	RootDir                        string           `json:"rootDir,omitempty"`
	OutputDir                      string           `json:"outputDir,omitempty"`
	OutputConfigured               bool             `json:"outputConfigured,omitempty"`
	PathVersion                    int              `json:"pathVersion,omitempty"`
	ShortcutSettings               ShortcutSettings `json:"shortcutSettings,omitempty"`
	UserProfile                    UserProfile      `json:"userProfile,omitempty"`
	UtilityMenu                    UtilityMenuState `json:"utilityMenu,omitempty"`
	GalleryPerformanceMode         string           `json:"galleryPerformanceMode,omitempty"`
	GalleryInitialBatchSize        int              `json:"galleryInitialBatchSize,omitempty"`
	GalleryPageSize                int              `json:"galleryPageSize,omitempty"`
	GalleryThumbPreferred          bool             `json:"galleryThumbPreferred,omitempty"`
	GalleryBackgroundVariantWarmup bool             `json:"galleryBackgroundVariantWarmup,omitempty"`
	GalleryMetadataLazy            bool             `json:"galleryMetadataLazy,omitempty"`
	AlwaysOnTop                    bool             `json:"alwaysOnTop,omitempty"`
	Wallpaper                      WallpaperConfig  `json:"wallpaper,omitempty"`
}

// WallpaperConfig 描述「背景壁纸」的完整配置。
// 与个人头像一样，图片本体复制到应用数据目录，这里只保存 __wallpaper__/ 开头的相对路径。
type WallpaperConfig struct {
	Enabled      bool    `json:"enabled"`
	ImagePath    string  `json:"imagePath,omitempty"`
	Fit          string  `json:"fit,omitempty"`          // cover | contain | auto
	Dim          float64 `json:"dim,omitempty"`          // 0~1，压暗强度
	SurfaceAlpha float64 `json:"surfaceAlpha,omitempty"` // 0.3~1，面板不透明度
	Blur         bool    `json:"blur,omitempty"`         // 面板毛玻璃
	// Version 是派生字段（图片文件的修改时间），只回给前端做缓存刷新用，不落盘
	Version int64 `json:"version,omitempty"`
}

type GalleryPerformanceSettings struct {
	Mode                    string `json:"mode"`
	InitialBatchSize        int    `json:"initialBatchSize"`
	PageSize                int    `json:"pageSize"`
	ThumbPreferred          bool   `json:"thumbPreferred"`
	BackgroundVariantWarmup bool   `json:"backgroundVariantWarmup"`
	MetadataLazy            bool   `json:"metadataLazy"`
}

type WindowBehaviorSettings struct {
	AlwaysOnTop bool `json:"alwaysOnTop"`
}

type UserProfile struct {
	DisplayName        string `json:"displayName,omitempty"`
	Headline           string `json:"headline,omitempty"`
	Bio                string `json:"bio,omitempty"`
	Location           string `json:"location,omitempty"`
	Website            string `json:"website,omitempty"`
	DailyGoal          int    `json:"dailyGoal,omitempty"`
	PreferredStartPage string `json:"preferredStartPage,omitempty"`
	ImagePath          string `json:"imagePath,omitempty"`
}

type UtilityMenuItem struct {
	ID      string `json:"id"`
	Visible bool   `json:"visible"`
	Order   int    `json:"order,omitempty"`
}

type UtilityMenuState struct {
	Items []UtilityMenuItem `json:"items,omitempty"`
}

type LauncherTool struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Icon string `json:"icon"`
	Args string `json:"args"`
}

type DirectoryBinding struct {
	RootDir       string `json:"rootDir"`
	OutputDir     string `json:"outputDir"`
	OutputRelPath string `json:"outputRelPath"`
	Configured    bool   `json:"configured"`
}
