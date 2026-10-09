# Comfy Manager 发布说明

当前版本：`v3.2.0`

远程仓库：

- `origin`: `https://github.com/dcajusteno-ops/comfyui-demo.git`
- 默认分支：`master`

## 1. v3.2.0 发布重点

`v3.2.0` 是在 `v3.1.5` 基础上的功能性更新，重点全部落在“看图体验”上：重做大图查看器、给图库加 3D 环绕预览模式、新增背景壁纸自定义，并修掉了几个影响可用性的缺陷。

### 本次核心变更

- **大图查看器（灯箱）增强**
  - 氛围画廊：按当前图片主色渲染两侧光晕、主图下方地面倒影、切换图片的“相纸显影”动效、鼠标轻微视差
  - 底部新增缩略图导航条（自动居中、点击跳转），并提供批次概览入口按钮
  - 批次概览层：`空格` 呼出半透明缩略图矩阵，支持同日 / 同模型 / 同提示词 / 连拍组分组，点击即跳
  - 参数 HUD：模型 / 采样器 / 调度器 / seed / steps / CFG，放大时显示当前视野在原图的位置框，可开关
  - 连拍组由“1 / 3 角标”改为 3D 层叠卡牌（悬停散开、点击抽卡）
  - 灵动小窗的图片详情补齐缩放与拖动，与主窗口体验对齐
- **图库新增 3D 环绕预览模式**
  - 工具栏“网格 / 环绕”切换，选择持久化
  - 图片沿中央光柱多层环绕，拖动旋转带惯性、滚轮上下滚动、空闲自动缓慢旋转
  - “每圈张数”可调（4–40），圈数与整体缩放自适应窗口，缩放窗口不会裁切卡片
  - 点击任意图片直接进入大图查看
- **新增背景壁纸自定义**（个人中心 → 外观与偏好 → 背景壁纸）
  - 本地选图，图片复制到 `data/wallpaper`，原图移动/删除不影响显示
  - 可调：启用开关、填充方式（铺满 / 完整 / 原始）、压暗、面板不透明度、面板毛玻璃
  - 侧栏、卡片、对话框半透明透出壁纸；3D 环绕等沉浸式深色场景同步适配
- **缺陷修复**
  - 开发模式下图片全部显示为破图（vite 的 SPA 回退把图片请求返回成 HTML）
  - 3D 环绕视图滚轮方向相反、卡片被放得过大而裁切、背景不跟随壁纸
- 版本号在本版统一：`project.nsi` / `package.json` / README / docs / 软件内使用文档全部对齐 `v3.2.0`

## 2. 标准发布流程

### 2.1 同步文档

正式发布前至少同步这些文件：

- 根目录 `README.md`
- `docs/README.md`
- `docs/PROJECT_CONTEXT.md`
- `docs/BACKEND_FILE_MAP.md`
- `docs/REFACTOR_PLAN.md`
- `docs/RELEASE.md`
- `docs/WINDOWS_INSTALLER.md`
- 软件内 `frontend/src/components/Documentation.vue`
- 版本号来源：`desktop-source/build/windows/installer/project.nsi` 的 `INFO_PRODUCTVERSION`（**决定安装包显示的版本**）与 `frontend/package.json`

### 2.2 构建桌面程序

在 `desktop-source` 下执行：

```powershell
wails build -clean
```

### 2.3 构建安装程序

确保系统可用 `makensis`（本机路径：`C:\Program Files (x86)\NSIS\makensis.exe`，不在 PATH 里时需要先加进去）后执行：

```powershell
wails build -clean -nsis
```

### 2.4 覆盖根目录发布产物

```powershell
Copy-Item .\desktop-source\build\bin\desktop-app.exe .\desktop-app.exe -Force
Copy-Item .\desktop-source\build\bin\ComfyManager-amd64-installer.exe .\ComfyManager-amd64-installer.exe -Force
```

### 2.5 Git 发布步骤

```powershell
git add -A
git commit -m "release: v3.2.0"
git tag -a v3.2.0 -m "v3.2.0"
git push origin master
git push origin v3.2.0
```

### 2.6 推送 GitHub Release

```powershell
gh release create v3.2.0 .\desktop-app.exe .\ComfyManager-amd64-installer.exe `
  --title "v3.2.0" --notes-file <发布说明>
```

## 3. v3.2.0 验收清单

- 软件内使用文档版本显示为 `v3.2.0`
- GitHub README 版本显示为 `v3.2.0`
- 安装包属性中的版本为 `3.2.0`（来自 `project.nsi`）
- 根目录 `desktop-app.exe` 已覆盖最新构建
- 根目录 `ComfyManager-amd64-installer.exe` 已覆盖最新安装包
- `wails build -clean` 成功
- `wails build -clean -nsis` 成功
- 根目录发布产物与 `desktop-source/build/bin/` 哈希一致
- **从安装包中解出 `desktop-app.exe`，其哈希与根目录一致**（确认安装包内真的是新程序）
- Git tag `v3.2.0` 已推送，GitHub Release 已附带两个产物

## 4. 历史版本

| 版本 | 日期 | 说明 |
|---|---|---|
| v3.2.0 | 2026-10-09 | 大图查看器增强（氛围画廊 / 缩略图导航条 / 批次概览层 / 参数 HUD / 层叠连拍组）、3D 环绕预览模式、背景壁纸自定义、图片与环绕视图缺陷修复 |
| v3.1.5 | 2026-05-08 | 第三方词库导入放开、随包附赠 NAI4 标签库、安装器同步 |
| v3.1.4 | 2026-04-27 | 提示词分类算法优化、分类下拉框卡顿修复 |
| v3.1.0 | 2026-05-08 | 灵动图库小窗、分页、应用内确认弹层、安装程序与文档同步更新 |
| v3.0.1 | 2026-04-27 | 日期归档黑屏修复、归档树折叠修复、根目录发布产物哈希校验、版本同步 |
| v3.0.0 | 2026-04-21 | 后端目录整理、文档全面重写、软件内文档升级、安装链恢复、发布收敛 |
| v2.2.1 | 2026-04-21 | 软件内使用文档与外部文档版本同步、重新打包桌面端和安装包 |
| v2.2.0 | 2026-04-21 | 大型图库性能模式、目录健康中心、预览变体、自定义目录增强 |
| v2.1.6 | 2026-04-19 | Windows 安装版同步、Prompt 增强、分页与缓存修复 |
