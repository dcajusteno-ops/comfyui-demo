<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'

const props = defineProps({
  images: { type: Array, default: () => [] },
})

const emit = defineEmits(['view'])

const CARD_WIDTH = 132
const CARD_HEIGHT = 176
const RING_GAP = 48
const MIN_RADIUS = 160
const AUTO_SPIN_SPEED = 0.045
// 拖动超过这个距离就当成旋转/滚动，不再触发点击
const CLICK_THRESHOLD = 6
const PER_RING_MIN = 4
const PER_RING_MAX = 40
const PER_RING_STEP = 2
// 滚到顶/底时给塔留的空白（屏幕像素）
const EDGE_PADDING = 100
// 透视距离 = 半径 × 该系数：近端放大率恒定为 ≈1.31×，不随每圈张数暴涨
const PERSPECTIVE_FACTOR = 4.2
const PERSPECTIVE_MIN = 1400
const PERSPECTIVE_MAX = 6400
// 环绕整体最多占视口宽度的比例，超出则整体等比缩小（自适应）
const FIT_WIDTH_RATIO = 0.94
const FIT_SCALE_MIN = 0.22
const PILLAR_Z_GAP = 60
// 每秒滚轮/拖动位移换算：屏幕像素 → 场景单位
const SCROLL_SPEED = 0.8

const storedPerRing = Number(localStorage.getItem('orbitPerRing'))
const perRing = ref(
  Number.isFinite(storedPerRing) && storedPerRing >= PER_RING_MIN && storedPerRing <= PER_RING_MAX
    ? storedPerRing
    : 10,
)

const rotation = ref(0)
const shiftY = ref(0)
const viewportHeight = ref(600)
const viewportWidth = ref(1000)
const isDragging = ref(false)
const viewportEl = ref(null)

let frame = 0
let resizeObserver = null
let dragging = false
let lastX = 0
let lastY = 0
let velocity = 0
let movedDistance = 0

const total = computed(() => props.images.length)

// 圈数完全由「每圈张数」决定，不再封顶
const ringCount = computed(() => {
  if (!total.value) return 0
  return Math.max(1, Math.ceil(total.value / perRing.value))
})

// 每圈实际张数（最后一圈可能不满）
const ringSizes = computed(() => {
  const sizes = []
  for (let i = 0; i < ringCount.value; i += 1) {
    sizes.push(Math.max(0, Math.min(perRing.value, total.value - i * perRing.value)))
  }
  return sizes
})

const radius = computed(() => {
  const per = Math.max(3, perRing.value)
  const ideal = CARD_WIDTH / 2 / Math.tan(Math.PI / per)
  return Math.max(MIN_RADIUS, ideal * 1.04)
})

// 透视距离跟着半径走，而不是写死 —— 写死会把近处卡片放得过大、撑破视口
const perspective = computed(() =>
  Math.min(PERSPECTIVE_MAX, Math.max(PERSPECTIVE_MIN, radius.value * PERSPECTIVE_FACTOR)),
)

// 近端卡片放大率，只由「透视 / 半径」之比决定，与窗口大小无关
const nearMagnify = computed(
  () => perspective.value / Math.max(1, perspective.value - radius.value),
)

// 整圈投影后的半宽（采样计算，含卡片自身宽度在透视下的贡献）
const projectedHalfWidth = computed(() => {
  const r = radius.value
  const p = perspective.value
  const half = CARD_WIDTH / 2
  let max = 0
  for (let i = 0; i <= 72; i += 1) {
    const theta = (Math.PI / 2) * (i / 72)
    const scale = p / Math.max(1, p - r * Math.cos(theta))
    const extent = (r * Math.sin(theta) + half * Math.abs(Math.cos(theta))) * scale
    if (extent > max) max = extent
  }
  return max
})

// 自适应：环绕整体必须塞进视口宽度（只缩不放）
const fitScale = computed(() => {
  const width = projectedHalfWidth.value * 2
  if (!width) return 1
  const target = viewportWidth.value * FIT_WIDTH_RATIO
  return Math.min(1, Math.max(FIT_SCALE_MIN, target / width))
})

// 塔在屏幕上的真实高度（把近端放大也算进去），用于计算可滚动范围
const contentScreenHeight = computed(() => {
  if (!ringCount.value) return CARD_HEIGHT
  const halfSpan = ((ringCount.value - 1) / 2) * (CARD_HEIGHT + RING_GAP)
  const halfHeight = halfSpan + (CARD_HEIGHT / 2) * nearMagnify.value
  return halfHeight * 2 * fitScale.value
})

// 塔超出视口才允许上下滚；结果换算回场景单位
const maxShift = computed(() => {
  const overflow = contentScreenHeight.value + EDGE_PADDING * 2 - viewportHeight.value
  if (overflow <= 0) return 0
  return overflow / 2 / fitScale.value
})

