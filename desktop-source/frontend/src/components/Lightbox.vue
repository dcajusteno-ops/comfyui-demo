<script setup>
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import * as App from '@/api'
import AmbientBackdrop from './AmbientBackdrop.vue'
import BatchOverview from './BatchOverview.vue'
import FavoriteGroupsDialog from './FavoriteGroupsDialog.vue'
import ImageMetadataPanel from './ImageMetadataPanel.vue'
import LightboxToolbar from './LightboxToolbar.vue'
import LightboxViewer from './LightboxViewer.vue'

const props = defineProps({
  image: Object,
  isOpen: Boolean,
  images: { type: Array, default: () => [] },
  currentIndex: Number,
  favoriteGroups: { type: Array, default: () => [] },
  tags: { type: Array, default: () => [] },
  imageTags: { type: Object, default: () => ({}) },
  imageNotes: { type: Object, default: () => ({}) },
  openTagsOnMount: { type: Boolean, default: false },
})

const emit = defineEmits(['close', 'navigate', 'toggle-favorite', 'add-tag', 'remove-tag', 'delete', 'open-location', 'favorite-groups-changed', 'open-prompt-assistant', 'jump-to-index'])

const currentStackImage = ref(null)

const currentDisplayImage = computed(() => currentStackImage.value || props.image)

watch(() => props.image, () => {
  currentStackImage.value = null
})

const scale = ref(1)
const offset = ref({ x: 0, y: 0 })
const isDragging = ref(false)
const lastMousePos = ref({ x: 0, y: 0 })
const metadata = ref(null)
const metadataLoading = ref(false)
const metadataError = ref('')
const favoriteGroupsDialogOpen = ref(false)
const hudVisible = ref(localStorage.getItem('lightboxHudVisible') !== 'false')
const batchOverlayOpen = ref(false)

watch(
  () => props.isOpen,
  isOpen => {
    if (!isOpen) batchOverlayOpen.value = false
  },
)
const displayImageSrc = ref('')
const fullImageLoading = ref(false)
let metadataRequestId = 0
let imageRequestId = 0
const lightboxPreviewSrc = computed(
  () => currentDisplayImage.value?.previewPath || currentDisplayImage.value?.thumbPath || currentDisplayImage.value?.path || '',
)
const lightboxFullSrc = computed(() => currentDisplayImage.value?.path || '')

const syncLightboxImageSource = () => {
  const requestId = ++imageRequestId
  const previewSrc = lightboxPreviewSrc.value
  const fullSrc = lightboxFullSrc.value

  if (!props.isOpen || !fullSrc) {
    displayImageSrc.value = ''
    fullImageLoading.value = false
    return
  }

  if (!previewSrc || previewSrc === fullSrc) {
    displayImageSrc.value = fullSrc
    fullImageLoading.value = false
    return
  }

  displayImageSrc.value = previewSrc
  fullImageLoading.value = true

  const loader = new Image()
  loader.decoding = 'async'
  loader.onload = () => {
    if (requestId !== imageRequestId) return
    displayImageSrc.value = fullSrc
    fullImageLoading.value = false
  }
  loader.onerror = () => {
    if (requestId !== imageRequestId) return
    displayImageSrc.value = fullSrc
    fullImageLoading.value = false
  }
  loader.src = fullSrc
}

const totalImages = computed(() => props.images?.length || 0)
const canGoPrev = computed(() => props.currentIndex > 0)
const canGoNext = computed(() => props.currentIndex < totalImages.value - 1)
const imageCounter = computed(() => `${Math.min((props.currentIndex || 0) + 1, Math.max(totalImages.value, 1))} / ${Math.max(totalImages.value, 1)}`)

const imageTags = computed(() => {
  if (!currentDisplayImage.value) return []
  const tagIds = props.imageTags[currentDisplayImage.value.relPath] || []
  return tagIds.map(id => props.tags.find(tag => tag.id === id)).filter(Boolean)
})

const availableTags = computed(() => {
  const assigned = imageTags.value.map(tag => tag.id)
  return props.tags.filter(tag => !assigned.includes(tag.id))
})

const resetZoom = () => {
  scale.value = 1
  offset.value = { x: 0, y: 0 }
  isDragging.value = false
}

