package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"golang.org/x/sys/windows"
)

var (
	appConfig *Config
	// Version 由编译时通过 -ldflags 注入
	Version = "dev"
)

func main() {
	// 检查管理员权限
	if !IsAdmin() {
		// 直接尝试以管理员身份运行,不弹窗询问
		if err := runAsAdmin(); err != nil {
			// 如果提升失败,显示错误并继续(以普通用户身份运行,只能查看不能修改)
			log.Printf("无法提升权限: %v\n", err)
		} else {
			// 提升成功,新进程已启动,当前进程退出
			return
		}
	}

	// 加载配置
	var err error
	appConfig, err = LoadConfig()
	if err != nil {
		// 如果加载失败,使用默认配置
		log.Printf("加载配置失败: %v,使用默认配置\n", err)
		appConfig = &Config{
			SelectedInterface: "",
			NetworkConfigs:    []NetworkConfig{},
		}
	}

	// 启动GUI
	if err := runMainWindow(); err != nil {
		// GUI启动失败
		walk.MsgBox(nil, "错误", fmt.Sprintf("启动失败: %v", err), walk.MsgBoxIconError)
		log.Fatal(err)
	}
}

// runAsAdmin 以管理员身份重新运行
func runAsAdmin() error {
	var (
		shell32  = windows.NewLazySystemDLL("shell32.dll")
		shellExe = shell32.NewProc("ShellExecuteW")
	)

	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := ""

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	ret, _, _ := shellExe.Call(
		0,
		uintptr(unsafe.Pointer(verbPtr)),
		uintptr(unsafe.Pointer(exePtr)),
		uintptr(unsafe.Pointer(argPtr)),
		uintptr(unsafe.Pointer(cwdPtr)),
		1, // SW_NORMAL
	)

	if ret <= 32 {
		return fmt.Errorf("无法以管理员身份运行(错误码: %d)", ret)
	}
	return nil
}

