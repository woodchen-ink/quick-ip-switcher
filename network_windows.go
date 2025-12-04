// +build windows

package main

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var (
	iphlpapi                 = windows.NewLazySystemDLL("iphlpapi.dll")
	procGetAdaptersInfo      = iphlpapi.NewProc("GetAdaptersInfo")
	procGetAdaptersAddresses = iphlpapi.NewProc("GetAdaptersAddresses")
)

const (
	MAX_ADAPTER_NAME_LENGTH        = 256
	MAX_ADAPTER_DESCRIPTION_LENGTH = 128
	MAX_ADAPTER_ADDRESS_LENGTH     = 8
	ERROR_BUFFER_OVERFLOW          = 111
	CREATE_NO_WINDOW               = 0x08000000
)

type IP_ADAPTER_INFO struct {
	Next                *IP_ADAPTER_INFO
	ComboIndex          uint32
	AdapterName         [MAX_ADAPTER_NAME_LENGTH + 4]byte
	Description         [MAX_ADAPTER_DESCRIPTION_LENGTH + 4]byte
	AddressLength       uint32
	Address             [MAX_ADAPTER_ADDRESS_LENGTH]byte
	Index               uint32
	Type                uint32
	DhcpEnabled         uint32
	CurrentIpAddress    *IP_ADDR_STRING
	IpAddressList       IP_ADDR_STRING
	GatewayList         IP_ADDR_STRING
	DhcpServer          IP_ADDR_STRING
	HaveWins            uint32
	PrimaryWinsServer   IP_ADDR_STRING
	SecondaryWinsServer IP_ADDR_STRING
	LeaseObtained       int64
	LeaseExpires        int64
}

type IP_ADDR_STRING struct {
	Next      *IP_ADDR_STRING
	IpAddress [16]byte
	IpMask    [16]byte
	Context   uint32
}

// GetFriendlyName 从注册表获取网卡的友好名称
func GetFriendlyName(adapterGUID string) string {
	keyPath := `SYSTEM\CurrentControlSet\Control\Network\{4D36E972-E325-11CE-BFC1-08002BE10318}\` + adapterGUID + `\Connection`

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.QUERY_VALUE)
	if err != nil {
		return adapterGUID
	}
	defer key.Close()

	name, _, err := key.GetStringValue("Name")
	if err != nil {
		return adapterGUID
	}

	return name
}

// GetNetworkAdaptersWinAPI 使用Windows API获取网络适配器
func GetNetworkAdaptersWinAPI() ([]NetworkAdapter, error) {
	var size uint32 = 15000
	buf := make([]byte, size)

	ret, _, _ := procGetAdaptersInfo.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
	)

	if ret == ERROR_BUFFER_OVERFLOW {
		buf = make([]byte, size)
		ret, _, _ = procGetAdaptersInfo.Call(
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&size)),
		)
	}

	if ret != 0 {
		return nil, fmt.Errorf("GetAdaptersInfo failed with code %d", ret)
	}

	var adapters []NetworkAdapter
	adapter := (*IP_ADAPTER_INFO)(unsafe.Pointer(&buf[0]))

	for adapter != nil {
		// 获取GUID
		guid := windows.BytePtrToString(&adapter.AdapterName[0])

		// 从注册表获取友好名称
		friendlyName := GetFriendlyName(guid)

		// 获取描述
		desc := windows.BytePtrToString(&adapter.Description[0])

		// 判断状态
		ipAddr := windows.BytePtrToString(&adapter.IpAddressList.IpAddress[0])
		status := "Disconnected"
		if ipAddr != "" && ipAddr != "0.0.0.0" {
			status = "Up"
		}

		adapters = append(adapters, NetworkAdapter{
			Name:        friendlyName, // 使用友好名称
			Status:      status,
			Description: desc,
		})

		adapter = adapter.Next
	}

	return adapters, nil
}

// GetCurrentConfigWinAPI 使用Windows API获取当前配置
func GetCurrentConfigWinAPI(interfaceName string) CurrentNetworkInfo {
	info := CurrentNetworkInfo{Exists: true}

	var size uint32 = 15000
	buf := make([]byte, size)

	ret, _, _ := procGetAdaptersInfo.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
	)

	if ret == ERROR_BUFFER_OVERFLOW {
		buf = make([]byte, size)
		ret, _, _ = procGetAdaptersInfo.Call(
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&size)),
		)
	}

	if ret != 0 {
		info.Exists = false
		info.Error = fmt.Sprintf("获取适配器信息失败: %d", ret)
		return info
	}

	adapter := (*IP_ADAPTER_INFO)(unsafe.Pointer(&buf[0]))

	for adapter != nil {
		guid := windows.BytePtrToString(&adapter.AdapterName[0])
		friendlyName := GetFriendlyName(guid)

		if friendlyName == interfaceName {
			// 获取IP地址
			info.IP = windows.BytePtrToString(&adapter.IpAddressList.IpAddress[0])

			// 获取网关
			info.Gateway = windows.BytePtrToString(&adapter.GatewayList.IpAddress[0])

			// 获取DNS (通过PowerShell但隐藏窗口)
			cmd := exec.Command("powershell", "-NoProfile", "-WindowStyle", "Hidden", "-Command",
				fmt.Sprintf("((Get-DnsClientServerAddress -InterfaceAlias '%s' -AddressFamily IPv4 -ErrorAction SilentlyContinue).ServerAddresses -join ', ')", interfaceName))
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: CREATE_NO_WINDOW}
			output, err := cmd.Output()
			if err == nil && len(output) > 0 {
				info.DNS = strings.TrimSpace(string(output))
			}

			if info.IP == "" || info.IP == "0.0.0.0" {
				info.Exists = false
				info.Error = "适配器未配置IP地址"
			}

			return info
		}

		adapter = adapter.Next
	}

	info.Exists = false
	info.Error = "未找到指定的网络适配器"
	return info
}

// SetIPConfigWinAPI 使用netsh设置IP配置
func SetIPConfigWinAPI(interfaceName string, config NetworkConfig) error {
	// 设置静态IP
	cmd := fmt.Sprintf(`netsh interface ip set address name="%s" static %s %s %s`,
		interfaceName, config.IP, config.SubnetMask, config.Gateway)

	if err := execCmd(cmd); err != nil {
		return fmt.Errorf("设置IP地址失败: %v", err)
	}

	// 设置DNS
	cmd = fmt.Sprintf(`netsh interface ip set dns name="%s" static %s`,
		interfaceName, config.DNS)

	if err := execCmd(cmd); err != nil {
		return fmt.Errorf("设置DNS失败: %v", err)
	}

	return nil
}

// execCmd 执行命令
func execCmd(command string) error {
	cmd := exec.Command("cmd", "/C", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: CREATE_NO_WINDOW,
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v (输出: %s)", err, string(output))
	}

	return nil
}

// IsAdminWinAPI 使用Windows API检查管理员权限
func IsAdminWinAPI() bool {
	var sid *windows.SID

	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)
	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}

	return member
}