const ringY = ringIndex => {
  const middle = (ringCount.value - 1) / 2
  return (ringIndex - middle) * (CARD_HEIGHT + RING_GAP)
}

const cardStyle = index => {
  const ringIndex = Math.floor(index / perRing.value)
  const indexInRing = index % perRing.value
  const size = ringSizes.value[ringIndex] || 1
  const angle = (indexInRing / size) * 360
  return {
    width: `${CARD_WIDTH}px`,
    height: `${CARD_HEIGHT}px`,
    marginLeft: `${-CARD_WIDTH / 2}px`,
    marginTop: `${-CARD_HEIGHT / 2}px`,
    transform: `translateY(${ringY(ringIndex).toFixed(1)}px) rotateY(${angle.toFixed(2)}deg) translateZ(${radius.value.toFixed(1)}px)`,
  }
}

// 外层 2D 缩放：把整个 3D 渲染结果等比缩到视口内（挂在透视容器之外，避免影响 3D 上下文）
const stageStyle = computed(() => ({
  transform: `scale(${fitScale.value.toFixed(4)})`,
}))

const perspectiveStyle = computed(() => ({ perspective: `${perspective.value.toFixed(0)}px` }))

// 柱子始终贯穿视口（含可滚动的那段位移），不会滚到一半就断头
const pillarStyle = computed(() => {
  const shrink = perspective.value / (perspective.value + radius.value + PILLAR_Z_GAP)
  const screenNeed = viewportHeight.value + 2 * maxShift.value * fitScale.value
  const sceneNeed = screenNeed / Math.max(0.12, fitScale.value * shrink)
  const towerScene = contentScreenHeight.value / Math.max(0.12, fitScale.value)
  const height = Math.max(towerScene, sceneNeed) * 1.04
  return {
    height: `${height.toFixed(1)}px`,
    transform: `translate(-50%, -50%) translateZ(${-(radius.value + PILLAR_Z_GAP)}px)`,
    background:
      'linear-gradient(to right, rgba(255,255,255,0.03), rgba(255,255,255,0.22), rgba(255,255,255,0.03))',
    boxShadow: '0 0 40px rgba(255,255,255,0.08)',
  }
})

const sceneStyle = computed(() => ({
  transform: `translateY(${shiftY.value.toFixed(1)}px)`,
}))

const clampShift = value => Math.min(maxShift.value, Math.max(-maxShift.value, value))

const measure = () => {
  if (!viewportEl.value) return
  viewportWidth.value = viewportEl.value.clientWidth || 1000
  viewportHeight.value = viewportEl.value.clientHeight || 600
  shiftY.value = clampShift(shiftY.value)
}

const setPerRing = delta => {
  const next = Math.min(PER_RING_MAX, Math.max(PER_RING_MIN, perRing.value + delta))
  if (next === perRing.value) return
  perRing.value = next
  localStorage.setItem('orbitPerRing', String(next))
  shiftY.value = clampShift(shiftY.value)
}

const tick = () => {
  if (!dragging) {
    if (Math.abs(velocity) > 0.02) {
      rotation.value += velocity
      velocity *= 0.93
    } else {
      velocity = 0
      rotation.value += AUTO_SPIN_SPEED
    }
  }
  frame = requestAnimationFrame(tick)
}

// 用 window 监听而不是 setPointerCapture：
// pointer capture 会把 click 事件也劫持到捕获元素上，导致卡片点不开大图
const onWindowPointerMove = event => {
  if (!dragging) return
  const dx = event.clientX - lastX
  const dy = event.clientY - lastY
  lastX = event.clientX
  lastY = event.clientY
  movedDistance += Math.abs(dx) + Math.abs(dy)
  // 横向分量旋转，纵向分量上下滚动（按屏幕距离 1:1 换算回场景单位）
  velocity = dx * 0.32
  rotation.value += dx * 0.32
  if (maxShift.value > 0) {
    shiftY.value = clampShift(shiftY.value + dy / Math.max(0.2, fitScale.value))
  }
}

const stopDragging = () => {
  dragging = false
  isDragging.value = false
  window.removeEventListener('pointermove', onWindowPointerMove)
  window.removeEventListener('pointerup', stopDragging)
  window.removeEventListener('pointercancel', stopDragging)
}

const onPointerDown = event => {
  dragging = true
  isDragging.value = true
  lastX = event.clientX
  lastY = event.clientY
  velocity = 0
  movedDistance = 0
  window.addEventListener('pointermove', onWindowPointerMove)
  window.addEventListener('pointerup', stopDragging)
  window.addEventListener('pointercancel', stopDragging)
}

const handleCardClick = img => {
  // 刚刚在拖拽，不算点击
  if (movedDistance > CLICK_THRESHOLD) return
  emit('view', img)
}

