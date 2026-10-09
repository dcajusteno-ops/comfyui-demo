<script setup>
import { computed, ref, watch } from 'vue'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Slider } from '@/components/ui/slider'
import { Switch } from '@/components/ui/switch'
import { ImageIcon, Loader2, RotateCcw, Trash2, Upload } from 'lucide-vue-next'
import {
  clearWallpaperImage,
  isWallpaperActive,
  resetWallpaperAppearance,
  selectWallpaperImage,
  updateWallpaper,
  wallpaper,
  wallpaperUrl,
} from '@/wallpaper'

const busy = ref(false)
const previewFailed = ref(false)

const fitOptions = [
  { value: 'cover', label: '铺满' },
  { value: 'contain', label: '完整' },
  { value: 'auto', label: '原始' },
]

const hasImage = computed(() => wallpaper.value.imagePath !== '')
const dimPercent = computed(() => Math.round(wallpaper.value.dim * 100))
const surfacePercent = computed(() => Math.round(wallpaper.value.surfaceAlpha * 100))

// reka-ui 的单滑块在不同版本里可能回传数组，这里统一取数字
const toNumber = (value) => (Array.isArray(value) ? Number(value[0]) : Number(value))

watch(wallpaperUrl, () => {
  previewFailed.value = false
})

const handleSelect = async () => {
  busy.value = true
  try {
    const config = await selectWallpaperImage()
    if (config?.imagePath) {
      toast.success('壁纸已更新')
    }
  } catch (error) {
    toast.error(`选择壁纸失败: ${error}`)
  } finally {
    busy.value = false
  }
}

const handleClear = async () => {
  busy.value = true
  try {
    await clearWallpaperImage()
    toast.success('壁纸已移除')
  } catch (error) {
    toast.error(`移除失败: ${error}`)
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <Card class="rounded-[28px] border-border/70 bg-card shadow-none">
    <CardHeader>
      <CardTitle class="text-lg">背景壁纸</CardTitle>
      <CardDescription>给整个界面铺一张底图，侧栏、卡片和浮层会半透明地透出它。</CardDescription>
    </CardHeader>
    <CardContent class="grid gap-5">
      <div class="overflow-hidden rounded-[24px] border border-border bg-background">
        <div v-if="hasImage && !previewFailed" class="relative aspect-[16/10] overflow-hidden bg-muted/30">
          <img
            :src="wallpaperUrl"
            alt="壁纸预览"
            class="h-full w-full object-cover"
            @error="previewFailed = true"
          />
          <div
            v-if="!isWallpaperActive"
            class="absolute inset-0 flex items-center justify-center bg-background/70 text-xs text-muted-foreground"
          >
            已停用，点下方开关启用
          </div>
        </div>
        <div v-else class="flex aspect-[16/10] items-center justify-center bg-muted/30 text-muted-foreground">
          <div class="flex flex-col items-center gap-2">
            <ImageIcon class="h-6 w-6" />
            <span class="text-sm">{{ previewFailed ? '图片读取失败，请重新选择' : '还没有设置壁纸' }}</span>
          </div>
        </div>
      </div>

      <div class="grid gap-2">
        <Button class="h-11 rounded-2xl" :disabled="busy" @click="handleSelect">
          <Loader2 v-if="busy" class="mr-2 h-4 w-4 animate-spin" />
          <Upload v-else class="mr-2 h-4 w-4" />
          {{ busy ? '处理中...' : (hasImage ? '更换壁纸' : '选择壁纸图片') }}
        </Button>
        <Button
          variant="outline"
          class="h-11 rounded-2xl shadow-none"
          :disabled="busy || !hasImage"
          @click="handleClear"
        >
          <Trash2 class="mr-2 h-4 w-4" />
          移除壁纸
        </Button>
      </div>

      <p class="text-xs leading-5 text-muted-foreground">
        图片会复制到应用数据目录（data/wallpaper），原图移动或删除都不影响显示。
      </p>

      <div class="grid gap-4 border-t border-border pt-4">
        <div class="flex items-center justify-between gap-3">
          <div class="min-w-0">
            <div class="text-sm font-medium">启用壁纸</div>
            <div class="mt-1 text-xs text-muted-foreground">关闭后立即恢复不透明界面，图片仍保留。</div>
          </div>
          <Switch
            :model-value="wallpaper.enabled"
            :disabled="!hasImage"
            @update:model-value="updateWallpaper({ enabled: $event })"
          />
        </div>

        <div class="grid gap-2">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium">填充方式</span>
            <span class="text-xs text-muted-foreground">
              {{ fitOptions.find((item) => item.value === wallpaper.fit)?.label }}
            </span>
          </div>
          <div class="flex items-center gap-1 rounded-2xl border border-border p-1">
            <button
              v-for="option in fitOptions"
              :key="option.value"
              type="button"
              class="flex-1 rounded-xl px-2 py-1.5 text-xs transition-colors disabled:opacity-40"
              :class="
                wallpaper.fit === option.value
                  ? 'bg-secondary text-foreground'
                  : 'text-muted-foreground hover:bg-muted/60'
              "
              :disabled="!hasImage"
              @click="updateWallpaper({ fit: option.value })"
            >
              {{ option.label }}
            </button>
          </div>
        </div>

        <div class="grid gap-2">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium">压暗</span>
            <span class="text-xs text-muted-foreground">{{ dimPercent }}%</span>
          </div>
          <Slider
            :model-value="[wallpaper.dim]"
            :min="0"
            :max="1"
            :step="0.05"
            :disabled="!hasImage"
            @update:model-value="updateWallpaper({ dim: toNumber($event) })"
          />
        </div>

        <div class="grid gap-2">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium">面板不透明度</span>
            <span class="text-xs text-muted-foreground">{{ surfacePercent }}%</span>
          </div>
          <Slider
            :model-value="[wallpaper.surfaceAlpha]"
            :min="0.3"
            :max="1"
            :step="0.02"
            :disabled="!hasImage"
            @update:model-value="updateWallpaper({ surfaceAlpha: toNumber($event) })"
          />
          <p class="text-xs text-muted-foreground">调低面板会更透、壁纸更明显；调高则界面更实。</p>
        </div>

        <div class="flex items-center justify-between gap-3">
          <div class="min-w-0">
            <div class="text-sm font-medium">面板毛玻璃</div>
            <div class="mt-1 text-xs text-muted-foreground">给侧栏与浮层加一层背景模糊。</div>
          </div>
          <Switch
            :model-value="wallpaper.blur"
            :disabled="!hasImage"
            @update:model-value="updateWallpaper({ blur: $event })"
          />
        </div>

        <Button
          variant="outline"
          class="h-10 rounded-2xl shadow-none"
          :disabled="!hasImage"
          @click="resetWallpaperAppearance"
        >
          <RotateCcw class="mr-2 h-4 w-4" />
          恢复默认外观
        </Button>
      </div>
    </CardContent>
  </Card>
</template>
