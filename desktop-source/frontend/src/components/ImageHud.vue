<script setup>
import { computed } from 'vue'
import { Sparkles } from 'lucide-vue-next'

const props = defineProps({
  metadata: { type: Object, default: null },
  imageSrc: { type: String, default: '' },
  scale: { type: Number, default: 1 },
  offset: { type: Object, default: () => ({ x: 0, y: 0 }) },
  visible: { type: Boolean, default: true },
  stage: { type: Object, default: () => ({ width: 0, height: 0 }) },
  imageSize: { type: Object, default: () => ({ width: 0, height: 0 }) },
})

defineEmits(['toggle'])

const chips = computed(() => {
  const m = props.metadata
  if (!m) return []

  const list = [
    { label: '模型', value: m.model },
    { label: 'LoRA', value: Array.isArray(m.loras) && m.loras.length ? m.loras.join(', ') : '' },
    { label: '采样器', value: m.sampler },
    { label: '调度器', value: m.scheduler },
    { label: 'Seed', value: m.seed },
    { label: 'Steps', value: m.steps },
    { label: 'CFG', value: m.cfg },
  ]

  return list.filter(item => item.value !== undefined && item.value !== null && item.value !== '')
})

const clamp01 = value => Math.min(1, Math.max(0, value))

const thumbStyle = computed(() => {
  const { width: iw, height: ih } = props.imageSize
  if (!iw || !ih) return { width: '112px', height: '84px' }

  const maxWidth = 112
  const maxHeight = 84
  const ratio = iw / ih
  let width = maxWidth
  let height = width / ratio
  if (height > maxHeight) {
    height = maxHeight
    width = height * ratio
  }
  return { width: `${width.toFixed(1)}px`, height: `${height.toFixed(1)}px` }
})

const viewportBox = computed(() => {
  if (!props.visible || props.scale <= 1) return null

  const { width: stageWidth, height: stageHeight } = props.stage
  const { width: imageWidth, height: imageHeight } = props.imageSize
  if (!stageWidth || !stageHeight || !imageWidth || !imageHeight) return null

  const scaledWidth = imageWidth * props.scale
  const scaledHeight = imageHeight * props.scale
  const left = stageWidth / 2 + props.offset.x - scaledWidth / 2
  const top = stageHeight / 2 + props.offset.y - scaledHeight / 2

  const x = clamp01(-left / scaledWidth)
  const y = clamp01(-top / scaledHeight)
  const width = Math.min(clamp01(stageWidth / scaledWidth), 1 - x)
  const height = Math.min(clamp01(stageHeight / scaledHeight), 1 - y)

  return {
    left: `${(x * 100).toFixed(2)}%`,
    top: `${(y * 100).toFixed(2)}%`,
    width: `${(width * 100).toFixed(2)}%`,
    height: `${(height * 100).toFixed(2)}%`,
  }
})
</script>

<template>
  <div class="pointer-events-none absolute inset-0 z-[60]">
    <button
      type="button"
      class="pointer-events-auto absolute left-4 top-4 flex h-8 w-8 items-center justify-center rounded-full border border-white/15 bg-black/50 backdrop-blur-xl transition-colors hover:bg-white/15 hover:text-white"
      :class="visible ? 'text-blue-300' : 'text-white/60'"
      :title="visible ? '隐藏生成参数' : '显示生成参数'"
      @click="$emit('toggle')"
    >
      <Sparkles class="h-4 w-4" />
    </button>

    <div
      v-if="visible && chips.length"
      class="absolute left-4 top-14 flex max-w-[220px] flex-col items-start gap-1.5"
    >
      <div
        v-for="chip in chips"
        :key="chip.label"
        class="flex max-w-full items-center gap-1.5 rounded-md border border-white/10 bg-black/55 px-2 py-1 text-[11px] shadow-lg backdrop-blur-xl"
      >
        <span class="shrink-0 text-white/45">{{ chip.label }}</span>
        <span class="truncate font-mono text-white/85">{{ chip.value }}</span>
      </div>
    </div>

    <div
      v-if="viewportBox && imageSrc"
      class="absolute bottom-12 right-4 overflow-hidden rounded-md border border-white/15 bg-black/60 shadow-2xl backdrop-blur-xl"
      :style="thumbStyle"
    >
      <img :src="imageSrc" alt="" aria-hidden="true" class="h-full w-full opacity-70" />
      <div class="absolute border border-blue-400 bg-blue-400/20" :style="viewportBox" />
    </div>
  </div>
</template>
