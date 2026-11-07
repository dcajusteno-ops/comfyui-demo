# Comfy Manager 发布说明

当前版本：`v3.3.0`

远程仓库：

- `origin`: `https://github.com/dcajusteno-ops/comfyui-demo.git`
- 默认分支：`master`

## 1. v3.3.0 发布重点

`v3.3.0` 是在 `v3.2.0` 基础上的体验与性能更新：移除图库 3D 环绕视图、重做全项目下拉菜单、完成大图库扫描与缓存优化，并补齐稳定性修复与后端单元测试。

### 本次核心变更

- **图库只保留网格视图**
  - 移除「网格 / 环绕」切换与 3D 环绕预览模式
  - 删除 `OrbitGallery.vue`、`.wp-scrim` 样式与相关入口
  - 缩略图尺寸、排序、连拍叠图、搜索与筛选保持不变
- **下拉菜单全面重做**
  - 新增 `AppSelect.vue`，替换项目中全部原生下拉框
  - 模型 / LoRA / 标签下拉支持输入关键字搜索
  - 固定宽度、深色 / 浅色主题自动适配，支持 ↑/↓、Enter、Esc 与空格操作
  - 选项显示数量、当前项勾选、长文本省略与悬停查看全称
- **大图库性能优化**
  - 元数据缓存改为单条读写；缓存写盘改为紧凑 JSON + 原子替换 + 防抖合并
  - 搜索 / 模型 / LoRA 筛选走轻量分页；图片总数计数缓存
  - 健康检查合并遍历；出图刷新跳过健康检查并并行化
  - 缩略图 4 并发限制；缩略图采样器改用双线性；目录监听重启防抖
  - 列表稳定排序，消除相同时间戳的翻页抖动
- **稳定性与修复**
  - 修复「我的词库」来源乱码，自定义提示词条目恢复识别与删除入口
  - 手动词条支持中文逗号 / 顿号 / 分号 / 问号拆分
  - 所有 JSON 配置改为原子写盘
  - 新增后端单元测试（排序 / 分页 / 目录范围 / 日期 / JSON 清洗 / 缓存写盘）
- 版本号在本版统一：`project.nsi` / `package.json` / README / docs / 软件内使用文档全部对齐 `v3.3.0`

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
git commit -m "release: v3.3.0"
git tag -a v3.3.0 -m "v3.3.0"
git push origin master
git push origin v3.3.0
```

### 2.6 推送 GitHub Release

```powershell
gh release create v3.3.0 .\desktop-app.exe .\ComfyManager-amd64-installer.exe `
  --title "v3.3.0" --notes-file <发布说明>
```

## 3. v3.3.0 验收清单

- 软件内使用文档版本显示为 `v3.3.0`
- GitHub README 版本显示为 `v3.3.0`
- 安装包属性中的版本为 `3.3.0`（来自 `project.nsi`）
- 根目录 `desktop-app.exe` 已覆盖最新构建
- 根目录 `ComfyManager-amd64-installer.exe` 已覆盖最新安装包
- `wails build -clean` 成功
- `wails build --nsis -clean` 成功
- 根目录发布产物与 `desktop-source/build/bin/` 哈希一致
- **从安装包中解出 `desktop-app.exe`，其哈希与根目录一致**（确认安装包内真的是新程序）
- Git tag `v3.3.0` 已推送，GitHub Release 已附带两个产物

## 4. 历史版本

| 版本 | 日期 | 说明 |
|---|---|---|
| v3.3.0 | 2026-11-07 | 移除 3D 环绕视图、统一可搜索下拉菜单、大图库扫描与缓存优化、编码与写盘修复、新增后端单元测试 |
| v3.2.0 | 2026-10-09 | 大图查看器增强（氛围画廊 / 缩略图导航条 / 批次概览层 / 参数 HUD / 层叠连拍组）、3D 环绕预览模式、背景壁纸自定义、图片与环绕视图缺陷修复 |
| v3.1.5 | 2026-05-08 | 第三方词库导入放开、随包附赠 NAI4 标签库、安装器同步 |
| v3.1.4 | 2026-04-27 | 提示词分类算法优化、分类下拉框卡顿修复 |
| v3.1.0 | 2026-05-08 | 灵动图库小窗、分页、应用内确认弹层、安装程序与文档同步更新 |
| v3.0.1 | 2026-04-27 | 日期归档黑屏修复、归档树折叠修复、根目录发布产物哈希校验、版本同步 |
| v3.0.0 | 2026-04-21 | 后端目录整理、文档全面重写、软件内文档升级、安装链恢复、发布收敛 |
| v2.2.1 | 2026-04-21 | 软件内使用文档与外部文档版本同步、重新打包桌面端和安装包 |
| v2.2.0 | 2026-04-21 | 大型图库性能模式、目录健康中心、预览变体、自定义目录增强 |
| v2.1.6 | 2026-04-19 | Windows 安装版同步、Prompt 增强、分页与缓存修复 |
