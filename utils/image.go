package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"ddt-copilot/defs"
	"errors"
	"fmt"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"unsafe"
)

// 图像处理相关接口

var (
	user32          = windows.NewLazySystemDLL("user32.dll")
	procPrintWindow = user32.NewProc("PrintWindow")
)

// CaptureWindow 对整个窗口截图，然后裁剪图像
func CaptureWindow(hwnd win.HWND, captureRect *win.RECT) (*image.RGBA, error) {
	hWindowDC := win.GetDC(hwnd)
	defer win.ReleaseDC(hwnd, hWindowDC)
	hCaptureDC := win.CreateCompatibleDC(hWindowDC)
	defer win.DeleteDC(hCaptureDC)

	var rc win.RECT
	_ = win.GetClientRect(hwnd, &rc)

	// 使用窗口的宽度和高度
	width := rc.Right - rc.Left
	height := rc.Bottom - rc.Top

	hBitmap := win.CreateCompatibleBitmap(hWindowDC, width, height)
	defer win.DeleteObject(win.HGDIOBJ(hBitmap))

	win.SelectObject(hCaptureDC, win.HGDIOBJ(hBitmap))

	// 捕获窗口的内容
	procPrintWindow.Call(uintptr(hwnd), uintptr(hCaptureDC), 0)

	var header win.BITMAPINFOHEADER
	header.BiSize = uint32(unsafe.Sizeof(header))
	header.BiPlanes = 1
	header.BiBitCount = 32
	header.BiWidth = width
	header.BiHeight = -height // 负值表示从顶部开始填充
	header.BiCompression = win.BI_RGB
	header.BiSizeImage = 0

	bitmapDataSize := uintptr(((int64(width)*int64(header.BiBitCount) + 31) / 32) * 4 * int64(height))
	hmem := win.GlobalAlloc(win.GMEM_MOVEABLE, bitmapDataSize)
	defer win.GlobalFree(hmem)
	memptr := win.GlobalLock(hmem)
	defer win.GlobalUnlock(hmem)

	if win.GetDIBits(hWindowDC, hBitmap, 0, uint32(height), (*uint8)(memptr), (*win.BITMAPINFO)(unsafe.Pointer(&header)), win.DIB_RGB_COLORS) == 0 {
		return nil, errors.New("GetDIBits 失败")
	}

	// 创建 RGBA 图像
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	i := 0
	src := uintptr(memptr)
	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			v0 := *(*uint8)(unsafe.Pointer(src))
			v1 := *(*uint8)(unsafe.Pointer(src + 1))
			v2 := *(*uint8)(unsafe.Pointer(src + 2))

			// BGRA => RGBA, and set A to 255
			img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3] = v2, v1, v0, 255

			i += 4
			src += 4
		}
	}

	// 如果 captureRect 不是 nil，裁剪图像
	if captureRect != nil {
		croppedImg := image.NewRGBA(image.Rect(0, 0, int(captureRect.Right-captureRect.Left), int(captureRect.Bottom-captureRect.Top)))
		for y := captureRect.Top; y < captureRect.Bottom; y++ {
			for x := captureRect.Left; x < captureRect.Right; x++ {
				croppedImg.Set(int(x-captureRect.Left), int(y-captureRect.Top), img.At(int(x), int(y)))
			}
		}
		img = croppedImg // 使用裁剪后的图像
	}

	return img, nil
}

const (
	encryptionKeyImg = "b058c6e56c148702"
	imgEncryptSuffix = "zz"
)

// encrypt encrypts the data using AES.
func encrypt(data []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	// Create a new GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// decrypt decrypts the data using AES.
func decrypt(data []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("密文太短")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	decrypted, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

// SaveImageToPng saves an image.RGBA as a PNG file with both encrypted and unencrypted versions.
func SaveImageToPng(img image.Image, filename string) error {
	// Generate the unencrypted filename
	unencryptedFilename := fmt.Sprintf("%s_unencrypted.png", filename)

	// Create and save the unencrypted PNG file
	unencryptedFile, err := os.Create(unencryptedFilename)
	if err != nil {
		return fmt.Errorf("无法创建未加密文件: %v", err)
	}
	defer unencryptedFile.Close()

	// Encode the image to the unencrypted PNG file
	if err := png.Encode(unencryptedFile, img); err != nil {
		return fmt.Errorf("保存未加密图片失败: %v", err)
	}

	// Prepare to create the encrypted file
	encryptedFilename := fmt.Sprintf("%s.%s", filename, imgEncryptSuffix)

	// Read the unencrypted PNG file to encrypt
	unencryptedFile.Seek(0, 0) // Reset file pointer to the beginning
	imgData, err := io.ReadAll(unencryptedFile)
	if err != nil {
		return fmt.Errorf("读取未加密文件失败: %v", err)
	}

	// Encrypt the image data
	encryptedData, err := encrypt(imgData, encryptionKeyImg)
	if err != nil {
		return fmt.Errorf("加密失败: %v", err)
	}

	// Write the encrypted data to the encrypted PNG file
	if err = os.WriteFile(encryptedFilename, encryptedData, 0644); err != nil {
		return fmt.Errorf("写入加密文件失败: %v", err)
	}

	return nil
}

// LoadPngToImage loads a PNG file and decrypts it to return an image.RGBA.
func LoadPngToImage(filename string) (image.Image, error) {
	filename = fmt.Sprintf("%s.%s", filename, imgEncryptSuffix)

	// Read the encrypted PNG file
	encryptedData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// Decrypt the image data
	decryptedData, err := decrypt(encryptedData, encryptionKeyImg)
	if err != nil {
		return nil, fmt.Errorf("解密失败: %v", err)
	}

	// Decode the decrypted PNG data into an image
	img, err := png.Decode(bytes.NewReader(decryptedData))
	if err != nil {
		return nil, fmt.Errorf("解码PNG失败: %v", err)
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

func CompareGrayImages(img1, img2 *image.Gray) (float64, error) {
	if img1.Bounds() != img2.Bounds() {
		return 0, errors.New("images must have the same dimensions")
	}

	var matchCount float64
	width, height := img1.Bounds().Max.X, img1.Bounds().Max.Y
	pixelCount := float64(width * height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gray1 := img1.GrayAt(x, y).Y
			gray2 := img2.GrayAt(x, y).Y
			if int(math.Abs(float64(gray1-gray2))) < defs.Colorthreshold {
				matchCount++
			}
		}
	}

	match := matchCount / pixelCount
	return match, nil
}

func IsImageSimilarity(img1, img2 *image.Gray, threshold float64) bool {
	match, err := CompareGrayImages(img1, img2)
	if err != nil {
		return false
	}
	return match >= threshold
}
