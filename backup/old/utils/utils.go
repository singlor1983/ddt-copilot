package utils

import (
	"errors"
	"fmt"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
	"image"
	"image/color"
	"image/draw"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

// 定义 Windows API 函数
var (
	user32 = windows.NewLazySystemDLL("user32.dll")

	procEnumChildWindows         = user32.NewProc("EnumChildWindows")
	procIsWindow                 = user32.NewProc("IsWindow")
	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
)

// GetProcessID 根据进程名获取进程ID
func GetProcessID(processName string) ([]uint32, error) {
	var processID []uint32

	// 创建进程快照
	processSnap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(processSnap)

	var processEntry windows.ProcessEntry32
	processEntry.Size = uint32(unsafe.Sizeof(processEntry))

	// 遍历进程
	for {
		if err = windows.Process32Next(processSnap, &processEntry); err != nil {
			break
		}

		exeFileName := windows.UTF16ToString(processEntry.ExeFile[:])
		if strings.EqualFold(exeFileName, processName) {
			processID = append(processID, processEntry.ProcessID)
		}
	}

	if len(processID) == 0 {
		return nil, fmt.Errorf("process:%s not found", processName)
	}
	return processID, nil
}

func GetWindowsByPID(pid uint32) ([]windows.Handle, error) {
	var wds []windows.Handle

	callback := func(hwnd windows.Handle, lParam uintptr) uintptr {
		var curPid uint32
		procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&curPid)))

		if curPid == uint32(lParam) {
			wds = append(wds, hwnd)
		}
		return 1 // 继续枚举
	}

	// 调用 EnumWindows
	_, _, _ = procEnumWindows.Call(syscall.NewCallback(callback), uintptr(pid))

	return wds, nil
}

func GetFirstWindowByPID(pid uint32) (windows.Handle, error) {
	wds, err := GetWindowsByPID(pid)
	if err != nil {
		return 0, err
	}
	if len(wds) == 0 {
		return 0, fmt.Errorf("not found windows. pid:%d", pid)
	}
	return wds[0], nil
}

func GetAllChildWindows(parentHwnd windows.Handle) ([]windows.Handle, error) {
	// 检查父窗口有效性
	isValid, _, _ := procIsWindow.Call(uintptr(parentHwnd))
	if isValid == 0 {
		return nil, fmt.Errorf("parent Window Handle %v is not valid", parentHwnd)
	}

	var childWindows []windows.Handle

	callback := func(hwnd windows.Handle, lParam uintptr) uintptr {
		childWindows = append(childWindows, hwnd)

		var pid uint32
		procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))

		return 1 // 继续枚举
	}

	procEnumChildWindows.Call(uintptr(parentHwnd), syscall.NewCallback(callback), 0)

	return childWindows, nil
}

func createImage(rect image.Rectangle) (img *image.RGBA, e error) {
	img = nil
	e = errors.New("cannot create image.RGBA")

	defer func() {
		err := recover()
		if err == nil {
			e = nil
		}
	}()
	// image.NewRGBA may panic if rect is too large.
	img = image.NewRGBA(rect)

	return img, e
}

// ExtractRegion 从原始图像中提取指定区域
func ExtractRegion(img *image.RGBA, captureRect *win.RECT) *image.RGBA {
	// 计算截取区域的宽度和高度
	width := int(captureRect.Right - captureRect.Left)
	height := int(captureRect.Bottom - captureRect.Top)

	// 确保 captureRect 在原始图像的边界内
	if captureRect.Left < 0 {
		captureRect.Left = 0
	}
	if captureRect.Top < 0 {
		captureRect.Top = 0
	}
	if captureRect.Right > int32(img.Rect.Max.X) {
		captureRect.Right = int32(img.Rect.Max.X)
	}
	if captureRect.Bottom > int32(img.Rect.Max.Y) {
		captureRect.Bottom = int32(img.Rect.Max.Y)
	}

	// 创建一个新的 RGBA 图像
	extractedImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// 定义源矩形和目标矩形
	srcRect := image.Rect(int(captureRect.Left), int(captureRect.Top), int(captureRect.Right), int(captureRect.Bottom))
	dstRect := image.Rect(0, 0, width, height)

	// 使用 draw 包将指定区域复制到新的图像
	draw.Draw(extractedImg, dstRect, img, srcRect.Min, draw.Src)

	return extractedImg
}

