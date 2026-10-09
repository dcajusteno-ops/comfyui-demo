<script setup>
import { computed } from 'vue'
import { useDominantColor } from '@/composables/useDominantColor'

const props = defineProps({
  imageSrc: { type: String, default: '' },
})

const { color } = useDominantColor(() => props.imageSrc)

const toGlow = alpha => {
  const c = color.value
  if (!c) return 'transparent'
  // 暗色 / 低饱和的图也要有可感知的氛围光：把最亮通道拉到 255（最多放大 2.2 倍）
  const peak = Math.max(c.r, c.g, c.b) || 1
  const gain = Math.min(255 / peak, 2.2)
  const r = Math.min(255, Math.round(c.r * gain))
  const g = Math.min(255, Math.round(c.g * gain))
  const b = Math.min(255, Math.round(c.b * gain))
  return `rgba(${r}, ${g}, ${b}, ${alpha})`
}

const leftGlow = computed(() => toGlow(0.66))
const rightGlow = computed(() => toGlow(0.5))
</script>

<template>
  <div class="pointer-events-none absolute inset-0 z-0 overflow-hidden">
    <div
      class="absolute -left-[12%] top-1/2 h-[78%] w-[46%] -translate-y-1/2 rounded-full blur-[100px] transition-colors duration-700 ease-out"
      :style="{ backgroundColor: leftGlow }"
    />
    <div
      class="absolute -right-[12%] top-1/2 h-[66%] w-[42%] -translate-y-1/2 rounded-full blur-[100px] transition-colors duration-700 ease-out"
      :style="{ backgroundColor: rightGlow }"
    />
  </div>
</template>
