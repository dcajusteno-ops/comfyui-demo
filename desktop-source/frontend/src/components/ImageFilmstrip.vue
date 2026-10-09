<script setup>
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { ChevronLeft, ChevronRight, LayoutGrid } from 'lucide-vue-next'

const props = defineProps({
  images: { type: Array, default: () => [] },
  currentIndex: { type: Number, default: 0 },
})

const emit = defineEmits(['jump', 'toggle-overview'])

// 只渲染当前图附近的缩略图：图库可能有上千张，全量渲染会明显卡顿
const TRACK_RADIUS = 24

const trackEl = ref(null)

const windowItems = computed(() => {
  const list = props.images || []
  const total = list.length
  if (!total) return []

  const span = TRACK_RADIUS * 2 + 1
  if (total <= span) return list.map((img, index) => ({ img, index }))

  let start = Math.max(0, props.currentIndex - TRACK_RADIUS)
  let end = start + span
  if (end > total) {
    end = total
    start = end - span
  }

  return list.slice(start, end).map((img, offset) => ({ img, index: start + offset }))
})

const scrollActiveIntoView = async () => {
  await nextTick()
  const root = trackEl.value
  if (!root) return
  const active = root.querySelector(`[data-index="${props.currentIndex}"]`)
  if (!active) return
  const target = active.offsetLeft - root.clientWidth / 2 + active.clientWidth / 2
  root.scrollTo({ left: Math.max(0, target), behavior: 'smooth' })
}

watch(() => props.currentIndex, scrollActiveIntoView)
onMounted(scrollActiveIntoView)
</script>

<template>
  <div
    v-if="images.length > 1"
    class="absolute bottom-0 left-0 right-0 z-[58] flex h-14 items-center gap-2 border-t border-white/10 bg-black/55 px-3 backdrop-blur-xl"
    @click.stop
    @wheel.stop
  >
    <button
      type="button"
      class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-white/50 transition-colors hover:bg-white/10 hover:text-white disabled:pointer-events-none disabled:opacity-25"
      :disabled="currentIndex <= 0"
      title="上一张"
      @click="emit('jump', currentIndex - 1)"
    >
      <ChevronLeft class="h-4 w-4" />
    </button>

    <div ref="trackEl" class="scrollbar-hide flex min-w-0 flex-1 items-center gap-1.5 overflow-x-auto">
      <button
        v-for="entry in windowItems"
        :key="entry.img.relPath"
        :data-index="entry.index"
        type="button"
        class="h-10 w-10 shrink-0 overflow-hidden rounded transition-all duration-150"
        :class="entry.index === currentIndex ? 'ring-2 ring-blue-400' : 'opacity-50 hover:opacity-100'"
        :title="entry.img.name"
        @click="emit('jump', entry.index)"
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

    <button
      type="button"
      class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-white/50 transition-colors hover:bg-white/10 hover:text-white disabled:pointer-events-none disabled:opacity-25"
      :disabled="currentIndex >= images.length - 1"
      title="下一张"
      @click="emit('jump', currentIndex + 1)"
    >
      <ChevronRight class="h-4 w-4" />
    </button>

    <button
      type="button"
      class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-white/50 transition-colors hover:bg-white/10 hover:text-white"
      title="批次概览（空格）"
      @click="emit('toggle-overview')"
    >
      <LayoutGrid class="h-4 w-4" />
    </button>
  </div>
</template>
