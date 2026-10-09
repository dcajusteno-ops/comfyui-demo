<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { ChevronLeft, ChevronRight, Loader2, RotateCcw } from 'lucide-vue-next'
import GroundReflection from './GroundReflection.vue'
import ImageFilmstrip from './ImageFilmstrip.vue'
import ImageHud from './ImageHud.vue'
import StackFan from './StackFan.vue'
import { useParallaxFloat } from '@/composables/useParallaxFloat'

const props = defineProps({
  image: { type: Object, default: null },
  currentDisplayImage: { type: Object, default: null },
  canGoPrev: { type: Boolean, default: false },
  canGoNext: { type: Boolean, default: false },
  displayImageSrc: { type: String, default: '' },
  fullImageLoading: { type: Boolean, default: false },
  scale: { type: Number, default: 1 },
  offset: {
    type: Object,
    default: () => ({ x: 0, y: 0 }),
  },
  isDragging: { type: Boolean, default: false },
  imageCounter: { type: String, default: '1 / 1' },
  stackCurrentIndex: { type: Number, default: 1 },
  metadata: { type: Object, default: null },
  hudVisible: { type: Boolean, default: true },
  images: { type: Array, default: () => [] },
  currentIndex: { type: Number, default: 0 },
})

defineEmits([
  'prev',
  'next',
  'prev-stack',
  'next-stack',
  'reset-zoom',
  'viewer-wheel',
  'viewer-mousedown',
  'viewer-mousemove',
  'viewer-mouseup',
  'toggle-hud',
  'select-stack',
  'toggle-overview',
  'jump-to-index',
])

const { x: parallaxX, y: parallaxY } = useParallaxFloat({
  strength: 6,
  disabled: () => props.scale > 1,
})

const parallaxStyle = computed(() => ({
  transform: `translate3d(${parallaxX.value}px, ${parallaxY.value}px, 0)`,
}))

// 相纸显影：换图时重播一次动画。只切换 class、不重建 img，避免额外闪一下
const developing = ref(false)
let developFrame = 0

watch(
  () => props.currentDisplayImage?.relPath,
  async () => {
    developing.value = false
    if (developFrame) cancelAnimationFrame(developFrame)
    await nextTick()
    developFrame = requestAnimationFrame(() => {
      developFrame = 0
      developing.value = true
    })
  },
  { immediate: true },
)

onUnmounted(() => {
  if (developFrame) cancelAnimationFrame(developFrame)
})

// 供 HUD 计算「当前视野在原图中的位置」使用
const stageEl = ref(null)
const mainImageEl = ref(null)
const stage = ref({ width: 0, height: 0 })
const imageSize = ref({ width: 0, height: 0 })
let resizeObserver = null

const measure = () => {
  if (stageEl.value) {
    stage.value = { width: stageEl.value.clientWidth, height: stageEl.value.clientHeight }
  }
  if (mainImageEl.value) {
    imageSize.value = { width: mainImageEl.value.offsetWidth, height: mainImageEl.value.offsetHeight }
  }
}

onMounted(() => {
  measure()
  if (typeof ResizeObserver === 'undefined') return
  resizeObserver = new ResizeObserver(measure)
  if (stageEl.value) resizeObserver.observe(stageEl.value)
  if (mainImageEl.value) resizeObserver.observe(mainImageEl.value)
})

onUnmounted(() => {
  if (resizeObserver) {
    resizeObserver.disconnect()
    resizeObserver = null
  }
})

watch(
  () => props.displayImageSrc,
  () => {
    nextTick(measure)
  },
)
</script>