// runMainWindow 运行主窗口
func runMainWindow() error {
	var mainWindow *walk.MainWindow
	var interfaceCombo *walk.ComboBox
	var configList *walk.ListBox
	var statusLabel *walk.Label
	var detailsEdit *walk.TextEdit

	// 配置列表模型 - 使用指针,使其能实时反映appConfig的变化
	configModel := &ConfigListModel{configs: &appConfig.NetworkConfigs}

	// 更新状态信息
	updateStatus := func() {
		if appConfig.SelectedInterface == "" {
			statusLabel.SetText("请先选择网络适配器")
			detailsEdit.SetText("未选择网络适配器")
			return
		}

		configName := GetCurrentConfigName(appConfig.SelectedInterface, appConfig.NetworkConfigs)
		current := GetCurrentConfig(appConfig.SelectedInterface)

		statusLabel.SetText(fmt.Sprintf("当前配置: %s", configName))

		if current.Exists {
			details := fmt.Sprintf("网络适配器: %s\r\n\r\nIP 地址    : %s\r\n默认网关   : %s\r\nDNS 服务器 : %s",
				appConfig.SelectedInterface, current.IP, current.Gateway, current.DNS)
			detailsEdit.SetText(details)
		} else {
			detailsEdit.SetText(fmt.Sprintf("错误: %s", current.Error))
		}
	}

	// 切换配置
	switchConfig := func(index int) {
		if appConfig.SelectedInterface == "" {
			walk.MsgBox(mainWindow, "错误", "请先选择网络适配器", walk.MsgBoxIconError)
			return
		}

		config, err := appConfig.GetNetworkConfig(index)
		if err != nil {
			walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
			return
		}

		err = ApplyConfig(appConfig.SelectedInterface, *config)
		if err != nil {
			walk.MsgBox(mainWindow, "配置失败", err.Error(), walk.MsgBoxIconError)
			return
		}

		walk.MsgBox(mainWindow, "配置成功",
			fmt.Sprintf("配置已成功切换到: %s", config.Name),
			walk.MsgBoxIconInformation)

		updateStatus()
	}

	// 初始化标志
	var initialized bool

	// 创建窗口
	title := fmt.Sprintf("网络配置管理工具 %s", Version)
	_, err := MainWindow{
		AssignTo: &mainWindow,
		Title:    title,
		MinSize:  Size{Width: 600, Height: 500},
		Size:     Size{Width: 600, Height: 500},
		Layout:   VBox{Margins: Margins{Left: 10, Top: 10, Right: 10, Bottom: 10}},
		OnSizeChanged: func() {
			// 窗口首次显示时加载数据(只执行一次)
			if !initialized && interfaceCombo != nil {
				initialized = true
				loadAdapters(interfaceCombo)
				updateStatus()
			}
		},
		Children: []Widget{
			Label{
				Text: "网络配置管理工具",
				Font: Font{PointSize: 16, Bold: true},
			},
			VSpacer{Size: 10},

			// 网卡选择
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text:    "网络适配器:",
						MinSize: Size{Width: 80},
					},
					ComboBox{
						AssignTo:      &interfaceCombo,
						Editable:      false,
						MinSize:       Size{Width: 300},
						OnCurrentIndexChanged: func() {
							if interfaceCombo.CurrentIndex() >= 0 {
								adapters, _ := GetNetworkAdapters()
								if interfaceCombo.CurrentIndex() < len(adapters) {
									appConfig.SelectedInterface = adapters[interfaceCombo.CurrentIndex()].Name
									SaveConfig(appConfig)
									updateStatus()
								}
							}
						},
					},
					PushButton{
						Text: "刷新",
						OnClicked: func() {
							loadAdapters(interfaceCombo)
							updateStatus()
						},
					},
				},
			},
			VSpacer{Size: 5},

			// 当前状态
			Label{
				AssignTo: &statusLabel,
				Text:     "当前配置: 未检测",
				Font:     Font{PointSize: 10, Bold: true},
			},
			VSpacer{Size: 5},

			TextEdit{
				AssignTo: &detailsEdit,
				ReadOnly: true,
				MinSize:  Size{Height: 120},
				MaxSize:  Size{Height: 120},
				Font:     Font{Family: "Consolas", PointSize: 9},
			},
			VSpacer{Size: 10},

			// 配置管理
			Label{
				Text: "保存的配置:",
				Font: Font{PointSize: 10, Bold: true},
			},
			VSpacer{Size: 5},

			Composite{
				Layout: HBox{},
				Children: []Widget{
					ListBox{
						AssignTo: &configList,
						Model:    configModel,
						MinSize:  Size{Width: 400, Height: 150},
						OnItemActivated: func() {
							if configList.CurrentIndex() >= 0 {
								switchConfig(configList.CurrentIndex())
							}
						},
					},
					Composite{
						Layout: VBox{},
						Children: []Widget{
							PushButton{
								Text:    "切换",
								MinSize: Size{Width: 120, Height: 35},
								OnClicked: func() {
									if configList.CurrentIndex() >= 0 {
										switchConfig(configList.CurrentIndex())
									} else {
										walk.MsgBox(mainWindow, "提示", "请先选择一个配置", walk.MsgBoxIconInformation)
									}
								},
							},
							VSpacer{Size: 5},
							PushButton{
								Text:    "新建配置",
								MinSize: Size{Width: 120, Height: 35},
								OnClicked: func() {
									showConfigDialog(mainWindow, nil, -1, func() {
										configModel.PublishItemsReset()
										configList.SetCurrentIndex(len(appConfig.NetworkConfigs) - 1)
									})
								},
							},
							VSpacer{Size: 5},
							PushButton{
								Text:    "编辑",
								MinSize: Size{Width: 120, Height: 35},
								OnClicked: func() {
									if configList.CurrentIndex() >= 0 {
										index := configList.CurrentIndex()
										config, _ := appConfig.GetNetworkConfig(index)
										showConfigDialog(mainWindow, config, index, func() {
											configModel.PublishItemsReset()
										})
									} else {
										walk.MsgBox(mainWindow, "提示", "请先选择一个配置", walk.MsgBoxIconInformation)
									}
								},
							},
							VSpacer{Size: 5},
							PushButton{
								Text:    "删除",
								MinSize: Size{Width: 120, Height: 35},
								OnClicked: func() {
									if configList.CurrentIndex() >= 0 {
										ret := walk.MsgBox(mainWindow, "确认删除",
											"确定要删除选中的配置吗?",
											walk.MsgBoxYesNo|walk.MsgBoxIconQuestion)
										if ret == walk.DlgCmdYes {
											appConfig.RemoveNetworkConfig(configList.CurrentIndex())
											SaveConfig(appConfig)
											configModel.PublishItemsReset()
										}
									} else {
										walk.MsgBox(mainWindow, "提示", "请先选择一个配置", walk.MsgBoxIconInformation)
									}
								},
							},
							VSpacer{},
						},
					},
				},
			},
			VSpacer{Size: 10},

			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						Text:    "刷新状态",
						MinSize: Size{Width: 100, Height: 35},
						OnClicked: func() {
							updateStatus()
						},
					},
				},
			},
		},
	}.Run()

	return err
}

