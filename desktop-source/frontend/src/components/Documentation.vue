<script setup>
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import {
  BookOpen,
  CalendarDays,
  Clock3,
  Filter,
  FolderTree,
  Heart,
  HelpCircle,
  Images,
  Keyboard,
  LayoutDashboard,
  Orbit,
  PencilRuler,
  Search,
  Settings2,
  Tag,
  Trash2,
  Wallpaper,
  Wand2,
} from 'lucide-vue-next'

const featureCards = [
  {
    icon: LayoutDashboard,
    title: '工作台总览',
    description: '打开软件后的主入口，用于快速确认最新作品、今日新增、图片总数和当前存储占用。',
  },
  {
    icon: CalendarDays,
    title: '日期产出工作台',
    description: '围绕今天、昨天、最近 7 天、本月和自定义日期范围回看近期创作，非常适合阶段性复盘。',
  },
  {
    icon: Filter,
    title: '模型 / LoRA 筛选',
    description: '自动从 PNG 元数据里提取模型与 LoRA 信息，用于图库、工作台和回看页的联合筛选。',
  },
  {
    icon: Search,
    title: '搜索与精查',
    description: '支持按文件名、路径、Prompt、模型、LoRA、标签和笔记搜索，适合在大图库中快速定位目标。',
  },
  {
    icon: Images,
    title: '大图查看器',
    description: '氛围光晕与地面倒影、底部缩略图导航条、批次概览层、参数 HUD、3D 层叠连拍组，看图细节一次补齐。',
  },
  {
    icon: Orbit,
    title: '3D 环绕预览',
    description: '图库可切换「网格 / 环绕」视图：图片沿中央光柱多层环绕，拖动旋转、滚轮上下滚动，每圈张数可调。',
  },
  {
    icon: Wallpaper,
    title: '背景壁纸自定义',
    description: '个人中心选一张本地图片当背景，可调压暗、面板不透明度和毛玻璃，侧栏与卡片会半透明地透出壁纸。',
  },
  {
    icon: Clock3,
    title: '大型图库性能模式',
    description: '当当前目录图片很多时，软件会倾向启用轻量分页、预览变体与更稳的加载策略，减少卡顿。',
  },
  {
    icon: FolderTree,
    title: '灵动图库小窗',
    description: '把最新 output 放到右下角小窗中查看，支持置顶、目录切换、搜索、分页、批量选择和快速恢复主窗口。',
  },
  {
    icon: FolderTree,
    title: '目录管理',
    description: '支持默认目录、日期归档目录、自定义目录并行浏览，也支持重新绑定任意 ComfyUI output。',
  },
  {
    icon: Heart,
    title: '收藏与分组',
    description: '把满意作品和常用内容加入收藏，再按分组做长期整理，方便从全量 output 中抽离重点结果。',
  },
  {
    icon: Tag,
    title: '标签与笔记',
    description: '支持自定义标签、颜色和图片笔记，适合记录题材、角色、风格、返工计划和筛图依据。',
  },
  {
    icon: PencilRuler,
    title: '提示词助手与模板',
    description: '围绕本地词库完成搜索、筛选、拼装、模板保存和重复使用，降低 Prompt 重复劳动。',
  },
  {
    icon: Wand2,
    title: '自动规则引擎',
    description: '按模型、LoRA、Prompt 或文件名自动打标签、加入收藏、归类目录，减少手工整理成本。',
  },
  {
    icon: Settings2,
    title: '设置中心',
    description: '统一管理快捷键、性能模式、工具菜单、缓存治理与目录维护，是长期使用时的重要入口。',
  },
  {
    icon: Trash2,
    title: '回收站保护',
    description: '删除图片默认先进回收站，可恢复、可批量清空，也可以按保留天数自动清理。',
  },
]

