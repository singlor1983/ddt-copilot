package utils

import (
	"ddt-copilot/defs"
	"github.com/lxn/win"
	"image"
	"time"
)

// GetDDTHwnds 获取所有的弹弹堂游戏窗口
func GetDDTHwnds() []win.HWND {
	var ret []win.HWND
	pids, _ := GetProcessID(string(defs.ProcessTgWeb))
	for _, pid := range pids {
		hwnd, err := GetFirstWindowByPID(pid)
		if err != nil {
			continue
		}
		wds, _ := GetAllChildWindows(hwnd)
		if len(wds) != 5 {
			continue
		}
		lastWd := win.HWND(wds[len(wds)-1])
		ret = append(ret, lastWd)
	}
	return ret
}

func FocusDDTWindow(hwnd win.HWND, dropBlock bool) {
	point := defs.GetPoint(defs.PointEmpty)
	LeftClick(hwnd, point.X, point.Y, 0) // 激活窗口，便于截图为亮色的
	if dropBlock {                       // 副本选择页面不能删除遮挡，因为这本身就是遮挡
		for i := 0; i < 5; i++ { // 截屏之前用ESC把其他遮挡界面关闭，避免影响判断，5次是因为可能有多层折叠
			KeyBoard(hwnd, defs.VK_ESCAPE, 0)
		}
	}
	time.Sleep(time.Millisecond * 100) // 等待那个继续操作图标消失
}

func CaptureWindowLight(hwnd win.HWND, captureRect *win.RECT, dropBlock bool) (*image.RGBA, error) {
	FocusDDTWindow(hwnd, dropBlock)
	return CaptureWindow(hwnd, captureRect)
}
