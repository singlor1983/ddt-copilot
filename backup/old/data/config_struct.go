package data

import (
	"ddt-copilot/defs"
	"github.com/lxn/win"
	"sync"
)

type SettingCustomFuben struct {
	Name    string         // 副本名
	Tp      defs.FubenType // 副本类型
	Page    int            // 第几页
	Count   int            // 这一页的第几个
	LvCount int            // 有几个难度
	LvIndex int            // 选第几个难度
}

type SettingLoopFuben struct {
}

type SettingCommon struct {
	IsBossFightEnable bool
	AttackCmd         []string // 自定义攻击指令
}

type SettingSpFuben struct {
}

type SettingMouse struct {
}

type SettingHWND struct {
	HWND   int32
	IsStop bool
}

type SettingHWNDS struct {
	HWNDS map[int32]SettingHWND
}

type SettingFubenPostion struct {
	Position map[defs.FubenID]defs.FubenPosition
}

type Setting struct {
	SettingCustomFuben  SettingCustomFuben  // 自定义副本设置区
	SettingLoopFuben    SettingLoopFuben    // 轮刷副本设置区
	SettingCommon       SettingCommon       // 通用设置区
	SettingSpFuben      SettingSpFuben      // 副本特殊设置区
	SettingMouse        SettingMouse        // 鼠标连点功能设置区
	SettingHWNDS        SettingHWNDS        // 窗口句柄对应的设置
	SettingFubenPostion SettingFubenPostion // 第一次打开页面时需要初始化这个副本位置
}

type Config struct {
	mu sync.RWMutex

	setting Setting
}

// Save 序列化，保存为文件
func (c *Config) Save() {

}

// Load 从文件中加载
func (c *Config) Load() {

}

func (c *Config) IsStop(hwnd win.HWND) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.setting.SettingHWNDS.HWNDS == nil {
		return true
	}
	v, ok := c.setting.SettingHWNDS.HWNDS[int32(hwnd)]
	if !ok {
		return true
	}
	return v.IsStop
}

func (c *Config) StopHWND(hwnd win.HWND) {
	c.mu.Lock()
	defer c.mu.Unlock()

}

func (c *Config) StartHWND(hwnd win.HWND) {
	c.mu.Lock()
	defer c.mu.Unlock()
}

func (c *Config) IsBossFightEnable() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.setting.SettingCommon.IsBossFightEnable
}

func (c *Config) SetBossFightEnable(enable bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.setting.SettingCommon.IsBossFightEnable = enable
}

func (c *Config) SetFubenPosition(position SettingFubenPostion) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.setting.SettingFubenPostion = position
}

func (c *Config) GetFubenPosition() map[defs.FubenID]defs.FubenPosition {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.setting.SettingFubenPostion.Position
}

func (c *Config) SetAttackCmd(cmds []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.setting.SettingCommon.AttackCmd = cmds
}

func (c *Config) GetAttackCmd() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.setting.SettingCommon.AttackCmd
}
