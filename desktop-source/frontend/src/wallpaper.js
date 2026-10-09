import { computed, ref } from 'vue'
import * as App from '@/api'

// ---------------------------------------------------------------------------
// 背景壁纸
//
// 配置存在后端 settings.json（和个人头像同一套模式），这里额外在 localStorage
// 留一份缓存，用于刷新时首帧就把壁纸铺上，避免先闪一下没有壁纸的界面。
// 所有外观参数都是即时生效 + 防抖落盘。
// ---------------------------------------------------------------------------

const STORAGE_KEY = 'wallpaperCache'

const FIT_VALUES = ['cover', 'contain', 'auto']
const MIN_SURFACE_ALPHA = 0.3
const MAX_SURFACE_ALPHA = 1
// 页面底色比面板更薄，否则壁纸会被整页的底色压掉
const PAGE_ALPHA_RATIO = 0.45
const MIN_PAGE_ALPHA = 0.08
const MAX_PAGE_ALPHA = 0.45
const SAVE_DEBOUNCE_MS = 400

const clamp = (value, min, max) => Math.min(max, Math.max(min, value))

const defaultWallpaper = () => ({
  enabled: false,
  imagePath: '',
  fit: 'cover',
  dim: 0.35,
  surfaceAlpha: 0.72,
  blur: false,
  version: 0,
})

const wallpaper = ref(defaultWallpaper())
const wallpaperLoaded = ref(false)

const normalizeWallpaper = (value) => {
  const fallback = defaultWallpaper()
  const source = value && typeof value === 'object' ? value : {}
  const dim = Number(source.dim)
  const surfaceAlpha = Number(source.surfaceAlpha)
  const imagePath = String(source.imagePath || '').replace(/^\/+/, '')

  return {
    enabled: source.enabled === true && imagePath !== '',
    imagePath,
    fit: FIT_VALUES.includes(source.fit) ? source.fit : fallback.fit,
    dim: Number.isFinite(dim) ? clamp(dim, 0, 1) : fallback.dim,
    surfaceAlpha: Number.isFinite(surfaceAlpha)
      ? clamp(surfaceAlpha, MIN_SURFACE_ALPHA, MAX_SURFACE_ALPHA)
      : fallback.surfaceAlpha,
    blur: source.blur === true,
    version: Number(source.version) || 0,
  }
}

const isWallpaperActive = computed(() => wallpaper.value.enabled && wallpaper.value.imagePath !== '')

// version 是后端按文件修改时间给的，换图后 URL 跟着变，避免浏览器拿旧缓存
const wallpaperUrl = computed(() => {
  const { imagePath, version } = wallpaper.value
  if (!imagePath) return ''
  return version ? `/${imagePath}?v=${version}` : `/${imagePath}`
})

const buildBackgroundLayers = (config) => {
  const layers = []
  if (config.dim > 0) {
    const tint = `rgba(0, 0, 0, ${config.dim.toFixed(3)})`
    layers.push(`linear-gradient(${tint}, ${tint})`)
  }
  const suffix = config.version ? `?v=${config.version}` : ''
  layers.push(`url("/${config.imagePath}${suffix}")`)
  return layers.join(', ')
}

const WALLPAPER_VARS = ['--wp-layers', '--wp-fit', '--wp-surface-pct', '--wp-page-pct']

const applyWallpaperToDom = (value) => {
  const config = normalizeWallpaper(value)
  const root = document.documentElement
  const active = config.enabled && config.imagePath !== ''

  root.classList.toggle('has-wallpaper', active)
  root.classList.toggle('wp-blur', active && config.blur)

  if (!active) {
    WALLPAPER_VARS.forEach((name) => root.style.removeProperty(name))
    return
  }

  const pageAlpha = clamp(config.surfaceAlpha * PAGE_ALPHA_RATIO, MIN_PAGE_ALPHA, MAX_PAGE_ALPHA)
  root.style.setProperty('--wp-layers', buildBackgroundLayers(config))
  root.style.setProperty('--wp-fit', config.fit)
  root.style.setProperty('--wp-surface-pct', `${Math.round(config.surfaceAlpha * 100)}%`)
  root.style.setProperty('--wp-page-pct', `${Math.round(pageAlpha * 100)}%`)
}

const readCache = () => {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}

const writeCache = (config) => {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(config))
  } catch {
    // 存储不可用时忽略：只是少了首帧缓存
  }
}

const setWallpaper = (value, options = {}) => {
  const config = normalizeWallpaper(value)
  wallpaper.value = config
  applyWallpaperToDom(config)
  if (options.cache !== false) writeCache(config)
  return config
}

// 首帧：先用缓存把壁纸铺上，随后 loadWallpaper 再用后端真实配置覆盖
const cachedWallpaper = readCache()
if (cachedWallpaper) {
  setWallpaper(cachedWallpaper, { cache: false })
}

let saveTimer = 0
let pendingSave = false

const flushSave = async () => {
  try {
    const saved = await App.SaveWallpaper(wallpaper.value)
    // 只回填派生字段，避免覆盖用户正在拖动的值
    wallpaper.value = {
      ...wallpaper.value,
      version: Number(saved?.version) || 0,
      imagePath: saved?.imagePath ?? wallpaper.value.imagePath,
      enabled: saved?.enabled ?? wallpaper.value.enabled,
    }
    writeCache(wallpaper.value)
  } catch (error) {
    console.warn('保存壁纸设置失败:', error)
  }
}

const scheduleSave = () => {
  pendingSave = true
  if (saveTimer) clearTimeout(saveTimer)
  saveTimer = window.setTimeout(() => {
    saveTimer = 0
    pendingSave = false
    flushSave()
  }, SAVE_DEBOUNCE_MS)
}

const loadWallpaper = async () => {
  try {
    const config = await App.GetWallpaper()
    setWallpaper(config)
  } catch (error) {
    console.warn('加载壁纸配置失败:', error)
  } finally {
    wallpaperLoaded.value = true
  }
}

const updateWallpaper = (patch) => {
  setWallpaper({ ...wallpaper.value, ...patch })
  scheduleSave()
}

const selectWallpaperImage = async () => {
  const saved = await App.SelectWallpaperImage()
  return setWallpaper(saved)
}

const clearWallpaperImage = async () => {
  const saved = await App.ClearWallpaperImage()
  return setWallpaper(saved)
}

const resetWallpaperAppearance = () => {
  const defaults = defaultWallpaper()
  updateWallpaper({
    fit: defaults.fit,
    dim: defaults.dim,
    surfaceAlpha: defaults.surfaceAlpha,
    blur: defaults.blur,
  })
}

export {
  wallpaper,
  wallpaperLoaded,
  wallpaperUrl,
  isWallpaperActive,
  defaultWallpaper,
  loadWallpaper,
  updateWallpaper,
  selectWallpaperImage,
  clearWallpaperImage,
  resetWallpaperAppearance,
}
