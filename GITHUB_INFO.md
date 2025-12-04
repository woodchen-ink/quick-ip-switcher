# GitHub 项目信息 / GitHub Project Information

## 项目名称推荐 / Recommended Project Names

1. **quick-ip-switcher** ⭐ (推荐/Recommended)
2. network-config-manager
3. ip-switcher-win
4. windows-ip-switcher

## 项目简介 / Project Description

### 中文
快速IP切换工具 - 一个轻量级的Windows GUI工具,用于管理和一键切换网络适配器IP配置。支持多网卡、多配置管理,使用Windows原生API,无窗口弹出,配置持久化存储。

### English
Quick IP Switcher - A lightweight Windows GUI tool for managing and switching network adapter IP configurations with one click. Supports multiple adapters and configurations, uses Windows native API, no popup windows, with persistent configuration storage.

## GitHub Topics / 主题标签

```
windows
go
golang
network-configuration
ip-switcher
network-manager
windows-gui
system-utility
network-tools
winapi
lxn-walk
network-adapter
ip-configuration
```

## GitHub About 设置

**Description (描述):**
```
🔧 A lightweight Windows GUI tool for managing and switching network IP configurations | Windows 网络IP配置一键切换工具
```

**Website (网站):**
```
留空或填写你的GitHub Pages链接
```

**Topics (主题):**
```
windows, go, network-configuration, ip-switcher, network-manager, windows-gui
```

## 发布第一个版本的步骤

1. **初始化Git仓库**
```bash
git init
git add .
git commit -m "Initial commit: Quick IP Switcher v2.1"
```

2. **创建GitHub仓库**
- 访问 https://github.com/new
- 仓库名: `quick-ip-switcher`
- 描述: 填写上面的Description
- 选择: Public
- 不勾选任何初始化选项(README, .gitignore, license)

3. **推送到GitHub**
```bash
git remote add origin https://github.com/YOUR_USERNAME/quick-ip-switcher.git
git branch -M main
git push -u origin main
```

4. **创建第一个Release**
```bash
# 创建tag
git tag -a v2.1.0 -m "Release v2.1.0 - Initial public release"
git push origin v2.1.0
```

GitHub Actions 会自动构建并创建 Release!

## 后续版本发布流程

每次想发布新版本时:

```bash
# 1. 修改代码并提交
git add .
git commit -m "Fix: some bug description"

# 2. 创建新tag
git tag -a v2.1.1 -m "Release v2.1.1 - Bug fixes"

# 3. 推送tag
git push origin v2.1.1
```

GitHub Actions 会自动:
- 编译 Windows 可执行文件
- 生成资源文件 (manifest)
- 创建 ZIP 压缩包
- 发布 GitHub Release
- 附加编译好的文件

## README Badge 替换

在 README.md 中,将 `YOUR_USERNAME` 替换为你的 GitHub 用户名:

```markdown
[![Release](https://img.shields.io/github/v/release/YOUR_USERNAME/quick-ip-switcher)](https://github.com/YOUR_USERNAME/quick-ip-switcher/releases)
```

替换为:

```markdown
[![Release](https://img.shields.io/github/v/release/你的用户名/quick-ip-switcher)](https://github.com/你的用户名/quick-ip-switcher/releases)
```

## 项目亮点 / Project Highlights

✨ **零依赖单文件** - Single executable, no dependencies
✨ **Windows原生API** - Uses Windows native API, no PowerShell popups
✨ **配置持久化** - Persistent JSON configuration storage
✨ **多网卡支持** - Support multiple network adapters
✨ **一键切换** - One-click IP configuration switching
✨ **自动构建** - Automated builds with GitHub Actions
✨ **开源免费** - Open source and free (MIT License)

## 可选: 添加截图

在 README 中可以添加程序截图:

1. 运行程序,截取主界面
2. 保存为 `screenshot.png`
3. 在 GitHub 仓库中创建 `docs/images/` 目录
4. 上传截图
5. 在 README.md 中添加:

```markdown
## 📸 Screenshots / 界面预览

![Main Interface](docs/images/screenshot.png)
```

## Star History (可选)

等项目有一定star后,可以添加:

```markdown
## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=YOUR_USERNAME/quick-ip-switcher&type=Date)](https://star-history.com/#YOUR_USERNAME/quick-ip-switcher&Date)
```
