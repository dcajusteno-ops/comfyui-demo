<script setup>
import { computed, ref } from 'vue'

const CARD_WIDTH = 38
const CARD_HEIGHT = 60
const MAX_VISIBLE = 9

const props = defineProps({
  image: { type: Object, default: null },
  currentIndex: { type: Number, default: 1 },
})

const emit = defineEmits(['select'])

const expanded = ref(false)

const items = computed(() => {
  if (!props.image?.isStackPrimary) return []
  return [props.image, ...(props.image.stackChildren || [])]
})

const activeIndex = computed(() =>
  Math.max(0, Math.min(items.value.length - 1, props.currentIndex - 1)),
)

// 以当前项为中心开窗，避免一次渲染上百张卡牌
const visibleCards = computed(() => {
  const list = items.value
  const total = list.length
  if (total <= MAX_VISIBLE) {
    return list.map((img, index) => ({ img, index }))
  }

  let start = Math.max(0, activeIndex.value - Math.floor(MAX_VISIBLE / 2))
  if (start + MAX_VISIBLE > total) start = total - MAX_VISIBLE

  return Array.from({ length: MAX_VISIBLE }, (_, offset) => {
    const index = start + offset
    return { img: list[index], index }
  })
})

const stageWidth = computed(() => {
  const spread = expanded.value ? 38 : 14
  const count = visibleCards.value.length
  return `${(count - 1) * spread + CARD_WIDTH + 28}px`
})

const cardStyle = (entry, position) => {
  const center = (visibleCards.value.length - 1) / 2
  const distance = position - center
  const isActive = entry.index === activeIndex.value

  const spread = expanded.value ? 38 : 14
  const lift = Math.abs(distance) * (expanded.value ? 3 : 5)
  const rotate = distance * (expanded.value ? 1.8 : 1.2)

  return {
    left: '50%',
    marginLeft: `${-CARD_WIDTH / 2}px`,
    width: `${CARD_WIDTH}px`,
    height: `${CARD_HEIGHT}px`,
    zIndex: String(100 - Math.round(Math.abs(distance) * 2)),
    transform: `translate(${(distance * spread).toFixed(1)}px, ${lift.toFixed(1)}px) rotate(${rotate.toFixed(2)}deg) scale(${isActive ? 1.18 : 1})`,
  }
}
</script>

<template>
  <div
    v-if="items.length > 1"
    class="absolute bottom-[4.5rem] left-[50%] z-[60] -translate-x-1/2"
    @click.stop
    @wheel.stop
  >
    <div class="mb-1.5 flex items-center justify-center gap-2">
      <span class="text-[10px] font-semibold uppercase tracking-widest text-white/45">连拍组</span>
      <span class="font-mono text-xs text-white/80">{{ currentIndex }} / {{ items.length }}</span>
    </div>

    <div
      class="relative flex h-[64px] items-end justify-center transition-[width] duration-300 ease-out"
      :style="{ width: stageWidth }"
      @mouseenter="expanded = true"
      @mouseleave="expanded = false"
    >
      <button
        v-for="(entry, position) in visibleCards"
        :key="entry.img.relPath"
        type="button"
        class="absolute bottom-0 cursor-pointer overflow-hidden rounded-md border bg-black/60 transition-all duration-300 ease-out"
        :class="
          entry.index === activeIndex
            ? 'border-blue-400 shadow-[0_0_18px_rgba(96,165,250,0.5)]'
            : 'border-white/25 hover:border-white/60'
        "
        :style="cardStyle(entry, position)"
        :title="entry.img.name"
        @click="emit('select', entry.index)"
      >
        <img
          :src="entry.img.thumbPath || entry.img.previewPath || entry.img.path"
          alt=""
          aria-hidden="true"
          draggable="false"
          loading="lazy"
          decoding="async"
          class="h-full w-full object-cover"
        />
      </button>
    </div>
  </div>
</template>