const quickGuides = [
  {
    title: '第一次进入怎么设置',
    steps: [
      '首次进入时，先绑定 ComfyUI 的真实 output 文件夹，而不是选择软件 exe 所在目录。',
      '绑定成功后，软件会自动识别 output 的上一级目录作为根目录，并刷新默认目录、日期归档和自定义目录。',
      '建议先看一遍工作台总览，确认图片总数、最新作品和存储占用都正常。',
    ],
  },
  {
    title: '按日期回看最近作品',
    steps: [
      '打开左侧“日期产出”页面，优先使用今天、昨天、最近 7 天或本月这些常用范围。',
      '需要更细时，可以继续叠加模型筛选或 LoRA 筛选。',
      '如果要回到更大范围继续细筛，可以再跳回图库页面继续操作。',
    ],
  },
  {
    title: '把满意作品沉淀下来',
    steps: [
      '在图库或大图查看里，把值得保留的图加入收藏。',
      '再为图片补上标签、笔记或收藏分组，形成后续回看和整理入口。',
      '这样可以把“全量 output” 与“长期保留结果”分层管理。',
    ],
  },
  {
    title: '用提示词助手减少重复劳动',
    steps: [
      '打开提示词助手或提示词模板入口，搜索目标词条或模板。',
      '把合适的内容加入正向 / 反向 Prompt 区域，再保存成模板。',
      '以后需要重复同类风格时，直接复用模板即可。',
    ],
  },
  {
    title: '边生成边看最新 output',
    steps: [
      '从侧边栏打开灵动图库小窗，需要时开启置顶。',
      '小窗顶部保留目录和搜索，筛选、排序、Output、清缓存收在工具按钮里。',
      '使用分页浏览更多图片，每页可在 60 / 120 / 240 张之间切换。',
    ],
  },
  {
    title: '在大图里快速过一批图',
    steps: [
      '点开任意一张图，想跳到附近的图就直接点底部缩略图导航条。',
      '想整体看一眼这一批，按空格打开批次概览层，还可以切换“同日 / 同模型 / 同提示词 / 连拍组”。',
      '点中目标图后概览层自动收起，回到大图继续看细节。',
    ],
  },
  {
    title: '用环绕模式逛一遍图库',
    steps: [
      '在图库工具栏把视图从“网格”切到“环绕”。',
      '拖动旋转看整体分布，滚轮上下滚动看塔的上下层，空闲时会自动缓慢旋转。',
      '左上角调整每圈张数：张数少塔就高，可以滚着逛；调大则更紧凑。',
    ],
  },
  {
    title: '换一张背景壁纸',
    steps: [
      '进入个人中心，在“外观与偏好”下面找到“背景壁纸”。',
      '点“选择壁纸图片”，图片会复制到 data/wallpaper，原图之后删掉也不影响。',
      '按需调压暗、面板不透明度；嫌面板太实就调低不透明度，想加模糊就打开毛玻璃。',
    ],
  },
]

const shortcuts = [
  { key: 'Esc', action: '关闭大图、弹窗或退出当前选择状态' },
  { key: 'Delete', action: '删除当前选中图片，删除前通常会先弹确认' },
  { key: '方向键', action: '在大图预览中切换上一张 / 下一张' },
  { key: '↑ / ↓', action: '在大图中切换同一连拍组里的上一张 / 下一张' },
  { key: '空格', action: '在大图中开合批次概览层' },
  { key: 'Ctrl + 0', action: '在大图中重置缩放状态' },
  { key: '滚轮 / 拖动', action: '在大图中缩放与平移（双击可快速放大 / 还原）' },
  { key: '批量模式 + 点击', action: '快速多选图片，用于批量操作' },
]

const faqs = [
  {
    q: '为什么侧边栏有数量，主区域却没有图？',
    a: '最常见原因不是图片丢了，而是顶部还残留着搜索词、日期范围、模型、LoRA、标签或收藏筛选。优先尝试“清空全部筛选”。',
  },
  {
    q: '为什么现在必须手动选择 output？',
    a: '从新版本开始，软件不再猜测 exe 的相对位置，而是要求绑定真实 output。这样可以适配任意安装位置，也能避免路径猜错导致的误扫描。',
  },
  {
    q: '提示词助手的数据来自哪里？',
    a: '运行时主要读取本地 `data/prompt-library/` 目录。它是可维护的数据目录，不依赖外部网络服务。',
  },
  {
    q: '为什么会看到性能模式提示？',
    a: '当当前目录图片很多时，软件会更偏向使用轻量分页、缩略图或预览图变体，目标是保持大图库下的可用性和流畅度。',
  },
  {
    q: '小窗里的确认框为什么不是系统弹窗？',
    a: '批量删除、清缓存等危险操作已经改为应用内确认弹层，避免出现 WebView 原生的 wails.localhost 确认框，也让视觉和主界面保持一致。',
  },
  {
    q: '小窗为什么要分页？',
    a: '小窗高度有限，分页可以减少一次渲染的图片数量，让滚动、选择和查看保持稳定。默认每页 120 张，也可以改成 60 或 240 张。',
  },
  {
    q: '切换了 output 之后要不要重装软件？',
    a: '不需要。重新绑定 output 即可，必要时再刷新图库或清理一次预览缓存。',
  },
  {
    q: '壁纸图片存在哪？删掉原图会怎样？',
    a: '选完壁纸后，软件会把图片复制一份到 `data/wallpaper/`，配置里只记录这个副本。所以原图移动、重命名或删除都不影响壁纸显示；想彻底清掉，在壁纸卡片里点“移除壁纸”即可。',
  },
  {
    q: '开了壁纸之后，为什么有些地方的背景没变？',
    a: '绝大多数面板都会跟随壁纸变半透明。唯一刻意保留纯色的是全屏大图查看器（灯箱）——它是“当前图片模糊放大垫底 + 压一层黑”的沉浸式环境，这样看图时不会被壁纸干扰。',
  },
]
</script>