// CaptureWindow 截取指定窗口的某一部分或整个窗口[截屏之前可以先激活窗口，以便截图的是亮色的图片，或者先模拟点击一下空白处]
func CaptureWindow(hwnd win.HWND, captureRect *win.RECT) (*image.RGBA, error) {
	var windRect win.RECT
	win.GetClientRect(hwnd, &windRect)

	var rect win.RECT
	if captureRect != nil {
		// 使用传入的 RECT，如果为 nil 则使用整个窗口
		rect = *captureRect

		// 确保截取区域在窗口范围内
		if rect.Left < 0 || rect.Top < 0 ||
			rect.Right > windRect.Right || rect.Bottom > windRect.Bottom {
			return nil, errors.New("capture rectangle is out of bounds")
		}
	} else {
		// 截取整个窗口
		rect = windRect
	}

	// 计算宽度和高度
	width := int(rect.Right - rect.Left)
	height := int(rect.Bottom - rect.Top)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	hdc := win.GetDC(hwnd)
	if hdc == 0 {
		return nil, errors.New("GetDC failed")
	}
	defer win.ReleaseDC(hwnd, hdc)

	memoryDevice := win.CreateCompatibleDC(hdc)
	if memoryDevice == 0 {
		return nil, errors.New("CreateCompatibleDC failed")
	}
	defer win.DeleteDC(memoryDevice)

	bitmap := win.CreateCompatibleBitmap(hdc, int32(width), int32(height))
	if bitmap == 0 {
		return nil, errors.New("CreateCompatibleBitmap failed")
	}
	defer win.DeleteObject(win.HGDIOBJ(bitmap))

	old := win.SelectObject(memoryDevice, win.HGDIOBJ(bitmap))
	if old == 0 {
		return nil, errors.New("SelectObject failed")
	}
	defer win.SelectObject(memoryDevice, old)

	win.SendMessage(hwnd, win.WM_PAINT, 0, 0)
	// 使用 BitBlt 来捕捉指定区域或整个窗口
	if !win.BitBlt(memoryDevice, 0, 0, int32(width), int32(height), hdc, rect.Left, rect.Top, win.SRCCOPY) {
		return nil, errors.New("BitBlt failed")
	}

	var header win.BITMAPINFOHEADER
	header.BiSize = uint32(unsafe.Sizeof(header))
	header.BiPlanes = 1
	header.BiBitCount = 32
	header.BiWidth = int32(width)
	header.BiHeight = int32(-height) // 负值表示从顶部开始填充
	header.BiCompression = win.BI_RGB
	header.BiSizeImage = 0

	bitmapDataSize := uintptr(((int64(width)*int64(header.BiBitCount) + 31) / 32) * 4 * int64(height))
	hmem := win.GlobalAlloc(win.GMEM_MOVEABLE, bitmapDataSize)
	defer win.GlobalFree(hmem)
	memptr := win.GlobalLock(hmem)
	defer win.GlobalUnlock(hmem)

	if win.GetDIBits(hdc, bitmap, 0, uint32(height), (*uint8)(memptr), (*win.BITMAPINFO)(unsafe.Pointer(&header)), win.DIB_RGB_COLORS) == 0 {
		return nil, errors.New("GetDIBits failed")
	}

	i := 0
	src := uintptr(memptr)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			v0 := *(*uint8)(unsafe.Pointer(src))
			v1 := *(*uint8)(unsafe.Pointer(src + 1))
			v2 := *(*uint8)(unsafe.Pointer(src + 2))

			// BGRA => RGBA, and set A to 255
			img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3] = v2, v1, v0, 255

			i += 4
			src += 4
		}
	}

	return img, nil
}

// ConvertToGray 将 RGBA 图像转换为灰度图像
func ConvertToGray(img *image.RGBA) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// 使用加权平均法计算灰度值
			grayValue := (r*299 + g*587 + b*114) / 1000
			grayImg.Set(x, y, color.Gray{Y: uint8(grayValue >> 8)}) // 右移以适应 0-255 范围
		}
	}

	return grayImg
}

// ConvertToGrayWithNormalization 转成灰度图之后再归一化，适用于数字归一化
func ConvertToGrayWithNormalization(img *image.RGBA) *image.Gray {
	imgGray := ConvertToGray(img)
	bounds := imgGray.Bounds()
	newImg := image.NewGray(bounds)
	targetColor := color.Gray{Y: 128} // 指定要变成的颜色（例如纯灰色）
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayValue := imgGray.At(x, y).(color.Gray).Y
			if grayValue == 0 { // 黑色像素
				newImg.Set(x, y, color.Gray{Y: 0})
			} else {
				newImg.Set(x, y, targetColor)
			}
		}
	}
	return newImg
}

// CompareGrayImages 计算两个相同尺寸的灰度图像的均方误差（MSE）
func CompareGrayImages(img1, img2 *image.Gray) (float64, error) {
	if img1.Bounds() != img2.Bounds() {
		return 0, errors.New("images must have the same dimensions")
	}

	var sumSquaredDiff float64
	width, height := img1.Bounds().Max.X, img1.Bounds().Max.Y
	pixelCount := float64(width * height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gray1 := img1.GrayAt(x, y).Y
			gray2 := img2.GrayAt(x, y).Y
			diff := float64(gray1) - float64(gray2)
			sumSquaredDiff += diff * diff
		}
	}

	mse := sumSquaredDiff / pixelCount
	return mse, nil
}

