package utils

import (
	"bytes"
	"ddt-copilot/defs"
	"fmt"
	"github.com/lxn/win"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test_GetProcessID(t *testing.T) {
	pid1, _ := GetProcessID(string(defs.ProcessTg3))
	pid2, _ := GetProcessID(string(defs.ProcessTgWeb))
	fmt.Println(pid1)
	fmt.Println(pid2)
}

func Test_GetFirstWindowByPID(t *testing.T) {
	pid, _ := GetProcessID(string(defs.ProcessTgWeb))

	hwnd, err := GetFirstWindowByPID(pid[0])
	if err != nil {
		fmt.Println("获取窗口句柄时出错:", err)
		return
	}

	if hwnd == 0 {
		fmt.Println("未找到窗口句柄")
	} else {
		fmt.Printf("进程ID %d 对应的窗口句柄是: %v\n", pid, hwnd)
	}
}

func GetLastWind() win.HWND {
	wds := GetDDTWindowsHWND()
	return win.HWND(wds[0])
}

func Test_GetAllChildWindows(t *testing.T) {
	lastWd := GetLastWind()
	fmt.Printf("lastWd:%d", lastWd)
}

func Test_Click(t *testing.T) {
	hwnd := 20388674
	win.NotifyWinEvent(win.EVENT_OBJECT_FOCUS, win.HWND(hwnd), win.OBJID_CLIENT, 0)
	return

	lastWd := GetLastWind()
	fmt.Printf("lastWd:%d", lastWd)

	ClickElement(lastWd, defs.ElementFubenEnter, 0)
}

func Test_KeyBoard(t *testing.T) {
	lastWd := GetLastWind()
	fmt.Printf("lastWd:%d", lastWd)

	KeyBoard(lastWd, defs.VK_SPACE, time.Second*3)
}

func Test_CaptureWindow(t *testing.T) {
	lastWd := GetLastWind()
	fmt.Printf("lastWd:%d\n", lastWd)

	//for i := -30; i < 91; i++ {
	//	if i < 0 {
	//		fmt.Printf("%d:ElementImgNum_%d,\n", i, int(math.Abs(float64(i))))
	//	} else {
	//		fmt.Printf("%d:ElementImgNum%d,\n", i, int(math.Abs(float64(i))))
	//	}
	//}
	//return

	//ConfirmDirection(lastWd, defs.DirectionLeft)
	//UpdateAngle(lastWd, 20)
	//Launch(lastWd, 25)
	////Move(lastWd, defs.DirectionRight, 10)
	//return

	// 获取屏幕的截图
	//img, err := CaptureWindowLight(win.HWND(lastWd), defs.GetElementRect(defs.DDTElementRectTypePassBtn))
	//img, err := CaptureWindowLight(win.HWND(lastWd), defs.GetElementRect(defs.DDTElementRectTypeFubenBtn1))
	//img, err := CaptureWindowLight(win.HWND(lastWd), defs.GetElementRect(defs.DDTElementRectTypeFubenBtn2))
	//img, err := CaptureWindowLight(win.HWND(lastWd), defs.GetElementRect(defs.DDTElementRectTypeFubenBtn3))
	//img, err := CaptureWindowLight(win.HWND(lastWd), defs.GetElementRect(defs.DDTElementRectTypeFubenBtn4))
	//img, err := CaptureWindowLight(win.HWND(lastWd), defs.GetElementRect(defs.DDTElementRectTypeFubenBtn5))
	//img, err := CaptureWindowLight(win.HWND(lastWd), defs.GetElementRect(defs.DDTElementRectTypeFubenBtn6))
	//img, err := CaptureWindowLight(win.HWND(lastWd), defs.GetElementRect(defs.DDTElementRectTypeFubenBtn7))
	//img, err := CaptureWindowLight(win.HWND(lastWd), defs.GetElementRect(defs.DDTElementRectTypeFubenBtn8))
	//x, y := 232+110+24+110+20+110+26, 489
	//x, y := 232+110+24+110+20, 489
	//x, y := 32, 558
	//w, h := 35, 15

	//var contnts []string
	//
	//for i := 60; i < 91; i++ {
	//	img, err := CaptureWindowLight(win.HWND(lastWd), &win.RECT{
	//		Left:   int32(x),
	//		Top:    int32(y),
	//		Right:  int32(x + w),
	//		Bottom: int32(y + h),
	//	})
	//
	//	if err != nil {
	//		log.Fatalf("无法截取屏幕: %v", err)
	//	}
	//
	//	// 转成灰色的，以适应不同的背景颜色定位
	//	grayImg := ConvertToGrayWithNormalization(img)
	//
	//	uids := []string{}
	//	for _, pix := range grayImg.Pix {
	//		uids = append(uids, strconv.Itoa(int(pix)))
	//	}
	//
	//	name := fmt.Sprintf("ElementImgNum%d", int(math.Abs(float64(i))))
	//
	//	content := `
	//%s = &image.Gray{
	//	Pix:    []uint8{%s},
	//	Stride: %d,
	//	Rect: image.Rectangle{
	//		Min: image.Point{
	//			X: 0,
	//			Y: 0,
	//		},
	//		Max: image.Point{
	//			X: %d,
	//			Y: %d,
	//		},
	//	},
	//}
	//`
	//	content = fmt.Sprintf(content, name, strings.Join(uids, ","), grayImg.Stride, grayImg.Rect.Max.X, grayImg.Rect.Max.Y)
	//	contnts = append(contnts, content)
	//
	//	// 创建文件并保存为 PNG 格式
	//	file, err := os.Create(fmt.Sprintf("%s.png", name))
	//	if err != nil {
	//		log.Fatalf("无法创建文件: %v", err)
	//	}
	//	defer file.Close()
	//
	//	// 将图像编码为 PNG 格式并写入文件
	//	err = png.Encode(file, grayImg)
	//	if err != nil {
	//		log.Fatalf("保存图片失败: %v", err)
	//	}
	//
	//	fmt.Printf("done %d\n", i)
	//
	//	time.Sleep(time.Second * 1)
	//}
	//
	//for _, contnt := range contnts {
	//	fmt.Printf("%s\n", contnt)
	//}
	//
	//return

	//	isLoop := true
	//	for isLoop {
	//		img, err := CaptureWindowLight(win.HWND(lastWd), defs.GetElementRect(defs.DDTElementRectTypeIsYourTurn), true)
	//		if err != nil {
	//			log.Fatalf("无法截取屏幕: %v", err)
	//		}
	//
	//		// 转成灰色的，以适应不同的背景颜色定位
	//		grayImg := ConvertToGray(img)
	//
	//		uids := []string{}
	//		for _, pix := range grayImg.Pix {
	//			uids = append(uids, strconv.Itoa(int(pix)))
	//		}
	//
	//		name := fmt.Sprintf("ElementImg%s", "DDTElementRectTypeJinjiFightSettle"[18:])
	//
	//		content := `
	//%s = &image.Gray{
	//	Pix:    []uint8{%s},
	//	Stride: %d,
	//	Rect: image.Rectangle{
	//		Min: image.Point{
	//			X: 0,
	//			Y: 0,
	//		},
	//		Max: image.Point{
	//			X: %d,
	//			Y: %d,
	//		},
	//	},
	//}
	//`
	//		content = fmt.Sprintf(content, name, strings.Join(uids, ","), grayImg.Stride, grayImg.Rect.Max.X, grayImg.Rect.Max.Y)
	//
	//		//fmt.Printf("%s\n", content)
	//
	//		diff, _ := CompareGrayImages(grayImg, defs.ElementImgIsYourTurn)
	//		fmt.Printf("diff:%f\n", diff)
	//
	//		// 创建文件并保存为 PNG 格式
	//		file, err := os.Create("fubenleveltmp.png")
	//		if err != nil {
	//			log.Fatalf("无法创建文件: %v", err)
	//		}
	//		defer file.Close()
	//
	//		// 将图像编码为 PNG 格式并写入文件
	//		//err = png.Encode(file, img)
	//		err = png.Encode(file, grayImg)
	//		if err != nil {
	//			log.Fatalf("保存图片失败: %v", err)
	//		}
	//		log.Println("截图成功，已保存为 screenshot.png")
	//	}

	img, err := CaptureWindowLight(win.HWND(lastWd), defs.GetElementRect(defs.DDTElementRectTypeFightReady), true)
	if err != nil {
		log.Fatalf("无法截取屏幕: %v", err)
	}

	// 转成灰色的，以适应不同的背景颜色定位
	grayImg := ConvertToGray(img)

	uids := []string{}
	for _, pix := range grayImg.Pix {
		uids = append(uids, strconv.Itoa(int(pix)))
	}

	name := fmt.Sprintf("ElementImg%s", "DDTElementRectTypeJinjiFightSettle"[18:])

	content := `
%s = &image.Gray{
	Pix:    []uint8{%s},
	Stride: %d,
	Rect: image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: %d,
			Y: %d,
		},
	},
}
`
	content = fmt.Sprintf(content, name, strings.Join(uids, ","), grayImg.Stride, grayImg.Rect.Max.X, grayImg.Rect.Max.Y)

	fmt.Printf("%s\n", content)

	diff, _ := CompareGrayImages(grayImg, defs.ElementImgIsYourTurn)
	fmt.Printf("diff:%f\n", diff)

	// 创建文件并保存为 PNG 格式
	file, err := os.Create("fubenleveltmp.png")
	if err != nil {
		log.Fatalf("无法创建文件: %v", err)
	}
	defer file.Close()

	// 将图像编码为 PNG 格式并写入文件
	//err = png.Encode(file, img)
	err = png.Encode(file, grayImg)
	if err != nil {
		log.Fatalf("保存图片失败: %v", err)
	}
	log.Println("截图成功，已保存为 screenshot.png")
}

func Test_MonitorWindows(t *testing.T) {
	lastWd := GetLastWind()
	fmt.Printf("lastWd:%d\n", lastWd)

	for {
		// 获取屏幕的截图
		img, err := CaptureWindowLight(win.HWND(lastWd), defs.GetElementRect(defs.DDTElementRectTypePassBtn), true)
		if err != nil {
			log.Fatalf("无法截取屏幕: %v", err)
		}

		// 转成灰色的，以适应不同的背景颜色定位
		grayImg := ConvertToGray(img)

		uids := []string{}
		for _, pix := range grayImg.Pix {
			uids = append(uids, strconv.Itoa(int(pix)))
		}

		fmt.Printf("%s\n", strings.Join(uids, ","))

		diff, _ := CompareGrayImages(grayImg, defs.ElementImgPassBtn)
		fmt.Printf("diff:%f\n", diff)
		if IsSimilarity(diff, defs.ImgSimilarityThresholdPassBtn) {
			fmt.Printf("now is your turn\n")
		}

		time.Sleep(time.Second)
	}
}

func Test_PostImgAndOCR(t *testing.T) {
	imagePaths := []string{
		"screenshot.png",
		"pass1.png",
		"pass2.png",
		"Snipaste_2024-11-08_15-44-24.jpg",
		"Snipaste_2024-11-08_16-58-46.jpg",
		"Snipaste_2024-11-08_17-00-55.jpg",
	}

	for _, imagePath := range imagePaths {

		// 读取 PNG 文件
		imageData, err := ioutil.ReadFile(imagePath)
		if err != nil {
			fmt.Println("Error reading image:", err)
			return
		}

		// 创建一个 POST 请求
		resp, err := http.Post("http://localhost:8000", "application/octet-stream", bytes.NewBuffer(imageData))
		if err != nil {
			fmt.Println("Error sending POST request:", err)
			return
		}
		defer resp.Body.Close()

		// 读取响应
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		// 打印响应内容
		fmt.Println("Response from server:", string(body))
	}
}

// 模拟一次脚本执行 Todo finish
func Test_Simulation(t *testing.T) {
	lastWd := GetLastWind()
	fmt.Printf("lastWd:%d\n", lastWd)

	ClickElementAfterLong(lastWd, defs.ElementFubenEnter)

	monitor := func() int {
		for { // 监听屏幕变化
			// Todo zhangzhihui 在识别之前把打开的资料窗口关闭
			img, err := CaptureWindowLight(lastWd, defs.GetElementRect(defs.DDTElementRectTypeFubenSelectText), true)
			if err != nil {
				log.Fatalf("无法截取屏幕: %v", err)
			}

			// 转成灰色的，以适应不同的背景颜色定位
			grayImg := ConvertToGray(img)

			diff, err := CompareGrayImages(grayImg, defs.ElementImgFubenSelectText)
			if err != nil {
				log.Fatalf("CompareGrayImages: %v", err)
			}
			fmt.Printf("diff:%f\n", diff)
			if IsSimilarity(diff, defs.ImgSimilarityThresholdFubenSelectText) {
				fmt.Printf("now is your turn\n")
				return 1
			}

			// 获取屏幕的截图
			img, err = CaptureWindowLight(lastWd, defs.GetElementRect(defs.DDTElementRectTypePassBtn), true)
			if err != nil {
				log.Fatalf("无法截取屏幕: %v", err)
			}

			// 转成灰色的，以适应不同的背景颜色定位
			grayImg = ConvertToGray(img)

			diff, err = CompareGrayImages(grayImg, defs.ElementImgPassBtn)
			if err != nil {
				log.Fatalf("CompareGrayImages: %v", err)
			}
			fmt.Printf("diff:%f\n", diff)
			if IsSimilarity(diff, defs.ImgSimilarityThresholdPassBtn) {
				fmt.Printf("now is your turn\n")
				return 2
			}

			time.Sleep(time.Second)
		}
	}

	var loop func()

	loop = func() {
		ClickElementAfterMid(lastWd, defs.ElementFubenSelect)

		//ctrl := NewCtrlWSMY(lastWd)
		//err := ctrl.SelectFubenMap(defs.FubenIDMY, defs.FubenLevelEasy, true)
		//if err != nil {
		//	fmt.Printf("err:%v\n", err)
		//	return
		//}
		//ClickElement(lastWd, defs.ElementFightEquipItem1)
		//defs.WaitItemSelect()
		//ClickElement(lastWd, defs.ElementFightEquipItem2)
		//defs.WaitItemSelect()
		//ClickElement(lastWd, defs.ElementFightEquipItem3)
		//defs.WaitItemSelect()
		//ClickElement(lastWd, defs.ElementFightSelectItem4)
		//defs.WaitItemSelect()
		//ClickElement(lastWd, defs.ElementFightSelectItem4)
		//defs.WaitItemSelect()
		//ClickElement(lastWd, defs.ElementFightSelectItem4)
		ClickElementAfterMid(lastWd, defs.ElementFightStart)
		ClickElementAfterMid(lastWd, defs.ElementFubenFightStartAck)

		for {
			bret := monitor()
			if bret == 2 {
				Move(lastWd, defs.DirectionLeft, 5)
				time.Sleep(time.Second)
				ConfirmDirection(lastWd, defs.DirectionRight)
				time.Sleep(time.Second)
				num := GetAngle(lastWd)
				needAngle := 30
				diff := needAngle - num
				UpdateAngle(lastWd, diff)
				UseSkill(lastWd, defs.VK_2) // 2
				UseSkill(lastWd, defs.VK_4) // 4
				UseSkill(lastWd, defs.VK_4) // 4
				UseSkill(lastWd, defs.VK_4) // 4
				UseSkill(lastWd, defs.VK_4) // 4
				UseSkill(lastWd, defs.VK_5) // 5
				UseSkill(lastWd, defs.VK_5) // 5
				UseSkill(lastWd, defs.VK_5) // 5
				UseSkill(lastWd, defs.VK_5) // 5
				UseSkill(lastWd, defs.VK_6) // 6
				UseSkill(lastWd, defs.VK_6) // 6
				UseSkill(lastWd, defs.VK_6) // 6
				UseSkill(lastWd, defs.VK_6) // 6
				UseSkill(lastWd, defs.VK_7) // 7
				UseSkill(lastWd, defs.VK_7) // 7
				UseSkill(lastWd, defs.VK_7) // 7
				UseSkill(lastWd, defs.VK_7) // 7
				Launch(lastWd, 75)
				time.Sleep(time.Second * 5) // 5秒不监听，因为这个时候还在出手
			} else if bret == 1 {
				loop()
			}
		}
	}

	loop()
}

//func Test_Simulation(t *testing.T) {
//	lastWd := GetLastWind()
//	fmt.Printf("lastWd:%d\n", lastWd)
//
//	RunScript(lastWd)
//	//for {
//	//	st := GetNowSenseType(lastWd)
//	//	fmt.Printf("now st:%d\n", st)
//	//	time.Sleep(time.Millisecond * 100)
//	//}
//
//}

// loadImage 从指定路径加载 JPG 图片
func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func Test_GetPointCenterFromMap(t *testing.T) {
	// 从本地读取 JPG 图片
	img, err := loadImage("Snipaste_2024-11-13_11-19-26.jpg") // 替换为实际图片路径
	if err != nil {
		panic(err.Error())
	}

	// 将 image.Image 转换为 *image.RGBA
	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		// 如果不是 *image.RGBA 类型，则进行转换
		rgbaImg = image.NewRGBA(img.Bounds())
		for y := 0; y < img.Bounds().Dy(); y++ {
			for x := 0; x < img.Bounds().Dx(); x++ {
				rgbaImg.Set(x, y, img.At(x, y))
			}
		}
	}

	// 打印结果
	println("红色圆点中心:")
	for _, center := range GetRedPointCenterFromMap(rgbaImg) {
		println(center.X, center.Y)
	}

	println("蓝色圆点中心:")
	for _, center := range GetBluePointCenterFromMap(rgbaImg) {
		println(center.X, center.Y)
	}
}

// Todo 写一个测试用例函数验证一下所有数值的diff是多少 确定合适的阈值
