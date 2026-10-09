<script setup>
import { X } from 'lucide-vue-next'
import { BATCH_MODES, useBatchGroups } from '@/composables/useBatchGroups'

const props = defineProps({
  images: { type: Array, default: () => [] },
  currentImage: { type: Object, default: null },
  currentIndex: { type: Number, default: 0 },
})

const emit = defineEmits(['jump', 'close'])

const { mode, group } = useBatchGroups(
  () => props.images,
  () => props.currentImage,
)
</script>

<template>
  <div
    class="absolute inset-0 z-[80] flex flex-col bg-black/85 backdrop-blur-xl"
    @click.self="emit('close')"
  >
    <div class="flex shrink-0 items-center gap-3 px-6 py-4">
      <span class="shrink-0 text-sm font-medium text-white/85">批次概览</span>
      <span v-if="group" class="truncate font-mono text-xs text-white/45">
        {{ group.label }} · {{ group.items.length }} 张
      </span>

      <div class="ml-auto flex shrink-0 items-center gap-0.5 rounded-full border border-white/10 bg-white/5 p-0.5">
        <button
          v-for="option in BATCH_MODES"
          :key="option.value"
          type="button"
          class="rounded-full px-3 py-1 text-xs transition-colors"
          :class="mode === option.value ? 'bg-white/15 text-white' : 'text-white/50 hover:text-white/85'"
          @click="mode = option.value"
        >
          {{ option.label }}
        </button>
      </div>

      <button
        type="button"
        class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-white/50 transition-colors hover:bg-white/10 hover:text-white"
        title="关闭（Esc 或空格）"
        @click="emit('close')"
      >
        <X class="h-4 w-4" />
      </button>
    </div>

    <div class="custom-scrollbar min-h-0 flex-1 overflow-y-auto px-6 pb-6" @click.self="emit('close')">
      <div
        v-if="group && group.items.length"
        class="grid gap-2"
        style="grid-template-columns: repeat(auto-fill, minmax(112px, 1fr))"
      >
        <button
          v-for="entry in group.items"
          :key="entry.img.relPath"
          type="button"
          class="group relative aspect-square overflow-hidden rounded-lg border transition-all duration-150"
          :class="
            entry.index === currentIndex
              ? 'border-blue-400 ring-2 ring-blue-400/40'
              : 'border-white/10 hover:border-white/40'
          "
          :title="entry.img.name"
          @click="emit('jump', entry.index)"
        >
          <img
            :src="entry.img.thumbPath || entry.img.previewPath || entry.img.path"
            alt=""
            loading="lazy"
            decoding="async"
            draggable="false"
            class="h-full w-full object-cover transition-transform duration-200 group-hover:scale-105"
          />
          <span
            class="pointer-events-none absolute inset-x-0 bottom-0 truncate bg-black/70 px-1.5 py-1 text-[10px] text-white/80 opacity-0 transition-opacity group-hover:opacity-100"
          >
            {{ entry.img.name }}
          </span>
        </button>
      </div>

      <div v-else class="flex h-full items-center justify-center text-sm text-white/40">
        当前批次没有可展示的图片
      </div>
    </div>
  </div>
</template>
