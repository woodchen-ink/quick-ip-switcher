# Quick IP Switcher

A lightweight Windows GUI tool for managing and switching network adapter IP configurations with a single click.

Built with Go, packaged as a single executable with no dependencies required.

[中文文档](README.md)

## ✨ Features

✅ **Adapter Selection** - Automatically lists all network adapters with friendly names (e.g., "Ethernet 3")
✅ **Configuration Management** - Add, edit, and delete multiple network configurations
✅ **One-Click Switching** - Quickly switch between different configurations
✅ **Persistent Storage** - Configurations automatically saved to JSON file
✅ **Real-Time Status** - Display current network configuration information
✅ **Native GUI** - Clean Windows native interface
✅ **Single Executable** - No dependencies, run directly
✅ **No Window Popups** - Uses Windows API, no PowerShell windows flash

## 📦 Quick Start

1. **Download** `quick-ip-switcher.exe` from [Releases](https://github.com/YOUR_USERNAME/quick-ip-switcher/releases)
2. **Double-click** to run (UAC prompt will appear for admin rights)
3. **Select** your network adapter from the dropdown
4. **Click** "New Config" to add IP configurations
5. **Select** a configuration and click "Switch" to apply

## 💡 Use Cases

### Example 1: Office and Home Network Switching
```
Config 1: Office Network
  IP: 10.0.2.100
  Gateway: 10.0.0.1
  DNS: 10.0.0.1

Config 2: Home Network
  IP: 192.168.1.100
  Gateway: 192.168.1.1
  DNS: 192.168.1.1
```

### Example 2: Testing Different Gateways
```
Config 1: Gateway 1
  IP: 10.0.2.2
  Gateway: 10.0.0.1

Config 2: Gateway 2
  IP: 10.0.2.2
  Gateway: 10.0.0.5
```

## 📱 Interface

```
┌──────────────────────────────────────────────────┐
│ Network Configuration Manager                    │
│                                                  │
│ Network Adapter: [Ethernet 3 (Up)    ▼] [Refresh]│
│                                                  │
│ Current Config: Config1 (Gateway 10.0.0.1)       │
│                                                  │
│ ┌────────────────────────────────────────────┐  │
│ │ Network Adapter: Ethernet 3                │  │
│ │                                            │  │
│ │ IP Address    : 10.0.2.2                  │  │
│ │ Default Gateway: 10.0.0.1                 │  │
│ │ DNS Server     : 10.0.0.1                 │  │
│ └────────────────────────────────────────────┘  │
│                                                  │
│ Saved Configurations:                            │
│                                                  │
│ ┌───────────────────────────┐  ┌──────────┐    │
│ │ Config1 (IP: 10.0.2.2...) │  │  Switch  │    │
│ │ Config2 (IP: 10.0.2.2...) │  │          │    │
│ │                           │  │New Config│    │
│ │                           │  │          │    │
│ │                           │  │   Edit   │    │
│ │                           │  │          │    │
│ │                           │  │  Delete  │    │
│ └───────────────────────────┘  └──────────┘    │
│                                                  │
│                               [Refresh Status]   │
└──────────────────────────────────────────────────┘
```

## ⚙️ Configuration File

Configurations are automatically saved in `config.json` in the program directory:

```json
{
  "selected_interface": "Ethernet 3",
  "network_configs": [
    {
      "name": "Config 1 (Gateway 10.0.0.1)",
      "ip": "10.0.2.2",
      "subnet_mask": "255.0.0.0",
      "gateway": "10.0.0.1",
      "dns": "10.0.0.1"
    }
  ]
}
```

## 🛠️ Development

### Requirements
- Go 1.16 or higher
- Windows 10/11

### Project Structure
```
switch_ips/
├── main.go              # Main program and GUI
├── config.go            # Configuration management
├── network.go           # Network operations interface
├── network_windows.go   # Windows API implementation
├── rsrc.syso            # Windows resource file (manifest)
├── go.mod               # Go modules
├── build.bat            # Build script
└── config.json          # Configuration file (generated at runtime)
```

### Build

```bash
# Method 1: Using build script (recommended)
build.bat

# Method 2: Manual build
go build -ldflags="-H windowsgui" -o quick-ip-switcher.exe
```

### Dependencies
- `github.com/lxn/walk` - Windows GUI library
- `golang.org/x/sys/windows` - Windows API
- `golang.org/x/sys/windows/registry` - Registry operations

## 🔧 System Requirements

- **Operating System**: Windows 10/11
- **Permissions**: Administrator rights (required for network configuration changes)
- **Runtime**: No runtime installation required

## ❓ FAQ

### Q: Why does it need administrator rights?
A: Modifying network configuration requires administrator privileges - this is a Windows security requirement. The program will automatically request UAC elevation.

### Q: Where are configurations saved?
A: Configurations are saved in `config.json` in the program directory.

### Q: Can I manage multiple network adapters?
A: Yes! Just select a different adapter from the dropdown. Each adapter can have its own list of configurations.

### Q: Does it support IPv6?
A: The current version only supports IPv4.

### Q: Why don't I see PowerShell windows?
A: The program uses Windows native API (GetAdaptersInfo) and hidden-window netsh commands - no black windows will appear.

## 🔍 Technical Details

### Core Technologies
- **Language**: Go 1.25
- **GUI Framework**: lxn/walk (Windows native controls)
- **Configuration Storage**: JSON format
- **Network Operations**: Windows API + netsh commands

### Windows APIs
- `GetAdaptersInfo` - Retrieve network adapter information (IP, gateway, etc.)
- `Registry API` - Read friendly adapter names from registry
- `AllocateAndInitializeSid` - Check administrator privileges
- `netsh` - Set IP configuration (hidden window execution)

### Key Code Snippets
```go
// 1. Get friendly adapter name
func GetFriendlyName(adapterGUID string) string {
    keyPath := `SYSTEM\CurrentControlSet\Control\Network\...`
    key, _ := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.QUERY_VALUE)
    name, _, _ := key.GetStringValue("Name")
    return name
}

// 2. Execute commands with hidden window
cmd.SysProcAttr = &syscall.SysProcAttr{
    HideWindow:    true,
    CreationFlags: CREATE_NO_WINDOW,
}
```

## 🎯 Advantages

### vs PowerShell Scripts
- ✅ No encoding issues, perfect Chinese support
- ✅ Graphical interface, more intuitive
- ✅ Persistent configuration, no manual code editing
- ✅ Support multiple adapters and configurations
- ✅ No window popups, better user experience

### vs Other Tools
- ✅ Single file, no installation needed
- ✅ Open source and free, transparent code
- ✅ Lightweight, small size (<10MB)
- ✅ Fast startup, smooth operation
- ✅ Uses Windows native API, optimal performance

## 🔐 Security

- ✅ Only modifies local network configuration
- ✅ No network communication involved
- ✅ All code is open source and auditable
- ✅ Configuration files stored locally
- ✅ Administrator rights required by Windows security
- ✅ Uses official Windows API, safe and reliable

## 📝 Changelog

### v2.1 (2024-12-04)
- ✨ Use Windows API instead of PowerShell, no window popups
- ✨ Read friendly adapter names from registry
- ✨ Fix configuration list real-time update issue
- ✨ Fix startup auto-detection of configuration
- ✨ Optimize details display format (proper line breaks)
- 🐛 Fix netsh command execution failure

### v2.0 (2024-12-04)
- ✨ Complete redesign, support adapter selection
- ✨ Support add, edit, delete configurations
- ✨ Persistent configuration storage
- ✨ Support managing multiple adapters
- ✨ Optimize UI and interaction

### v1.0 (2024-12-04)
- 🎉 Initial release
- ✅ Basic configuration switching functionality

## 📄 License

MIT License - Free to use, modify, and distribute

## 👨‍💻 Contributing

Issues and Pull Requests are welcome!

---

**Project Name Suggestions**: quick-ip-switcher, network-config-manager, ip-switcher-win
**Recommended**: `quick-ip-switcher`
