import { onMounted, onUnmounted, ref } from 'vue'

/**
 * 鼠标视差浮动。
 *
 * - 输出 ±strength 像素的位移，用 requestAnimationFrame 节流，不会每个 mousemove 都触发渲染
 * - 传入 disabled 返回 true 时（例如图片已放大、正在拖拽平移）位移恒为 0，避免和拖拽打架
 * - 鼠标移出窗口自动归零
 */
export function useParallaxFloat({ strength = 6, disabled = () => false } = {}) {
  const x = ref(0)
  const y = ref(0)
  let frame = 0
  let targetX = 0
  let targetY = 0

  const commit = () => {
    frame = 0
    x.value = targetX
    y.value = targetY
  }

  const schedule = () => {
    if (!frame) frame = requestAnimationFrame(commit)
  }

  const handleMouseMove = event => {
    if (disabled()) {
      if (targetX !== 0 || targetY !== 0) {
        targetX = 0
        targetY = 0
        schedule()
      }
      return
    }
    targetX = (event.clientX / window.innerWidth - 0.5) * 2 * strength
    targetY = (event.clientY / window.innerHeight - 0.5) * 2 * strength
    schedule()
  }

  const handleMouseLeave = () => {
    targetX = 0
    targetY = 0
    schedule()
  }

  onMounted(() => {
    window.addEventListener('mousemove', handleMouseMove, { passive: true })
    window.addEventListener('mouseleave', handleMouseLeave)
  })

  onUnmounted(() => {
    window.removeEventListener('mousemove', handleMouseMove)
    window.removeEventListener('mouseleave', handleMouseLeave)
    if (frame) cancelAnimationFrame(frame)
  })

  return { x, y }
}