const toggleHud = () => {
  hudVisible.value = !hudVisible.value
  try {
    localStorage.setItem('lightboxHudVisible', String(hudVisible.value))
  } catch {
    // 存储不可用时忽略，不影响本次会话
  }
}

const handleBatchJump = index => {
  batchOverlayOpen.value = false
  emit('jump-to-index', index)
}

const loadImageMetadata = async () => {
  const relPath = currentDisplayImage.value?.relPath
  if (!props.isOpen || !relPath) {
    metadata.value = null
    metadataLoading.value = false
    metadataError.value = ''
    return
  }

  const requestId = ++metadataRequestId
  metadataLoading.value = true
  metadataError.value = ''

  try {
    const result = await App.GetImageMetadata(relPath)
    if (requestId !== metadataRequestId) return
    metadata.value = result
  } catch (error) {
    if (requestId !== metadataRequestId) return
    metadata.value = null
    metadataError.value = error?.message || `${error}`
  } finally {
    if (requestId === metadataRequestId) {
      metadataLoading.value = false
    }
  }
}

const goToPrev = () => {
  if (!canGoPrev.value) return
  resetZoom()
  emit('navigate', 'prev')
}

const goToNext = () => {
  if (!canGoNext.value) return
  resetZoom()
  emit('navigate', 'next')
}

const goToPrevStackItem = () => {
  if (!props.image.isStackPrimary || props.image.stackCount <= 1) return
  const stackItems = [props.image, ...props.image.stackChildren]
  const currentIndex = stackItems.findIndex(i => i.relPath === currentDisplayImage.value.relPath)
  if (currentIndex > 0) {
    currentStackImage.value = stackItems[currentIndex - 1]
  } else {
    // Wrap around
    currentStackImage.value = stackItems[stackItems.length - 1]
  }
}

const goToNextStackItem = () => {
  if (!props.image.isStackPrimary || props.image.stackCount <= 1) return
  const stackItems = [props.image, ...props.image.stackChildren]
  const currentIndex = stackItems.findIndex(i => i.relPath === currentDisplayImage.value.relPath)
  if (currentIndex < stackItems.length - 1) {
    currentStackImage.value = stackItems[currentIndex + 1]
  } else {
    // Wrap around
    currentStackImage.value = stackItems[0]
  }
}

const getStackCurrentIndex = () => {
  if (!props.image?.isStackPrimary) return 1
  const stackItems = [props.image, ...(props.image.stackChildren || [])]
  const idx = stackItems.findIndex(i => i.relPath === currentDisplayImage.value?.relPath)
  return idx >= 0 ? idx + 1 : 1
}

const goToStackIndex = index => {
  if (!props.image?.isStackPrimary || props.image.stackCount <= 1) return
  const stackItems = [props.image, ...(props.image.stackChildren || [])]
  if (index < 0 || index >= stackItems.length) return
  currentStackImage.value = stackItems[index]
}

const handleWheel = event => {
  if (!props.isOpen) return
  const delta = -event.deltaY
  const zoomFactor = 1.1
  const nextScale = delta > 0 ? scale.value * zoomFactor : scale.value / zoomFactor
  scale.value = Math.min(Math.max(nextScale, 0.5), 10)
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

const handleKey = event => {
  if (!props.isOpen) return

  // 空格：开合批次概览层（原先未占用该键）
  if (event.code === 'Space') {
    event.preventDefault()
    batchOverlayOpen.value = !batchOverlayOpen.value
    return
  }

  // 概览层打开时 Esc 优先关闭它
  if (event.key === 'Escape' && batchOverlayOpen.value) {
    batchOverlayOpen.value = false
    return
  }

  if (event.key === 'Escape') emit('close')
  if (event.key === 'ArrowLeft') {
    goToPrev()
  }
  if (event.key === 'ArrowRight') {
    goToNext()
  }
  if (event.key === 'ArrowUp') {
    if (props.image?.isStackPrimary && props.image.stackCount > 1) {
      goToPrevStackItem()
    }
  }
  if (event.key === 'ArrowDown') {
    if (props.image?.isStackPrimary && props.image.stackCount > 1) {
      goToNextStackItem()
    }
  }
  if (event.key === '0' && event.ctrlKey) resetZoom()
}

watch(
  () => [props.isOpen, currentDisplayImage.value?.relPath],
  () => {
    void loadImageMetadata()
  },
  { immediate: true }
)

watch(
  () => [props.isOpen, currentDisplayImage.value?.relPath, currentDisplayImage.value?.previewPath, currentDisplayImage.value?.path],
  () => {
    syncLightboxImageSource()
  },
  { immediate: true }
)

onMounted(() => {
  window.addEventListener('keydown', handleKey)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKey)
})
</script>

