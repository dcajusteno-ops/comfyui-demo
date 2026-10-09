import { ref, watch } from 'vue'

const SAMPLE_SIZE = 32
const CACHE_LIMIT = 60
const colorCache = new Map()

function rgbToHue(r, g, b) {
  const max = Math.max(r, g, b)
  const min = Math.min(r, g, b)
  if (max === min) return 0
  const delta = max - min
  let hue
  if (max === r) hue = ((g - b) / delta) % 6
  else if (max === g) hue = (b - r) / delta + 2
  else hue = (r - g) / delta + 4
  hue *= 60
  if (hue < 0) hue += 360
  return hue
}

/**
 * 从 32x32 采样数据里挑主色。
 * 做法：跳过过暗 / 过亮 / 低饱和的像素，按色相每 15 度分桶，取像素最多的桶求平均。
 * 没有有效像素时回退到全图平均色（避免出现纯黑 / 纯灰）。
 */
function pickDominantColor(data) {
  const buckets = new Map()
  let rSum = 0
  let gSum = 0
  let bSum = 0
  let total = 0

  for (let i = 0; i < data.length; i += 4) {
    if (data[i + 3] < 128) continue
    const r = data[i]
    const g = data[i + 1]
    const b = data[i + 2]
    rSum += r
    gSum += g
    bSum += b
    total += 1

    const max = Math.max(r, g, b)
    const min = Math.min(r, g, b)
    const lightness = (max + min) / 2
    const saturation = max === min ? 0 : (max - min) / (255 - Math.abs(max + min - 255))
    if (lightness < 24 || lightness > 240 || saturation < 0.12) continue

    const key = Math.floor(rgbToHue(r, g, b) / 15) % 24
    const bucket = buckets.get(key) || { r: 0, g: 0, b: 0, n: 0 }
    bucket.r += r
    bucket.g += g
    bucket.b += b
    bucket.n += 1
    buckets.set(key, bucket)
  }

  let best = null
  for (const bucket of buckets.values()) {
    if (!best || bucket.n > best.n) best = bucket
  }

  if (best && best.n > 0) {
    return {
      r: Math.round(best.r / best.n),
      g: Math.round(best.g / best.n),
      b: Math.round(best.b / best.n),
    }
  }

  if (total > 0) {
    return {
      r: Math.round(rSum / total),
      g: Math.round(gSum / total),
      b: Math.round(bSum / total),
    }
  }

  return { r: 96, g: 104, b: 128 }
}

function sampleColorFromImage(src) {
  return new Promise((resolve, reject) => {
    const image = new Image()
    image.decoding = 'async'
    image.onload = () => {
      try {
        const canvas = document.createElement('canvas')
        canvas.width = SAMPLE_SIZE
        canvas.height = SAMPLE_SIZE
        const ctx = canvas.getContext('2d')
        if (!ctx) {
          reject(new Error('canvas 2d context unavailable'))
          return
        }
        ctx.drawImage(image, 0, 0, SAMPLE_SIZE, SAMPLE_SIZE)
        const { data } = ctx.getImageData(0, 0, SAMPLE_SIZE, SAMPLE_SIZE)
        resolve(pickDominantColor(data))
      } catch (error) {
        reject(error)
      }
    }
    image.onerror = () => reject(new Error('image load failed'))
    image.src = src
  })
}

/**
 * 采样图片主色（供氛围光晕使用）。
 *
 * - 只应传 thumbPath 这类小图，别传原图
 * - 内部带 requestId 竞态守卫：快速切图时旧结果会被丢弃，不会覆盖新图颜色
 * - 结果按完整 URL（含 ?v= 版本号）缓存，回看旧图零成本
 * - 采样失败时保留上一次颜色，避免闪黑
 */
export function useDominantColor(getSrc) {
  const color = ref(null)
  let requestId = 0

  const resolveSrc = () => (typeof getSrc === 'function' ? getSrc() : getSrc?.value)

  const load = () => {
    const src = resolveSrc()
    if (!src) {
      color.value = null
      return
    }

    const cached = colorCache.get(src)
    if (cached) {
      color.value = cached
      return
    }

    const rid = ++requestId
    sampleColorFromImage(src)
      .then(rgb => {
        if (rid !== requestId) return
        if (colorCache.size >= CACHE_LIMIT) {
          const oldest = colorCache.keys().next().value
          if (oldest !== undefined) colorCache.delete(oldest)
        }
        colorCache.set(src, rgb)
        color.value = rgb
      })
      .catch(() => {
        // 静默降级：保留当前颜色
      })
  }

  watch(
    () => resolveSrc(),
    () => {
      load()
    },
    { immediate: true },
  )

  return { color }
}
