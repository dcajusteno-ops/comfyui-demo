# Comfy Manager 项目上下文

当前版本：`v3.3.1`  
更新时间：`2026-11-08`

## 1. 项目定位

Comfy Manager 是一个基于 **Wails v2 + Go + Vue 3** 的桌面应用，面向 ComfyUI 出图后的本地整理场景。

当前主线很明确：

- 浏览 ComfyUI output 图片
- 按日期、模型、LoRA、标签、收藏和笔记回看作品
- 管理目录、回收站、缓存和自定义目录
- 提供提示词助手、模板与自动规则，提升复用效率
- 提供右下角灵动图库小窗，支持置顶、分页、快速筛选和轻量批量操作
- 提供更完整的大图查看体验（氛围画廊、缩略图导航条、批次概览层、参数 HUD、层叠连拍组）
- 提供背景壁纸自定义

这不是云相册，也不是远程协作系统。  
它当前的产品重心仍然是：

`ComfyUI output 的本地整理器`

## 2. 技术栈

### 后端

- Go
- Wails v2
- fsnotify
- golang.org/x/image
- github.com/google/uuid

### 前端

- Vue 3
- Vite
- Tailwind CSS
- shadcn-vue
- lucide-vue-next
- vue-sonner

## 3. 关键目录

```text
comfy-manager/
├─ README.md
├─ docs/
├─ data/
├─ .trash/
├─ desktop-app.exe
├─ ComfyManager-amd64-installer.exe
└─ desktop-source/
   ├─ main.go
   ├─ backend/
   ├─ frontend/
   ├─ build/
   └─ wails.json
```

## 4. 当前后端结构

当前后端已经不再把主要逻辑堆在根目录 `app.go`。

新的结构是：

- `desktop-source/main.go`
  Wails 入口，只负责启动应用和绑定 `backend.App`
- `desktop-source/backend/`
  后端主要实现目录
- `desktop-source/backend/exports.go`
  为根入口暴露启动、关闭和图片服务包装

后端分组规则：

- `app.go`
  `App` 壳子与共享状态
- `app_core_*`
  生命周期、运行时、常量
- `app_feature_*`
  业务功能实现
- `app_support_*`
  内部辅助与基础设施
- `app_types_*`
  类型定义

详细映射见：[BACKEND_FILE_MAP.md](./BACKEND_FILE_MAP.md)

## 5. 当前前端结构

### 根级页面

- `frontend/src/App.vue`
  根级页面装配、视图切换、绑定刷新链路

### 核心页面 / 组件

- `frontend/src/components/Home.vue`
  工作台总览
- `frontend/src/components/ImageGallery.vue`
  图库主视图，含网格视图与缩略图尺寸调节
- `frontend/src/components/DateWorkbench.vue`
  日期产出工作台
- `frontend/src/components/StatisticsDashboard.vue`
  数据视界
- `frontend/src/components/Documentation.vue`
  软件内使用文档
- `frontend/src/components/CompactWindow.vue`
  灵动图库小窗，承载置顶窗口、目录选择、搜索、分页、批量选择和图片详情入口
- `frontend/src/components/ProfileCenter.vue`
  个人中心，含外观与偏好、背景壁纸入口
- `frontend/src/components/PromptAssistantPage.vue`
  提示词助手
- `frontend/src/components/AutoRulesPanel.vue`
  自动规则引擎

### 大图查看器（灯箱）相关

- `Lightbox.vue`
  灯箱状态中心（缩放、切换、批次概览层开关、键盘）
- `LightboxViewer.vue`
  画布层（承载氛围背景、导航条、HUD、层叠卡牌）
- `LightboxToolbar.vue` / `ImageMetadataPanel.vue`
  右侧工具栏与元数据面板
- `AmbientBackdrop.vue` / `GroundReflection.vue`
  主色氛围光晕与地面倒影
- `ImageFilmstrip.vue`
  底部缩略图导航条
- `BatchOverview.vue`
  批次概览层（按同日 / 同模型 / 同提示词 / 连拍组分组）
- `ImageHud.vue`
  参数 HUD
- `StackFan.vue`
  连拍组 3D 层叠卡牌
- `AppSelect.vue`
  统一的自绘下拉（固定宽度、可搜索、键盘导航）
- `WallpaperPanel.vue`
  背景壁纸编辑器

### 当前已拆分的 composables

- `useImages.js`
  主图库状态入口
- `useGalleryData.js`
  图库数据与分页请求
- `useGalleryHelpers.js`
  辅助格式化与筛选工具
- `useLibraryMeta.js`
  标签、收藏、笔记等资料元数据
- `useWorkbenchFilters.js`
  工作台日期 / 模型 / LoRA 筛选
- `useImageZoom.js`
  图片缩放 / 平移（灯箱与灵动小窗共用）
- `useDominantColor.js`
  图片主色提取（氛围光晕用）
