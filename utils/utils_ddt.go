package utils

import (
	"ddt-copilot/defs"
	"github.com/lxn/win"
	"image"
	"math"
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

func CaptureWindowLightWithGray(hwnd win.HWND, captureRect *win.RECT, dropBlock bool) (*image.Gray, error) {
	img, err := CaptureWindowLight(hwnd, captureRect, dropBlock)
	if err != nil {
		return nil, err
	}
	return ConvertToGray(img), nil
}

func CaptureWindowLightWithNormalization(hwnd win.HWND, captureRect *win.RECT, dropBlock bool) (*image.Gray, error) {
	img, err := CaptureWindowLight(hwnd, captureRect, dropBlock)
	if err != nil {
		return nil, err
	}
	return ConvertToGrayWithNormalization(img), nil
}

// UpdateAngle 改变力度
func UpdateAngle(hwnd win.HWND, angle int) {
	direction := defs.DirectionUp
	if angle < 0 {
		direction = defs.DirectionDown
		angle = int(math.Abs(float64(angle)))
	}
	switch direction {
	case defs.DirectionUp:
		for i := 0; i < angle; i++ {
			KeyBoard(hwnd, defs.VK_UP, 0)
		}
	case defs.DirectionDown:
		for i := 0; i < angle; i++ {
			KeyBoard(hwnd, defs.VK_DOWN, 0)
		}
	}
}

// ConfirmDirection 转向、确认方向
func ConfirmDirection(hwnd win.HWND, direction defs.Direction) {
	switch direction {
	case defs.DirectionLeft:
		KeyBoard(hwnd, defs.VK_RIGHT, 0)
		KeyBoard(hwnd, defs.VK_LEFT, 0)
	case defs.DirectionRight:
		KeyBoard(hwnd, defs.VK_LEFT, 0)
		KeyBoard(hwnd, defs.VK_RIGHT, 0)
	}
}

// Move 移动。1格距离100毫秒
func Move(hwnd win.HWND, direction defs.Direction, distance int) {
	ConfirmDirection(hwnd, direction)
	if distance < 0 {
		distance = 0
	}
	ts := time.Duration(distance*100) * time.Millisecond
	switch direction {
	case defs.DirectionLeft:
		KeyBoard(hwnd, defs.VK_LEFT, ts)
	case defs.DirectionRight:
		KeyBoard(hwnd, defs.VK_RIGHT, ts)
	}
}

func UseSkill(hwnd win.HWND, skill uintptr) {
	KeyBoard(hwnd, skill, 0)
}

// Launch 攻击、发射，根据力度计算出需要按压持续的时间。1度40毫秒
func Launch(hwnd win.HWND, power int) {
	if power < 0 {
		power = 0
	}
	if power > 100 {
		power = 100
	}
	ts := time.Duration(power*40) * time.Millisecond
	KeyBoard(hwnd, defs.VK_SPACE, ts)
}
