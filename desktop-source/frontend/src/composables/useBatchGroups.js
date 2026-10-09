import { computed, ref } from 'vue'
import { getImageDateKey } from './useGalleryHelpers'

export const BATCH_MODES = [
  { value: 'day', label: '同日' },
  { value: 'model', label: '同模型' },
  { value: 'prompt', label: '同提示词' },
  { value: 'stack', label: '连拍组' },
]

const UNKNOWN_KEY = '__unknown__'

const truncate = (value, max = 36) => {
  const text = String(value || '').trim()
  return text.length > max ? `${text.slice(0, max)}…` : text
}

/**
 * 由当前图与图库列表派生「当前批次」。
 *
 * 四种口径：
 * - day    同一天产出（复用 getImageDateKey，文件夹日期优先，回退修改时间）
 * - model  同一模型
 * - prompt 同一提示词
 * - stack  当前连拍组（用已有的 stackChildren）
 *
 * 每个条目都带上它在 images 数组里的原始下标，供点击跳转使用。
 * 注意 images 是 Lightbox 收到的那个数组（连拍折叠后的），下标必须与之一致。
 */
export function useBatchGroups(getImages, getCurrentImage) {
  const mode = ref('day')

  const group = computed(() => {
    const images = Array.isArray(getImages()) ? getImages() : []
    const current = getCurrentImage()
    if (!current || !images.length) return null

    const currentIndex = images.findIndex(img => img.relPath === current.relPath)

    if (mode.value === 'stack') {
      const stackItems = current.isStackPrimary
        ? [current, ...(current.stackChildren || [])]
        : [current]

      const items = stackItems
        .map(img => ({ img, index: images.findIndex(entry => entry.relPath === img.relPath) }))
        .filter(entry => entry.index >= 0)

      return {
        key: 'stack',
        label: '当前连拍组',
        items: items.length ? items : [{ img: current, index: currentIndex }],
      }
    }

    const keyOf = img => {
      if (mode.value === 'model') return img?.model || UNKNOWN_KEY
      if (mode.value === 'prompt') return img?.prompt || UNKNOWN_KEY
      return getImageDateKey(img) || UNKNOWN_KEY
    }

    const targetKey = keyOf(current)
    const items = images
      .map((img, index) => ({ img, index }))
      .filter(entry => keyOf(entry.img) === targetKey)

    let label = '当前批次'
    if (targetKey === UNKNOWN_KEY) label = '未识别'
    else if (mode.value === 'prompt') label = truncate(targetKey)
    else if (mode.value === 'model') label = truncate(targetKey)
    else label = targetKey

    return { key: targetKey, label, items }
  })

  return { mode, group }
}