<template>
  <div class="flex h-full flex-col bg-background text-foreground">
    <div class="flex h-16 shrink-0 items-center border-b bg-card/50 px-6 backdrop-blur-sm">
      <BookOpen class="mr-3 h-5 w-5 text-primary" />
      <div class="flex items-center gap-3">
        <h1 class="text-xl font-bold tracking-tight">使用文档</h1>
        <Badge variant="outline" class="rounded-full px-3 py-1 text-xs">v3.2.0</Badge>
      </div>
    </div>

    <div class="flex-1 overflow-y-auto">
      <div class="mx-auto flex max-w-5xl flex-col gap-10 px-6 py-6 pb-20">
        <section class="rounded-[28px] border border-border/70 bg-gradient-to-br from-primary/8 via-transparent to-transparent p-7">
          <div class="space-y-3">
            <div class="flex items-center gap-2 text-sm text-muted-foreground">
              <Clock3 class="h-4 w-4" />
              <span>当前版本说明</span>
            </div>
            <h2 class="text-3xl font-semibold tracking-tight">v3.2.0 使用说明</h2>
            <p class="max-w-3xl text-sm leading-7 text-muted-foreground">
              你现在看到的软件内文档已经与 GitHub README、发布说明、安装器版本和当前桌面程序同步到
              <span class="font-medium text-foreground">v3.2.0</span>。这一版是一次围绕“看图体验”的更新：重做了大图查看器（氛围光晕与倒影、缩略图导航条、
              批次概览层、参数 HUD、层叠连拍组），给图库加了 3D 环绕预览模式，并新增了背景壁纸自定义。
            </p>
          </div>
        </section>

        <section class="space-y-5">
          <div class="flex items-center gap-3">
            <div class="h-6 w-1 rounded-full bg-primary" />
            <h3 class="text-xl font-semibold">核心能力</h3>
          </div>
          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
            <Card
              v-for="feature in featureCards"
              :key="feature.title"
              class="border-border/70 bg-card/70 shadow-sm transition hover:border-primary/30 hover:bg-accent/20"
            >
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-3 text-base">
                  <span class="rounded-xl bg-primary/10 p-2 text-primary">
                    <component :is="feature.icon" class="h-4 w-4" />
                  </span>
                  <span>{{ feature.title }}</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription class="text-sm leading-7 text-muted-foreground">
                  {{ feature.description }}
                </CardDescription>
              </CardContent>
            </Card>
          </div>
        </section>

        <Separator />

        <section class="space-y-5">
          <div class="flex items-center gap-3">
            <div class="h-6 w-1 rounded-full bg-primary" />
            <h3 class="text-xl font-semibold">快速上手</h3>
          </div>
          <div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
            <Card
              v-for="guide in quickGuides"
              :key="guide.title"
              class="border-border/70 bg-card/70 shadow-sm"
            >
              <CardHeader class="pb-3">
                <CardTitle class="text-base">{{ guide.title }}</CardTitle>
              </CardHeader>
              <CardContent>
                <ol class="space-y-2 text-sm leading-7 text-muted-foreground">
                  <li v-for="(step, index) in guide.steps" :key="step">
                    {{ index + 1 }}. {{ step }}
                  </li>
                </ol>
              </CardContent>
            </Card>
          </div>
        </section>

        <Separator />

        <section class="space-y-5">
          <div class="flex items-center gap-3">
            <div class="h-6 w-1 rounded-full bg-primary" />
            <h3 class="flex items-center text-xl font-semibold">
              <Keyboard class="mr-2 h-5 w-5" />
              常用快捷操作
            </h3>
          </div>
          <div class="grid grid-cols-1 gap-3 md:grid-cols-2">
            <div
              v-for="shortcut in shortcuts"
              :key="shortcut.key"
              class="flex items-center justify-between rounded-2xl border border-border/70 bg-card/70 px-4 py-3"
            >
              <span class="pr-4 text-sm text-muted-foreground">{{ shortcut.action }}</span>
              <Badge variant="outline" class="rounded-full px-3 py-1 font-mono text-xs">
                {{ shortcut.key }}
              </Badge>
            </div>
          </div>
        </section>

        <Separator />

        <section class="space-y-5">
          <div class="flex items-center gap-3">
            <div class="h-6 w-1 rounded-full bg-primary" />
            <h3 class="flex items-center text-xl font-semibold">
              <HelpCircle class="mr-2 h-5 w-5" />
              常见问题
            </h3>
          </div>
          <div class="space-y-4">
            <Card
              v-for="item in faqs"
              :key="item.q"
              class="border-border/70 bg-card/70 shadow-sm"
            >
              <CardContent class="space-y-3 p-5">
                <div class="text-base font-semibold">{{ item.q }}</div>
                <p class="text-sm leading-7 text-muted-foreground">{{ item.a }}</p>
              </CardContent>
            </Card>
          </div>
        </section>

        <div class="pt-4 text-center text-sm text-muted-foreground">
          Comfy Manager v3.2.0 / 面向 ComfyUI 出图整理、大图查看、灵动小窗、目录治理与提示词工作流
        </div>
      </div>
    </div>
  </div>
</template>
