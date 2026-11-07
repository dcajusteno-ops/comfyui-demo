package backend

const customPromptSource = "我的词库"

const defaultFavoriteGroupID = "default"
const defaultFavoriteGroupName = "Default Favorites"

const profileAssetPrefix = "__profile__/"
const variantAssetPrefix = "__variant__/"
const thumbVariantAssetPrefix = variantAssetPrefix + "thumb/"
const previewVariantAssetPrefix = variantAssetPrefix + "preview/"
const thumbVariantMaxDimension = 640
const previewVariantMaxDimension = 1600
const trashAssetPrefix = "__trash__/"
const wallpaperAssetPrefix = "__wallpaper__/"

// 壁纸相关默认值与取值范围（对齐源项目 comfyui-xyz-demo 的 WallpaperConfig）
const wallpaperFileStem = "wallpaper"
const defaultWallpaperFit = "cover"
const defaultWallpaperDim = 0.35
const defaultWallpaperSurfaceAlpha = 0.72
const minWallpaperSurfaceAlpha = 0.3
const maxWallpaperSurfaceAlpha = 1.0
const minWallpaperDim = 0.0
const maxWallpaperDim = 1.0
const maxWallpaperImageBytes = 32 * 1024 * 1024

const pathVersionRootRelative = 2
const customRootsVersion = 3
const builtinDateArchiveRootID = "builtin-date-archive"

const mainWindowWidth = 1440
const mainWindowHeight = 900
const mainWindowMinWidth = 900
const mainWindowMinHeight = 600
const compactWindowMinWidth = 320
const compactWindowMinHeight = 260