<template>
  <div class="absolute inset-y-0 left-0 right-[408px]">
    <button
      v-if="canGoPrev"
      class="absolute left-8 top-1/2 z-[70] flex h-12 w-12 -translate-y-1/2 items-center justify-center rounded-full bg-white/10 text-white transition-all hover:scale-110 hover:bg-white/20"
      @click="$emit('prev')"
    >
      <ChevronLeft class="h-8 w-8" />
    </button>

    <button
      v-if="canGoNext"
      class="absolute right-8 top-1/2 z-[70] flex h-12 w-12 -translate-y-1/2 items-center justify-center rounded-full bg-white/10 text-white transition-all hover:scale-110 hover:bg-white/20"
      @click="$emit('next')"
    >
      <ChevronRight class="h-8 w-8" />
    </button>

      <div
        ref="stageEl"
        class="relative flex h-full w-full flex-col items-center justify-center overflow-hidden px-24"
        :class="{ 'cursor-grab': scale > 1 && !isDragging, 'cursor-grabbing': isDragging }"
      @wheel="$emit('viewer-wheel', $event)"
      @mousedown="$emit('viewer-mousedown', $event)"
      @mousemove="$emit('viewer-mousemove', $event)"
      @mouseup="$emit('viewer-mouseup')"
      @mouseleave="$emit('viewer-mouseup')"
    >
      <div
        v-if="fullImageLoading"
        class="absolute top-8 z-[65] flex items-center gap-2 rounded-full border border-white/10 bg-black/55 px-3 py-1.5 text-sm text-white/80 shadow-lg backdrop-blur-xl"
      >
        <Loader2 class="h-4 w-4 animate-spin" />
        <span>正在载入原图</span>
      </div>
      <div class="relative select-none" :style="parallaxStyle">
        <div
          class="relative select-none transition-transform duration-75 ease-out"
          :style="{ transform: `translate(${offset.x}px, ${offset.y}px) scale(${scale})` }"
        >
          <div class="relative" :class="{ 'paper-develop': developing }">
            <img
              ref="mainImageEl"
              :src="displayImageSrc || currentDisplayImage?.path"
              :alt="currentDisplayImage?.name || ''"
              loading="eager"
              decoding="async"
              class="pointer-events-none max-h-[calc(100vh-120px)] max-w-full rounded object-contain shadow-2xl"
              @load="measure"
            />
            <GroundReflection :image-src="displayImageSrc || currentDisplayImage?.path" />
          </div>
        </div>
      </div>
    </div>

    <ImageFilmstrip
      :images="images"
      :current-index="currentIndex"
      @jump="$emit('jump-to-index', $event)"
      @toggle-overview="$emit('toggle-overview')"
    />

    <ImageHud
      :metadata="metadata"
      :image-src="displayImageSrc || currentDisplayImage?.path || ''"
      :scale="scale"
      :offset="offset"
      :visible="hudVisible"
      :stage="stage"
      :image-size="imageSize"
      @toggle="$emit('toggle-hud')"
    />
  </div>

  <div
    class="group absolute bottom-[4.5rem] left-8 z-[60] flex min-w-[220px] max-w-md flex-col gap-3 rounded-xl border border-white/10 bg-black/70 p-4 text-white shadow-2xl backdrop-blur-xl transition-all hover:bg-black/80"
    @click.stop
  >
    <div class="flex items-center">
      <span class="truncate text-base font-semibold tracking-wide text-white/90">{{ currentDisplayImage?.name }}</span>
    </div>

    <div class="h-px w-full bg-white/10" />

    <div class="flex items-center justify-between text-xs">
      <div class="flex items-center gap-3 font-mono text-white/60">
        <span>{{ imageCounter }}</span>
        <span class="h-3 w-px bg-white/20" />
        <span :class="{ 'font-bold text-blue-400': scale !== 1 }">{{ Math.round(scale * 100) }}%</span>
      </div>

      <button
        v-if="scale !== 1 || offset.x !== 0 || offset.y !== 0"
        class="-mr-2 flex items-center gap-1.5 rounded px-2 py-1 text-xs font-medium text-blue-300 transition-all hover:bg-white/20 hover:text-blue-200"
        title="重置缩放"
        @click="$emit('reset-zoom')"
      >
        <RotateCcw class="h-3 w-3" />
        重置视图
      </button>
    </div>
  </div>

  <StackFan
    v-if="image?.isStackPrimary && image.stackCount > 1"
    :image="image"
    :current-index="stackCurrentIndex"
    @select="$emit('select-stack', $event)"
  />
</template>

<style scoped>
@keyframes paper-develop {
  from {
    opacity: 0;
    filter: blur(22px) brightness(1.35);
    transform: scale(1.04);
  }
  to {
    opacity: 1;
    filter: blur(0) brightness(1);
    transform: scale(1);
  }
}

.paper-develop {
  animation: paper-develop 0.55s cubic-bezier(0.22, 0.61, 0.36, 1) both;
}
</style>