- `useParallaxFloat.js`
  鼠标视差浮动
- `useBatchGroups.js`
  批次分组派生（同日 / 同模型 / 同提示词 / 连拍组）

### 根级状态模块

- `frontend/src/theme.js`
  亮 / 暗主题与 View Transition 切换
- `frontend/src/wallpaper.js`
  背景壁纸配置、CSS 变量注入与持久化

## 6. 数据持久化

常见持久化文件位于 `data/`：

- `favorites.json`
- `tags.json`
- `image-tags.json`
- `image-notes.json`
- `custom-roots.json`
- `settings.json`（用户资料、性能设置、**壁纸配置**）
- `auto-rules.json`
- `trash-metadata.json`
- `image-meta-cache.json`
- `prompt-library/`
- `profile/`（头像图片）
- `wallpaper/`（壁纸图片，settings 里只存 `__wallpaper__/` 相对路径）

## 7. 关键业务链路

### 图片刷新链

`fsnotify` -> `images:changed` 事件 -> 前端订阅 -> 图库 / 工作台 / 统计刷新

### output 绑定链

用户选择 output -> 后端校验目录 -> 保存 settings -> 重启 watcher -> 前端刷新当前视图

### 大图库性能链

目录扫描 -> 轻量分页 / 预览变体 -> 延迟元数据读取 -> Lightbox 查看原图与细节

### 灵动图库小窗链

侧边栏打开小窗 -> Wails 调整窗口大小和置顶状态 -> 前端小窗读取最新图库数据 -> 目录 / 搜索 / 筛选 -> 分页显示 -> 图片详情或批量操作

### 提示词复用链

本地词库 -> 搜索 / 分类 / 模板 -> 拼装 Prompt -> 回写使用上下文

### 背景壁纸链

个人中心选图 -> 原生文件对话框 -> 复制进 `data/wallpaper/wallpaper.<ext>` -> settings 存 `__wallpaper__/` 相对路径 -> 前端注入 `html.has-wallpaper` 类与 CSS 变量 -> 表面令牌半透明化让壁纸透出

## 8. v3.3.x 当前状态总结

`v3.3.1` 在 `v3.3.0` 基础上完成了一轮 UI 可访问性与视觉一致性打磨：

- 全项目图标按钮、开关、滑杆、输入框补齐可访问名称与状态语义
- 统计图表 / 热力图 / 图片卡片补充 ARIA 语义，图片卡片改为真实按钮覆盖层
- 控件边界 / 焦点环 / 次要文字对比度统一提升，深色层级重新收敛
- 修复壁纸模式下提示气泡与首页大图文字遮罩失效
- 圆角与字号收敛：卡片 `rounded-2xl`、控件 `rounded-xl`、最小字号 11px
- Toast 与确认框跟随主题

`v3.3.0` 的主要成果：

- 图库移除 3D 环绕视图，只保留网格视图；同步清理组件、样式与文档
- 新增 `AppSelect.vue` 自绘下拉，替换全部原生下拉框；模型 / LoRA / 标签支持搜索，固定宽度并适配主题
- 元数据缓存改为单条读写；缓存写盘改为紧凑 JSON + 原子替换 + 防抖合并
- 轻量分页支持搜索 / 模型 / LoRA；图片总数计数缓存；列表稳定排序
- 健康检查合并遍历；出图事件跳过健康检查并并行刷新
- 缩略图 4 并发限制与采样器分级；目录监听重启防抖
- 修复「我的词库」乱码与手动词条中文标点拆分
- 所有 JSON 配置原子写盘；新增后端单元测试
- 版本号统一：`project.nsi` / `package.json` / README / docs / 软件内使用文档全部对齐 `v3.3.1`

## 9. 当前边界与原则

- 当前阶段不继续做更深的 Go 包拆分
- 以后新增后端能力，优先继续放进 `backend/` 现有分组
- 优先维持外部行为稳定，再做结构优化
- 新增控制类交互时优先复用语义令牌（`bg-card` / `text-muted-foreground` 等），不要在组件里写死颜色——写死的深色底会把背景壁纸盖死
- 发布时必须同时同步：
  - 根目录 exe
  - 安装程序
  - 根 README
  - `docs/README.md`
  - `docs/PROJECT_CONTEXT.md`
  - `docs/RELEASE.md`
  - 软件内 `Documentation.vue`

## 10. 下一步建议

若继续推进，优先级建议如下：

1. 清理剩余编码历史包袱，统一前端中文文案
2. 继续收敛前端大组件的职责边界
3. 检查小窗和主窗口图库筛选逻辑是否可以进一步复用
4. 评估是否需要补一份图库 / 工作台 / 小窗数据流图
5. 评估灯箱是否也要透出背景壁纸（当前灯箱是「模糊当前图 + 纯黑」的沉浸式设计，刻意未跟随壁纸）