// loadAdapters 加载网络适配器列表
func loadAdapters(combo *walk.ComboBox) {
	adapters, err := GetNetworkAdapters()
	if err != nil {
		walk.MsgBox(nil, "错误", fmt.Sprintf("获取网络适配器失败: %v", err), walk.MsgBoxIconError)
		return
	}

	items := make([]string, 0, len(adapters))
	selectedIndex := -1
	for i, adapter := range adapters {
		item := fmt.Sprintf("%s (%s)", adapter.Name, adapter.Status)
		items = append(items, item)
		if adapter.Name == appConfig.SelectedInterface {
			selectedIndex = i
		}
	}

	combo.SetModel(items)
	if selectedIndex >= 0 {
		combo.SetCurrentIndex(selectedIndex)
	}
}

// ConfigListModel 配置列表模型
type ConfigListModel struct {
	walk.ListModelBase
	configs *[]NetworkConfig
}

func (m *ConfigListModel) ItemCount() int {
	if m.configs == nil {
		return 0
	}
	return len(*m.configs)
}

func (m *ConfigListModel) Value(index int) interface{} {
	if m.configs == nil || index < 0 || index >= len(*m.configs) {
		return ""
	}
	cfg := (*m.configs)[index]
	return fmt.Sprintf("%s (IP: %s, 网关: %s)", cfg.Name, cfg.IP, cfg.Gateway)
}

// showConfigDialog 显示配置编辑对话框
func showConfigDialog(owner walk.Form, config *NetworkConfig, index int, onSave func()) {
	var dlg *walk.Dialog
	var nameEdit, ipEdit, maskEdit, gatewayEdit, dnsEdit *walk.LineEdit
	var acceptPB, cancelPB *walk.PushButton

	isNew := config == nil
	if isNew {
		config = &NetworkConfig{}
	}

	Dialog{
		AssignTo:  &dlg,
		Title:     func() string { if isNew { return "新建配置" } else { return "编辑配置" } }(),
		MinSize:   Size{Width: 400, Height: 300},
		Layout:    VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "配置名称:"},
					LineEdit{AssignTo: &nameEdit, Text: config.Name},

					Label{Text: "IP地址:"},
					LineEdit{AssignTo: &ipEdit, Text: config.IP},

					Label{Text: "子网掩码:"},
					LineEdit{AssignTo: &maskEdit, Text: config.SubnetMask},

					Label{Text: "默认网关:"},
					LineEdit{AssignTo: &gatewayEdit, Text: config.Gateway},

					Label{Text: "DNS服务器:"},
					LineEdit{AssignTo: &dnsEdit, Text: config.DNS},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "保存",
						OnClicked: func() {
							// 验证输入
							if nameEdit.Text() == "" {
								walk.MsgBox(dlg, "错误", "配置名称不能为空", walk.MsgBoxIconError)
								return
							}
							if ipEdit.Text() == "" {
								walk.MsgBox(dlg, "错误", "IP地址不能为空", walk.MsgBoxIconError)
								return
							}

							// 保存配置
							newConfig := NetworkConfig{
								Name:       nameEdit.Text(),
								IP:         ipEdit.Text(),
								SubnetMask: maskEdit.Text(),
								Gateway:    gatewayEdit.Text(),
								DNS:        dnsEdit.Text(),
							}

							if isNew {
								appConfig.AddNetworkConfig(newConfig)
							} else {
								appConfig.UpdateNetworkConfig(index, newConfig)
							}

							if err := SaveConfig(appConfig); err != nil {
								walk.MsgBox(dlg, "错误", fmt.Sprintf("保存配置失败: %v", err), walk.MsgBoxIconError)
								return
							}

							if onSave != nil {
								onSave()
							}

							dlg.Accept()
						},
					},
					PushButton{
						AssignTo: &cancelPB,
						Text:     "取消",
						OnClicked: func() {
							dlg.Cancel()
						},
					},
				},
			},
		},
	}.Run(owner)
}
