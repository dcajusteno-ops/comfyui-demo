<script setup>
defineProps({
  imageSrc: { type: String, default: '' },
})

const maskStyle = {
  maskImage: 'linear-gradient(to bottom, rgba(0, 0, 0, 0.72), transparent 92%)',
  WebkitMaskImage: 'linear-gradient(to bottom, rgba(0, 0, 0, 0.72), transparent 92%)',
}
</script>

<template>
  <div
    v-if="imageSrc"
    class="pointer-events-none absolute inset-x-0 top-full h-16 overflow-hidden opacity-70"
    :style="maskStyle"
  >
    <img
      :src="imageSrc"
      alt=""
      aria-hidden="true"
      draggable="false"
      class="h-full w-full -scale-y-100 object-cover blur-[3px]"
    />
  </div>
</template>
