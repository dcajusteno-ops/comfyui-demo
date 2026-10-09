# Windows 安装器说明

当前版本：`v3.2.0`

更新时间：`2026-10-09`

## 1. 安装器目标

本项目提供标准 Windows 安装流程，目标是：

- 用户能选择安装目录
- 程序文件和运行数据保持在同一安装目录附近
- 安装完成后可直接从桌面或开始菜单启动

## 2. 默认安装布局

安装完成后，典型目录结构如下：

```text
H:\Comfy Manager\
├─ desktop-app.exe
├─ data\
└─ .trash\
```

说明：

- 程序主文件：`desktop-app.exe`
- 运行数据：`data\`
- 回收站：`.trash\`

## 3. 如何构建安装器

在 `desktop-source/` 目录执行：

```powershell
wails build -clean -nsis
```

前提：

- 系统已安装 NSIS（本机：`C:\Program Files (x86)\NSIS\makensis.exe`）
- `makensis` 可在命令行中调用；本机它**不在 PATH 里**，需要先把 `C:\Program Files (x86)\NSIS` 加进 PATH，否则 NSIS 步骤会失败

构建成功后，输出位于：

```text
desktop-source\build\bin\ComfyManager-amd64-installer.exe
```

## 4. 发布时如何覆盖根目录安装器

```powershell
Copy-Item .\desktop-source\build\bin\ComfyManager-amd64-installer.exe .\ComfyManager-amd64-installer.exe -Force
```

## 5. 安装器行为

当前安装器会：

- 显示欢迎页
- 允许用户选择安装目录（默认 `$PROFILE\Comfy Manager`）
- 自动安装或检测 WebView2 运行时
- 复制程序文件
- 复制 `data\prompt-library\` 到安装目录
- 创建开始菜单快捷方式
- 创建桌面快捷方式

安装器带 `requireAdministrator` 清单，运行时会请求管理员权限。

## 6. 当前版本注意事项

- `v3.2.0` 的安装包与根目录 `desktop-app.exe` 使用同一次构建结果，安装包内已包含大图查看器增强、3D 环绕预览模式与背景壁纸功能
- 安装包显示的版本来自 `build/windows/installer/project.nsi` 的 `INFO_PRODUCTVERSION` —— **改版本号必须改这里**，否则安装包属性里还是旧版本
- 发布前必须确保根目录安装包是最新构建结果
- 如果安装器构建失败，优先检查 `makensis` 是否可用

## 7. 本次发布校验点

- `wails build -clean -nsis` 成功生成 `desktop-source\build\bin\ComfyManager-amd64-installer.exe`
- 根目录 `ComfyManager-amd64-installer.exe` 已由 `build\bin` 新产物覆盖
- 根目录 `desktop-app.exe` 已由 `build\bin` 新产物覆盖
- **把 `desktop-app.exe` 从安装包里解出来比对 sha256**（NSIS 默认压缩，直接搜字符串搜不到；可用 Bandizip CLI：`bz x -o:<dir> <installer.exe>`），确认与根目录一致
- 软件内使用文档与 GitHub 文档同步描述大图查看器、3D 环绕与背景壁纸
