<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { cn } from '@/lib/utils'
import { Check, ChevronDown, Search } from 'lucide-vue-next'

const props = defineProps({
  id: { type: String, default: undefined },
  modelValue: { type: [String, Number, null], default: '' },
  options: { type: Array, default: () => [] },
  placeholder: { type: String, default: '请选择' },
  searchable: { type: Boolean, default: false },
  searchPlaceholder: { type: String, default: '搜索...' },
  emptyText: { type: String, default: '没有匹配项' },
  title: { type: String, default: undefined },
  disabled: { type: Boolean, default: false },
  triggerClass: { type: [String, Array, Object], default: '' },
  contentClass: { type: [String, Array, Object], default: '' },
  maxHeight: { type: String, default: '18rem' },
})

const emit = defineEmits(['update:modelValue', 'change'])

const open = ref(false)
const query = ref('')
const activeIndex = ref(-1)
const triggerRef = ref(null)
const searchRef = ref(null)
const listRef = ref(null)
const contentStyle = ref({})

const normalizedOptions = computed(() =>
  (props.options || []).map((option) => ({
    value: option.value,
    label: option.label ?? String(option.value ?? ''),
    count: option.count,
    disabled: !!option.disabled,
  })),
)

const selectedOption = computed(() =>
  normalizedOptions.value.find((option) => option.value === props.modelValue),
)

const displayLabel = computed(() => {
  if (selectedOption.value) return selectedOption.value.label
  if (props.modelValue === '' || props.modelValue === null || props.modelValue === undefined) {
    return props.placeholder
  }
  return String(props.modelValue)
})

const filteredOptions = computed(() => {
  const keyword = query.value.trim().toLowerCase()
  if (!keyword) return normalizedOptions.value
  return normalizedOptions.value.filter(
    (option) =>
      option.label.toLowerCase().includes(keyword) ||
      String(option.value ?? '').toLowerCase().includes(keyword),
  )
})

watch(filteredOptions, (list) => {
  if (!list.length) {
    activeIndex.value = -1
    return
  }
  if (activeIndex.value < 0 || activeIndex.value >= list.length) {
    activeIndex.value = 0
  }
})

const measureContentWidth = () => {
  const width = triggerRef.value?.getBoundingClientRect().width || 200
  contentStyle.value = { width: `${Math.max(width, 160)}px` }
}

const scrollActiveIntoView = () => {
  const active = listRef.value?.querySelector('[data-active="true"]')
  active?.scrollIntoView({ block: 'nearest' })
}

const setOpen = (value) => {
  open.value = value
  if (!value) return

  query.value = ''
  measureContentWidth()

  const selectedIndex = normalizedOptions.value.findIndex(
    (option) => option.value === props.modelValue,
  )
  activeIndex.value = selectedIndex >= 0 ? selectedIndex : 0

  nextTick(() => {
    measureContentWidth()
    if (props.searchable) {
      searchRef.value?.focus()
    }
    scrollActiveIntoView()
  })
}

const moveActive = (delta) => {
  const list = filteredOptions.value
  if (!list.length) return

  let next = activeIndex.value
  for (let step = 0; step < list.length; step += 1) {
    next = (next + delta + list.length) % list.length
    if (!list[next].disabled) break
  }
  activeIndex.value = next
  nextTick(scrollActiveIntoView)
}

const selectOption = (option) => {
  if (!option || option.disabled) return
  emit('update:modelValue', option.value)
  emit('change', option.value)
  setOpen(false)
}

const handleKeydown = (event) => {
  if (props.disabled) return

  switch (event.key) {
    case 'ArrowDown':
    case 'ArrowUp':
      event.preventDefault()
      if (!open.value) {
        setOpen(true)
        return
      }
      moveActive(event.key === 'ArrowDown' ? 1 : -1)
      break
    case 'Enter':
      if (!open.value) {
        event.preventDefault()
        setOpen(true)
        return
      }
      event.preventDefault()
      selectOption(filteredOptions.value[activeIndex.value])
      break
    case ' ':
      if (!open.value) {
        event.preventDefault()
        setOpen(true)
      }
      break
    case 'Escape':
      if (open.value) {
        event.preventDefault()
        setOpen(false)
      }
      break
    default:
      break
  }
}
</script>

<template>
  <Popover :open="open" @update:open="setOpen">
    <PopoverTrigger as-child>
      <button
        ref="triggerRef"
        :id="id"
        type="button"
        :disabled="disabled"
        :title="title"
        role="combobox"
        aria-haspopup="listbox"
        :aria-expanded="open"
        :class="
          cn(
            'flex items-center justify-between gap-2 border-control-border text-left text-foreground outline-none transition focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/60 focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:cursor-not-allowed disabled:opacity-60',
            triggerClass,
          )
        "
        @keydown="handleKeydown"
      >
        <span
          class="min-w-0 flex-1 truncate"
          :class="{ 'text-muted-foreground': !selectedOption }"
          :title="displayLabel"
        >
          {{ displayLabel }}
        </span>
        <ChevronDown
          class="h-4 w-4 shrink-0 opacity-60 transition-transform duration-200"
          :class="{ 'rotate-180': open }"
        />
      </button>
    </PopoverTrigger>

    <PopoverContent
      data-slot="app-select-content"
      side="bottom"
      align="start"
      :side-offset="6"
      avoid-collisions
      :style="contentStyle"
      :class="cn('z-50 rounded-xl border border-border bg-popover p-1.5 shadow-lg', contentClass)"
      @open-auto-focus="$event.preventDefault()"
    >
      <div v-if="searchable" class="p-0.5 pb-1">
        <div class="relative">
          <Search class="pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
          <input
            ref="searchRef"
            v-model="query"
            type="text"
            :placeholder="searchPlaceholder"
            class="h-8 w-full rounded-md border border-control-border/70 bg-background pl-8 pr-2 text-xs text-foreground outline-none transition focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/50 focus-visible:ring-inset"
            @keydown="handleKeydown"
          />
        </div>
      </div>

      <div
        ref="listRef"
        role="listbox"
        class="custom-scrollbar overflow-y-auto p-0.5"
        :style="{ maxHeight }"
      >
        <button
          v-for="(option, index) in filteredOptions"
          :key="String(option.value)"
          type="button"
          role="option"
          :aria-selected="option.value === modelValue"
          :data-active="index === activeIndex"
          :disabled="option.disabled"
          class="flex w-full items-center gap-2 rounded-lg px-2.5 py-1.5 text-left text-sm transition-colors hover:bg-accent hover:text-accent-foreground disabled:opacity-50"
          :class="{ 'bg-accent text-accent-foreground': index === activeIndex }"
          @click="selectOption(option)"
          @mousemove="activeIndex = index"
        >
          <span class="min-w-0 flex-1 truncate" :title="option.label">{{ option.label }}</span>
          <span
            v-if="option.count !== undefined"
            class="shrink-0 text-[11px] tabular-nums text-muted-foreground"
          >
            {{ option.count }}
          </span>
          <Check v-if="option.value === modelValue" class="h-3.5 w-3.5 shrink-0 text-primary" />
        </button>

        <div v-if="!filteredOptions.length" class="px-3 py-4 text-center text-xs text-muted-foreground">
          {{ emptyText }}
        </div>
      </div>
    </PopoverContent>
  </Popover>
</template>