const onWheel = event => {
  // 按住 Shift：改为旋转
  if (event.shiftKey) {
    rotation.value += event.deltaY * 0.1
    velocity = 0
    return
  }
  // 塔没超出视口时，滚轮也用来旋转
  if (maxShift.value <= 0) {
    rotation.value += event.deltaY * 0.1
    velocity = 0
    return
  }
  // 滚轮向下（deltaY > 0）＝ 想看塔的下方，内容应当上移（和正常滚动一致）
  const step = (event.deltaY * SCROLL_SPEED) / Math.max(0.2, fitScale.value)
  shiftY.value = clampShift(shiftY.value - step)
}

onMounted(() => {
  frame = requestAnimationFrame(tick)
  measure()
  if (typeof ResizeObserver !== 'undefined' && viewportEl.value) {
    resizeObserver = new ResizeObserver(measure)
    resizeObserver.observe(viewportEl.value)
  }
})

onUnmounted(() => {
  if (frame) cancelAnimationFrame(frame)
  if (resizeObserver) {
    resizeObserver.disconnect()
    resizeObserver = null
  }
  stopDragging()
})
</script>

<template>
  <div
    ref="viewportEl"
    class="wp-scrim relative h-full w-full select-none overflow-hidden bg-neutral-950"
    :class="isDragging ? 'cursor-grabbing' : 'cursor-grab'"
    @pointerdown="onPointerDown"
    @wheel.prevent="onWheel"
  >
    <div
      class="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_at_center,rgba(255,255,255,0.07),transparent_65%)]"
    />

    <!-- 自适应缩放层：把整幅 3D 画面等比缩到视口内 -->
    <div class="absolute inset-0 transition-transform duration-200 ease-out" :style="stageStyle">
      <div class="absolute inset-0" :style="perspectiveStyle">
        <div
          class="absolute left-1/2 top-1/2"
          style="transform-style: preserve-3d"
          :style="sceneStyle"
        >
          <div class="absolute left-0 top-0 w-4 rounded-full" :style="pillarStyle" />

          <div
            class="absolute left-0 top-0"
            style="transform-style: preserve-3d"
            :style="{ transform: `rotateY(${rotation}deg)` }"
          >
            <button
              v-for="(img, index) in images"
              :key="img.relPath"
              type="button"
              class="absolute left-0 top-0 overflow-hidden rounded-xl border border-white/15 bg-neutral-900 shadow-[0_20px_45px_rgba(0,0,0,0.6)] transition-[box-shadow,filter,border-color] duration-200 hover:border-white/50 hover:brightness-110"
              :style="cardStyle(index)"
              :title="img.name"
              @click="handleCardClick(img)"
            >
              <img
                :src="img.thumbPath || img.previewPath || img.path"
                alt=""
                aria-hidden="true"
                draggable="false"
                loading="lazy"
                decoding="async"
                class="pointer-events-none h-full w-full object-cover"
              />
            </button>
          </div>
        </div>
      </div>
    </div>

    <div
      class="absolute left-5 top-5 z-20 flex items-center gap-1.5 rounded-xl border border-white/10 bg-black/65 px-3 py-2 text-xs text-white/70 backdrop-blur-xl"
      @pointerdown.stop
      @wheel.stop
    >
      <span>每圈</span>
      <button
        type="button"
        class="flex h-6 w-6 items-center justify-center rounded-md text-white/70 transition-colors hover:bg-white/15 hover:text-white disabled:pointer-events-none disabled:opacity-30"
        :disabled="perRing <= PER_RING_MIN"
        title="减少每圈张数"
        @click="setPerRing(-PER_RING_STEP)"
      >
        −
      </button>
      <span class="w-7 text-center font-mono text-sm text-white">{{ perRing }}</span>
      <button
        type="button"
        class="flex h-6 w-6 items-center justify-center rounded-md text-white/70 transition-colors hover:bg-white/15 hover:text-white disabled:pointer-events-none disabled:opacity-30"
        :disabled="perRing >= PER_RING_MAX"
        title="增加每圈张数"
        @click="setPerRing(PER_RING_STEP)"
      >
        +
      </button>
      <span>张</span>
      <span class="ml-1 border-l border-white/15 pl-2 text-white/45">共 {{ ringCount }} 圈</span>
    </div>

    <div class="pointer-events-none absolute bottom-5 left-6 text-xs leading-relaxed text-white/40">
      拖动旋转 · {{ maxShift > 0 ? '滚轮上下滚动' : '滚轮旋转' }} · 点击图片查看大图
      <br />
      共 {{ images.length }} 张 · 每圈 {{ perRing }} 张 · {{ ringCount }} 圈
    </div>
  </div>
</template>
