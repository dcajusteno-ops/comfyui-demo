import { ref } from 'vue'

/**
 * 图片缩放 / 平移交互。
 *
 * 行为与 Lightbox.vue 中既有实现保持一致，供全屏灯箱与灵动小窗共用：
 * - 滚轮缩放，区间默认 0.5x – 10x
 * - 仅在 scale > 1 时允许拖拽平移
 */
export function useImageZoom({ min = 0.5, max = 10, factor = 1.1 } = {}) {
  const scale = ref(1)
  const offset = ref({ x: 0, y: 0 })
  const isDragging = ref(false)
  const lastMousePos = ref({ x: 0, y: 0 })

  const resetZoom = () => {
    scale.value = 1
    offset.value = { x: 0, y: 0 }
    isDragging.value = false
  }

  const handleWheel = event => {
    const delta = -event.deltaY
    const nextScale = delta > 0 ? scale.value * factor : scale.value / factor
    scale.value = Math.min(Math.max(nextScale, min), max)
  }

  const handleMouseDown = event => {
    if (scale.value <= 1) return
    isDragging.value = true
    lastMousePos.value = { x: event.clientX, y: event.clientY }
    event.preventDefault()
  }

  const handleMouseMove = event => {
    if (!isDragging.value) return
    const dx = event.clientX - lastMousePos.value.x
    const dy = event.clientY - lastMousePos.value.y
    offset.value = {
      x: offset.value.x + dx,
      y: offset.value.y + dy,
    }
    lastMousePos.value = { x: event.clientX, y: event.clientY }
  }

  const handleMouseUp = () => {
    isDragging.value = false
  }

  return {
    scale,
    offset,
    isDragging,
    handleWheel,
    handleMouseDown,
    handleMouseMove,
    handleMouseUp,
    resetZoom,
  }
}