// IsSimilarity 计算相似性比例，可以定义一个阈值来判断相似性
func IsSimilarity(mse float64, threshold float64) bool {
	return mse <= threshold
}

func Click(hWnd win.HWND, x int32, y int32, duration time.Duration) {
	lParam := uintptr(y<<16 | x)
	win.PostMessage(hWnd, win.WM_LBUTTONDOWN, win.MK_LBUTTON, lParam)
	time.Sleep(duration)
	win.PostMessage(hWnd, win.WM_LBUTTONUP, 0, lParam)
}

func KeyBoard(hWnd win.HWND, key uintptr, duration time.Duration) {
	win.PostMessage(hWnd, win.WM_KEYDOWN, key, 0)
	time.Sleep(duration)
	win.PostMessage(hWnd, win.WM_KEYUP, key, 0)
}

// findColorPixels 查找指定颜色的所有像素点
func findColorPixels(img *image.RGBA, targetColor color.RGBA, tolerance int) []image.Point {
	var pixels []image.Point
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if isSimilarColorWithRange(img.RGBAAt(x, y), targetColor, tolerance) {
				pixels = append(pixels, image.Point{X: x, Y: y})
			}
		}
	}

	return pixels
}

// isSimilarColor 检查两个颜色之间的相似度
func isSimilarColorWithRange(c color.Color, target color.RGBA, tolerance int) bool {
	r, g, b, _ := c.RGBA()
	targetR, targetG, targetB, _ := target.RGBA()

	r, g, b = r>>8, g>>8, b>>8 // 将 0 到 65535 的值缩放到 0 到 255 的范围
	targetR, targetG, targetB = targetR>>8, targetG>>8, targetB>>8

	return (abs(int(r)-int(targetR)) <= tolerance) &&
		(abs(int(g)-int(targetG)) <= tolerance) &&
		(abs(int(b)-int(targetB)) <= tolerance)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// findCircleCenters 从像素集合中计算圆点的中心，指定一个半径
func findCircleCenters(pixels []image.Point, radius int) []image.Point {
	var centers []image.Point
	visited := make(map[image.Point]bool)

	for _, p := range pixels {
		if visited[p] {
			continue
		}

		// 查找半径内的所有像素点
		var cluster []image.Point
		findClusterInRadius(pixels, p, visited, &cluster, radius)

		if len(cluster) > 0 {
			center := calculateCenter(cluster)
			centers = append(centers, center)
		}
	}

	return centers
}

// findClusterInRadius 找到在指定半径内的所有相关像素
func findClusterInRadius(pixels []image.Point, start image.Point, visited map[image.Point]bool, cluster *[]image.Point, radius int) {
	for _, p := range pixels {
		if !visited[p] && isWithinRadius(start, p, radius) {
			visited[p] = true
			*cluster = append(*cluster, p)
		}
	}
}

// isWithinRadius 检查点是否在指定半径内
func isWithinRadius(center, point image.Point, radius int) bool {
	dx := center.X - point.X
	dy := center.Y - point.Y
	return (dx*dx + dy*dy) <= radius*radius
}

// calculateCenter 计算给定像素集合的中心
func calculateCenter(points []image.Point) image.Point {
	if len(points) == 0 {
		return image.Point{}
	}

	var sumX, sumY int
	for _, p := range points {
		sumX += p.X
		sumY += p.Y
	}
	centerX := sumX / len(points)
	centerY := sumY / len(points)

	return image.Point{X: centerX, Y: centerY}
}

// GetRedPointCenterFromMap 从小地图中获取红点的中心点坐标
// 除了在单人竞技战斗中用，还可以用来判断副本中的boss是否已经阵亡
func GetRedPointCenterFromMap(img *image.RGBA) []image.Point {
	if img == nil {
		return nil
	}
	redPixels := findColorPixels(img, color.RGBA{R: 255, G: 10, B: 10, A: 255}, 50)
	return findCircleCenters(redPixels, 18) // 可以调整radius
}

// GetBluePointCenterFromMap 从小地图中获取蓝点的中心点坐标
// 如果蓝点有多个的话 就根据力度来确认初始位置吧，这个接口不能滥用 稳定性不大，可能只在单人竞技战斗中比较有用
func GetBluePointCenterFromMap(img *image.RGBA) []image.Point {
	if img == nil {
		return nil
	}
	bluePixels := findColorPixels(img, color.RGBA{R: 50, G: 50, B: 255, A: 255}, 50)
	return findCircleCenters(bluePixels, 18) // 可以调整radius
}
