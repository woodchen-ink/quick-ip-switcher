# Quick IP Switcher - 快速IP切换工具

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![Platform](https://img.shields.io/badge/platform-Windows-0078D6?logo=windows)](https://www.microsoft.com/windows)
[![Release](https://img.shields.io/github/v/release/YOUR_USERNAME/quick-ip-switcher)](https://github.com/YOUR_USERNAME/quick-ip-switcher/releases)

一个功能完整的 Windows GUI 工具,用于管理和切换网络适配器的IP配置。

使用 Go 语言开发,单个 exe 文件,无需安装任何依赖。

[English Documentation](README_EN.md)

## ✨ 功能特性

✅ **网卡选择** - 自动列出所有网络适配器,显示友好名称(如"以太网3")
✅ **配置管理** - 支持添加、编辑、删除多个网络配置
✅ **一键切换** - 快速在不同配置间切换
✅ **配置持久化** - 配置自动保存到 JSON 文件
✅ **实时状态** - 显示当前网络配置信息
✅ **原生GUI** - Windows 原生界面,简洁易用
✅ **单文件** - 无依赖,直接运行
✅ **无窗口弹出** - 使用Windows API,不弹出PowerShell窗口

## 📦 下载安装

### 方式一: 下载发布版本(推荐)

从 [Releases](https://github.com/YOUR_USERNAME/quick-ip-switcher/releases) 页面下载最新版本:
- `quick-ip-switcher.exe` - 单文件可执行程序(推荐)
- `quick-ip-switcher-vX.X.X-windows-amd64.zip` - 包含文档的完整包

### 方式二: 从源码编译

```bash
git clone https://github.com/YOUR_USERNAME/quick-ip-switcher.git
cd quick-ip-switcher
build.bat
```

### 快速开始

1. **双击 `quick-ip-switcher.exe`**
2. 在 UAC 提示中点击"是"授予管理员权限
3. 在下拉框中选择要管理的网络适配器
4. 点击"新建配置"添加IP配置
5. 选择配置后点击"切换"

### 详细步骤

#### 1. 选择网络适配器
- 程序启动后会自动列出所有网卡
- 在"网络适配器"下拉框中选择要管理的网卡(显示友好名称,如"以太网3")
- 程序会自动保存你的选择

#### 2. 添加配置
- 点击"新建配置"按钮
- 填写以下信息:
  - **配置名称**: 例如"公司网络"、"家庭网络"
  - **IP地址**: 例如 10.0.2.2
  - **子网掩码**: 例如 255.0.0.0
  - **默认网关**: 例如 10.0.0.1
  - **DNS服务器**: 例如 10.0.0.1
- 点击"保存"

#### 3. 切换配置
- 在配置列表中选择要切换的配置
- 点击"切换"按钮或双击配置项
- 等待切换完成提示

#### 4. 管理配置
- **编辑**: 选中配置后点击"编辑"
- **删除**: 选中配置后点击"删除"
- **查看状态**: 点击"刷新状态"查看当前网络信息

## 💡 使用场景

### 示例1: 公司和家庭网络切换
```
配置1: 公司网络
  IP: 10.0.2.100
  网关: 10.0.0.1
  DNS: 10.0.0.1

配置2: 家庭网络
  IP: 192.168.1.100
  网关: 192.168.1.1
  DNS: 192.168.1.1
```

### 示例2: 测试不同网关
```
配置1: 网关1
  IP: 10.0.2.2
  网关: 10.0.0.1

配置2: 网关2
  IP: 10.0.2.2
  网关: 10.0.0.5
```

## 📱 界面说明

```
┌──────────────────────────────────────────────────┐
│ 网络配置管理工具                                 │
│                                                  │
│ 网络适配器: [以太网3 (Up)        ▼] [刷新]     │
│                                                  │
│ 当前配置: 配置1 (网关 10.0.0.1)                  │
│                                                  │
│ ┌────────────────────────────────────────────┐  │
│ │ 网络适配器: 以太网3                        │  │
│ │                                            │  │
│ │ IP 地址    : 10.0.2.2                     │  │
│ │ 默认网关   : 10.0.0.1                     │  │
│ │ DNS 服务器 : 10.0.0.1                     │  │
│ └────────────────────────────────────────────┘  │
│                                                  │
│ 保存的配置:                                      │
│                                                  │
│ ┌───────────────────────────┐  ┌──────────┐    │
│ │ 配置1 (IP: 10.0.2.2...)   │  │  切换    │    │
│ │ 配置2 (IP: 10.0.2.2...)   │  │          │    │
│ │                           │  │ 新建配置  │    │
│ │                           │  │          │    │
│ │                           │  │  编辑    │    │
│ │                           │  │          │    │
│ │                           │  │  删除    │    │
│ └───────────────────────────┘  └──────────┘    │
│                                                  │
│                               [刷新状态]         │
└──────────────────────────────────────────────────┘
```

## ⚙️ 配置文件

配置自动保存在程序所在目录的 `config.json` 文件中:

```json
{
  "selected_interface": "以太网3",
  "network_configs": [
    {
      "name": "配置1 (网关 10.0.0.1)",
      "ip": "10.0.2.2",
      "subnet_mask": "255.0.0.0",
      "gateway": "10.0.0.1",
      "dns": "10.0.0.1"
    },
    {
      "name": "配置2 (网关 10.0.0.5)",
      "ip": "10.0.2.2",
      "subnet_mask": "255.0.0.0",
      "gateway": "10.0.0.5",
      "dns": "10.0.0.5"
    }
  ]
}
```

## 🛠️ 开发说明

### 环境要求
- Go 1.16 或更高版本
- Windows 10/11

### 项目结构
```
switch_ips/
├── main.go              # 主程序和GUI
├── config.go            # 配置管理
├── network.go           # 网络操作接口
├── network_windows.go   # Windows API实现
├── rsrc.syso            # Windows资源文件(manifest)
├── go.mod               # Go 模块
├── build.bat            # 编译脚本
├── quick-ip-switcher.exe # 可执行文件
├── config.json          # 配置文件(运行时生成)
├── README.md            # 说明文档
└── 使用说明.txt         # 快速参考
```

### 编译

```bash
# 方法1: 使用编译脚本(推荐)
build.bat

# 方法2: 手动编译
go build -ldflags="-H windowsgui" -o quick-ip-switcher.exe
```

### 依赖库
- `github.com/lxn/walk` - Windows GUI 库
- `golang.org/x/sys/windows` - Windows API
- `golang.org/x/sys/windows/registry` - 注册表操作

## 🔧 系统要求

- **操作系统**: Windows 10/11
- **权限**: 管理员权限(修改网络配置需要)
- **运行时**: 无需安装任何运行时

## ❓ 常见问题

### Q: 为什么需要管理员权限?
A: 修改网络配置需要管理员权限,这是 Windows 系统的安全限制。程序会自动请求UAC提升权限。

### Q: 配置保存在哪里?
A: 配置保存在程序所在目录的 `config.json` 文件中。

### Q: 如何查看我的网络适配器名称?
A: 程序会自动列出所有网络适配器,显示友好名称(如"以太网3"、"WLAN"),你只需要在下拉框中选择即可。

### Q: 切换失败怎么办?
A: 请检查:
1. 是否以管理员权限运行
2. 网络适配器是否已启用
3. IP配置是否正确(IP、子网掩码、网关格式)
4. 查看错误提示信息

### Q: 可以管理多个网卡吗?
A: 可以!只需在下拉框中切换不同的网卡,每个网卡可以有自己的配置列表。

### Q: 支持 IPv6 吗?
A: 当前版本仅支持 IPv4。

### Q: 配置会丢失吗?
A: 不会,所有配置自动保存到 JSON 文件,程序重启后自动加载。

### Q: 为什么不会弹出PowerShell窗口?
A: 程序使用Windows原生API (GetAdaptersInfo) 和隐藏窗口的netsh命令,不会有任何黑窗口闪现。

## 🔍 技术实现

### 核心技术
- **语言**: Go 1.25
- **GUI框架**: lxn/walk (Windows 原生控件)
- **配置存储**: JSON 格式
- **网络操作**: Windows API + netsh命令

### Windows API
- `GetAdaptersInfo` - 获取网络适配器信息(IP、网关等)
- `Registry API` - 从注册表读取友好名称
- `AllocateAndInitializeSid` - 检查管理员权限
- `netsh` - 设置IP配置(隐藏窗口执行)

### 关键代码
```go
// 1. 获取友好名称
func GetFriendlyName(adapterGUID string) string {
    keyPath := `SYSTEM\CurrentControlSet\Control\Network\...`
    key, _ := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.QUERY_VALUE)
    name, _, _ := key.GetStringValue("Name")
    return name
}

// 2. 隐藏窗口执行命令
cmd.SysProcAttr = &syscall.SysProcAttr{
    HideWindow:    true,
    CreationFlags: CREATE_NO_WINDOW,
}
```

## 🎯 功能亮点

### vs PowerShell 脚本
- ✅ 无编码问题,完美支持中文
- ✅ 图形化界面,操作更直观
- ✅ 配置持久化,无需手动编辑代码
- ✅ 支持多网卡多配置管理
- ✅ 无窗口弹出,用户体验更好

### vs 其他工具
- ✅ 单文件,无需安装
- ✅ 开源免费,代码透明
- ✅ 轻量级,体积小(<10MB)
- ✅ 启动快速,操作流畅
- ✅ 使用Windows原生API,性能最优

## 🔐 安全说明

- ✅ 本工具仅修改本地网络配置
- ✅ 不涉及任何网络传输
- ✅ 所有代码开源,可自行审查
- ✅ 配置文件存储在本地
- ✅ 需要管理员权限是系统安全要求
- ✅ 使用Windows官方API,安全可靠

## 📝 更新日志

### v2.1 (2024-12-04)
- ✨ 使用Windows API替代PowerShell,无窗口弹出
- ✨ 从注册表读取网卡友好名称
- ✨ 修复配置列表实时更新问题
- ✨ 修复启动时自动检测配置
- ✨ 优化详情显示格式(正确换行)
- 🐛 修复netsh命令执行失败问题

### v2.0 (2024-12-04)
- ✨ 全新设计,支持网卡选择
- ✨ 支持添加、编辑、删除配置
- ✨ 配置持久化存储
- ✨ 支持管理多个网卡
- ✨ 优化界面和交互

### v1.0 (2024-12-04)
- 🎉 首次发布
- ✅ 基本的配置切换功能

## 📄 许可证

MIT License - 可自由使用、修改和分发

## 👨‍💻 贡献

欢迎提交 Issue 和 Pull Request!

---

**开发者**: Claude Code
**日期**: 2024-12-04
**版本**: 2.1
**项目地址**: 本地工具