<template>
  <transition
    enter-active-class="transition duration-200 ease-out"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition duration-200 ease-in"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="isOpen && image"
      class="fixed inset-0 z-[50] flex items-center justify-center overflow-hidden bg-black"
      @click.self="$emit('close')"
    >
      <div class="absolute inset-0 z-0 overflow-hidden pointer-events-none">
        <transition
          enter-active-class="transition-opacity duration-700 ease-in-out"
          leave-active-class="transition-opacity duration-700 ease-in-out"
          enter-from-class="opacity-0"
          leave-to-class="opacity-0"
        >
          <img
            :key="displayImageSrc || currentDisplayImage.path"
            :src="displayImageSrc || currentDisplayImage.path"
            class="h-full w-full scale-110 object-cover opacity-40 blur-[60px]"
            loading="eager"
            decoding="async"
          />
        </transition>
        <div class="absolute inset-0 bg-black/40 backdrop-blur-[2px]" />
      </div>

      <AmbientBackdrop
        :image-src="displayImageSrc || currentDisplayImage.previewPath || currentDisplayImage.thumbPath || currentDisplayImage.path || ''"
      />

      <LightboxToolbar
        :is-open="isOpen"
        :open-tags-on-mount="openTagsOnMount"
        :current-display-image="currentDisplayImage"
        :image-tags="imageTags"
        :available-tags="availableTags"
        :metadata="metadata"
        @remove-tag="$emit('remove-tag', ...$event)"
        @add-tag="$emit('add-tag', ...$event)"
        @open-location="$emit('open-location', $event)"
        @open-favorite-groups="favoriteGroupsDialogOpen = true"
        @delete="$emit('delete', $event)"
        @toggle-favorite="$emit('toggle-favorite', $event)"
        @open-prompt-assistant="$emit('open-prompt-assistant', $event)"
        @close="$emit('close')"
      />

      <LightboxViewer
        :image="image"
        :current-display-image="currentDisplayImage"
        :can-go-prev="canGoPrev"
        :can-go-next="canGoNext"
        :display-image-src="displayImageSrc"
        :full-image-loading="fullImageLoading"
        :scale="scale"
        :offset="offset"
        :is-dragging="isDragging"
        :image-counter="imageCounter"
        :stack-current-index="getStackCurrentIndex()"
        :metadata="metadata"
        :hud-visible="hudVisible"
        :images="images"
        :current-index="currentIndex || 0"
        @prev="goToPrev"
        @next="goToNext"
        @prev-stack="goToPrevStackItem"
        @next-stack="goToNextStackItem"
        @select-stack="goToStackIndex"
        @jump-to-index="handleBatchJump"
        @toggle-overview="batchOverlayOpen = !batchOverlayOpen"
        @reset-zoom="resetZoom"
        @toggle-hud="toggleHud"
        @viewer-wheel="handleWheel"
        @viewer-mousedown="handleMouseDown"
        @viewer-mousemove="handleMouseMove"
        @viewer-mouseup="handleMouseUp"
      />

      <ImageMetadataPanel
        :current-display-image="currentDisplayImage"
        :image-notes="imageNotes"
        :metadata="metadata"
        :metadata-loading="metadataLoading"
        :metadata-error="metadataError"
      />

      <FavoriteGroupsDialog
        v-model:open="favoriteGroupsDialogOpen"
        :groups="favoriteGroups"
        :image="currentDisplayImage"
        @change="$emit('favorite-groups-changed')"
      />

      <BatchOverview
        v-if="batchOverlayOpen"
        :images="images"
        :current-image="currentDisplayImage"
        :current-index="currentIndex || 0"
        @jump="handleBatchJump"
        @close="batchOverlayOpen = false"
      />
    </div>
  </transition>
</template>


