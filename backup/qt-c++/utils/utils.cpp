#ifndef UTILS_H
#define UTILS_H

#include <windows.h>
#include <tlhelp32.h>
#include <string>
#include <vector>
#include <iostream>
#include <stdexcept>
#include <mutex>
#include <fstream>
#include <sstream>
#include <iomanip>
#include <locale>
#include <codecvt>
#include <cmath>
#include <thread>
#include <defs/constants.h>
#include <defs/struct.h>

std::string format(const std::string& fmt) {
    return fmt; // 处理没有参数的情况
}

// 将窄字符字符串转换为宽字符字符串
std::wstring StringToWString(const std::string& str) {
    std::wstring_convert<std::codecvt_utf8<wchar_t>> converter;
    return converter.from_bytes(str);
}


// 将宽字符字符串转换为窄字符字符串
std::string WStringToString(const std::wstring& wstr) {
    int size_needed = WideCharToMultiByte(CP_UTF8, 0, &wstr[0], (int)wstr.size(), NULL, 0, NULL, NULL);
    std::string str(size_needed, 0);
    WideCharToMultiByte(CP_UTF8, 0, &wstr[0], (int)wstr.size(), &str[0], size_needed, NULL, NULL);
    return str;
}

// 根据进程名获取进程id
std::vector<DWORD> GetProcessID(const std::wstring& processName) {
    std::vector<DWORD> processIDs;

    // 创建进程快照
    HANDLE processSnap = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0);
    if (processSnap == INVALID_HANDLE_VALUE) {
        throw std::runtime_error("Failed to create process snapshot");
    }

    PROCESSENTRY32 processEntry;
    processEntry.dwSize = sizeof(processEntry);

    // 遍历进程
    if (Process32First(processSnap, &processEntry)) {
        do {
            // 将进程名转换为 std::wstring
            std::wstring exeFileName(processEntry.szExeFile);
            if (_wcsicmp(exeFileName.c_str(), processName.c_str()) == 0) {
                processIDs.push_back(processEntry.th32ProcessID);
            }
        } while (Process32Next(processSnap, &processEntry));
    }

    CloseHandle(processSnap);

    if (processIDs.empty()) {
        throw std::runtime_error("Process not found");
    }

    return processIDs;
}

// 根据进程id获取包含的窗口句柄
std::vector<HWND> GetHwndsByPID(DWORD pid) {
    static std::vector<HWND> windowHandles;
    static std::mutex mtx; // 互斥锁

    std::lock_guard<std::mutex> lock(mtx); // 确保线程安全，对象析构时自动解锁
    windowHandles.clear();

    auto callback = [](HWND hwnd, LPARAM lParam) -> BOOL {
        DWORD curPid;
        GetWindowThreadProcessId(hwnd, &curPid);

        if (curPid == lParam) {
            windowHandles.push_back(hwnd);
        }
        return TRUE; // 继续枚举
    };

    EnumWindows(callback, static_cast<LPARAM>(pid));

    return windowHandles;
}

std::vector<HWND> GetChildHwnds(HWND parentHwnd) {
    static std::vector<HWND> windowHandles;
    static std::mutex mtx; // 互斥锁

    std::lock_guard<std::mutex> lock(mtx); // 确保线程安全，对象析构时自动解锁
    windowHandles.clear();

    if (!IsWindow(parentHwnd)){
        return windowHandles;
    }

    auto callback = [](HWND hwnd, LPARAM) -> BOOL {
        windowHandles.push_back(hwnd);
        return TRUE; // 继续枚举
    };

    EnumChildWindows(parentHwnd, callback, 0);

    return windowHandles;
}

HBITMAP CaptureWindow(HWND hwnd, const RECT* captureRect = nullptr) {
    HDC hWindowDC = GetWindowDC(hwnd);
    HDC hCaptureDC = CreateCompatibleDC(hWindowDC);

    RECT rc;
    GetClientRect(hwnd, &rc);

    int width = rc.right - rc.left;
    int height = rc.bottom - rc.top;

    HBITMAP hBitmap = CreateCompatibleBitmap(hWindowDC, width, height);
    SelectObject(hCaptureDC, hBitmap);

    // 捕捉窗口内容，不使用BitBlt，因为可能导致后台截图失败。但是PrintWindow不能指定截取的区域，所以要裁剪的话先截取整个屏幕，再裁剪成所需的区域
    PrintWindow(hwnd, hCaptureDC, 0);

    HBITMAP hCroppedBitmap = nullptr;
    if (captureRect) { // 如果需要裁剪，创建一个新的位图
        int cropWidth = captureRect->right - captureRect->left;
        int cropHeight = captureRect->bottom - captureRect->top;
        hCroppedBitmap = CreateCompatibleBitmap(hWindowDC, cropWidth, cropHeight);

        if (hCroppedBitmap) {
            HDC hCroppedDC = CreateCompatibleDC(hWindowDC);
            SelectObject(hCroppedDC, hCroppedBitmap);
            BitBlt(hCroppedDC, 0, 0, cropWidth, cropHeight, hCaptureDC, captureRect->left, captureRect->top, SRCCOPY); // 使用 BitBlt 从 hCaptureDC 复制指定区域到裁剪位图
            DeleteDC(hCroppedDC); // 清理
        }
    }

    DeleteDC(hCaptureDC);
    ReleaseDC(hwnd, hWindowDC);

    HBITMAP result;

    if (hCroppedBitmap) { // 未返回位图的要释放资源
        result = hCroppedBitmap;
        DeleteObject(hBitmap);
    } else {
        result = hBitmap;
        DeleteObject(hCroppedBitmap);
    }

    return result;
}

void PrintBitmapInfo(HBITMAP hBitmap) {
    BITMAP bmp;
    if (GetObject(hBitmap, sizeof(BITMAP), &bmp) == 0) {
        std::cerr << "Failed to get bitmap information!" << std::endl;
        return;
    }

    std::cout << "-----------Bitmap Information begin-----------" << std::endl;
    std::cout << "Type: " << bmp.bmType << std::endl;
    std::cout << "Width: " << bmp.bmWidth << " pixels" << std::endl;
    std::cout << "Height: " << bmp.bmHeight << " pixels" << std::endl;
    std::cout << "Planes: " << bmp.bmPlanes << std::endl;
    std::cout << "Bits per Pixel: " << bmp.bmBitsPixel << " bits" << std::endl;
    std::cout << "Size: " << bmp.bmWidthBytes << " bytes per scan line" << std::endl;
    std::cout << "Total Size: " << bmp.bmHeight * bmp.bmWidthBytes << " bytes" << std::endl;
    std::cout << "-----------Bitmap Information end-------------" << std::endl;
}

// 统一处理各像素点的颜色值
void ConvertBitmapColorWithNormalization(HBITMAP hBitmap) {
    if (!hBitmap) {
        std::cerr << "Invalid HBITMAP!" << std::endl;
        return;
    }

    BITMAP bmp;
    GetObject(hBitmap, sizeof(BITMAP), &bmp);

    // 创建位图信息头
    BITMAPINFO bi;
    ZeroMemory(&bi, sizeof(BITMAPINFO));
    bi.bmiHeader.biSize = sizeof(BITMAPINFOHEADER);
    bi.bmiHeader.biWidth = bmp.bmWidth;
    bi.bmiHeader.biHeight = bmp.bmHeight;
    bi.bmiHeader.biPlanes = bmp.bmPlanes;
    bi.bmiHeader.biBitCount = bmp.bmBitsPixel;
    bi.bmiHeader.biCompression = BI_RGB;

    // 获取位图数据，Todo：这里可以封装成一个函数，传入hBitmap返回pPixels，位图信息头都是可以从hBitmap中获取的
    BYTE* pPixels = new BYTE[bmp.bmWidthBytes * bmp.bmHeight];
    GetDIBits(GetDC(NULL), hBitmap, 0, bmp.bmHeight, pPixels, &bi, DIB_RGB_COLORS);

    int bytesPerPixel = bmp.bmBitsPixel / 8; // 根据位深度计算每个像素的字节数
    // 修改每个像素的颜色值，
    for (int y = 0; y < bmp.bmHeight; y++) {
        for (int x = 0; x < bmp.bmWidth; x++)   {
            int index = (y * bmp.bmWidthBytes) + (x * bytesPerPixel); // 每个像素3个字节

            BYTE blue = pPixels[index];         // 蓝色分量
            BYTE green = pPixels[index + 1];    // 绿色分量
            BYTE red = pPixels[index + 2];      // 红色分量

            // 计算灰度值
            BYTE grayValue = static_cast<BYTE>((red * 299 + green * 587 + blue * 114) / 1000);

            // 根据灰度值设置颜色
            BYTE color = (grayValue == 0) ? 0 : 128; // 如果灰度值为0，则为黑色，否则为128【中灰色】

            pPixels[index] = color;        // 蓝色分量
            pPixels[index + 1] = color;    // 绿色分量
            pPixels[index + 2] = color;    // 红色分量
        }
    }

    // 将修改后的数据设置回位图
    SetDIBits(GetDC(NULL), hBitmap, 0, bmp.bmHeight, pPixels, &bi, DIB_RGB_COLORS);

    // 清理
    delete[] pPixels;
}

bool SaveBitmapToBmp(HBITMAP hBitmap, const std::wstring& filename) {
    if (!hBitmap) {
        std::cerr << "Invalid HBITMAP!" << std::endl;
        return false;
    }

    BITMAP bmp;
    if (GetObject(hBitmap, sizeof(BITMAP), &bmp) == 0) {
        std::cerr << "Failed to get bitmap information!" << std::endl;
        return false;
    }

    // 创建位图文件头
    BITMAPFILEHEADER bmfHeader;
    bmfHeader.bfType = 0x4D42; // "BM"
    bmfHeader.bfSize = sizeof(BITMAPFILEHEADER) + sizeof(BITMAPINFOHEADER) + bmp.bmWidthBytes * bmp.bmHeight;
    bmfHeader.bfReserved1 = 0;
    bmfHeader.bfReserved2 = 0;
    bmfHeader.bfOffBits = sizeof(BITMAPFILEHEADER) + sizeof(BITMAPINFOHEADER);

    // 创建位图信息头
    BITMAPINFOHEADER bi;
    bi.biSize = sizeof(BITMAPINFOHEADER);
    bi.biWidth = bmp.bmWidth;
    bi.biHeight = bmp.bmHeight;
    bi.biPlanes = 1;
    bi.biBitCount = bmp.bmBitsPixel; // 使用实际的位深度
    bi.biCompression = BI_RGB; // 不压缩
    bi.biSizeImage = 0; // 自动计算
    bi.biXPelsPerMeter = 0;
    bi.biYPelsPerMeter = 0;
    bi.biClrUsed = 0;
    bi.biClrImportant = 0;

    // 打开文件以写入
    HANDLE hFile = CreateFileW(filename.c_str(), GENERIC_WRITE, 0, NULL, CREATE_ALWAYS, FILE_ATTRIBUTE_NORMAL, NULL);
    if (hFile == INVALID_HANDLE_VALUE) {
        std::cerr << "Could not create file!" << std::endl;
        return false;
    }

    // 写入位图文件头和位图信息头
    DWORD written;
    WriteFile(hFile, &bmfHeader, sizeof(BITMAPFILEHEADER), &written, NULL);
    WriteFile(hFile, &bi, sizeof(BITMAPINFOHEADER), &written, NULL);

    // 根据位图的位深度处理数据
    int bitsPerPixel = bmp.bmBitsPixel;
    BYTE* pPixels = nullptr;

    if (bitsPerPixel == 24 || bitsPerPixel == 32) {
        // 为 24 位或 32 位位图分配内存
        pPixels = new BYTE[bmp.bmWidthBytes * bmp.bmHeight];
        GetDIBits(GetDC(NULL), hBitmap, 0, bmp.bmHeight, pPixels, (BITMAPINFO*)&bi, DIB_RGB_COLORS);

        // 处理行对齐
        for (int i = 0; i < bmp.bmHeight; ++i) {
            WriteFile(hFile, pPixels + (i * bmp.bmWidthBytes), bmp.bmWidthBytes, &written, NULL);
        }
    } else if (bitsPerPixel == 8) {
        // 8 位位图
        pPixels = new BYTE[bmp.bmWidthBytes * bmp.bmHeight];
        GetDIBits(GetDC(NULL), hBitmap, 0, bmp.bmHeight, pPixels, (BITMAPINFO*)&bi, DIB_RGB_COLORS);

        // 写入调色板（256 色）
        RGBQUAD palette[256];
        for (int i = 0; i < 256; i++) {
            palette[i].rgbRed = i;   // 示例调色板
            palette[i].rgbGreen = i;
            palette[i].rgbBlue = i;
            palette[i].rgbReserved = 0;
        }
        WriteFile(hFile, palette, sizeof(palette), &written, NULL);

        // 写入像素数据
        WriteFile(hFile, pPixels, bmp.bmWidthBytes * bmp.bmHeight, &written, NULL);
    } else {
        std::cerr << "Unsupported bit count: " << bitsPerPixel << std::endl;
        CloseHandle(hFile);
        return false;
    }

    // 释放资源
    delete[] pPixels;
    CloseHandle(hFile);
    return true;
}

// Todo：已生成的bmp文件需要加密和解密
HBITMAP LoadBitmapFromBmp(const std::wstring& filename) {
    HANDLE hFile = CreateFileW(filename.c_str(), GENERIC_READ, 0, NULL, OPEN_EXISTING, FILE_ATTRIBUTE_NORMAL, NULL);
    if (hFile == INVALID_HANDLE_VALUE) {
        std::cerr << "Failed to open file! filename: " << WStringToString(filename) << std::endl;
        return NULL;
    }

    // 读取位图文件头
    BITMAPFILEHEADER bmfHeader;
    DWORD bytesRead;
    ReadFile(hFile, &bmfHeader, sizeof(BITMAPFILEHEADER), &bytesRead, NULL);
    if (bmfHeader.bfType != 0x4D42) { // 检查文件类型
        CloseHandle(hFile);
        std::cerr << "Not a valid BMP file!" << std::endl;
        return NULL;
    }

    // 读取位图信息头
    BITMAPINFOHEADER bi;
    ReadFile(hFile, &bi, sizeof(BITMAPINFOHEADER), &bytesRead, NULL);

    // 分配内存用于位图数据
    BYTE* pPixels = new BYTE[bmfHeader.bfSize - sizeof(BITMAPFILEHEADER) - sizeof(BITMAPINFOHEADER)];
    ReadFile(hFile, pPixels, bmfHeader.bfSize - sizeof(BITMAPFILEHEADER) - sizeof(BITMAPINFOHEADER), &bytesRead, NULL);

    // 创建 HBITMAP
    HBITMAP hBitmap = CreateDIBitmap(GetDC(NULL), &bi, CBM_INIT, pPixels, (BITMAPINFO*)&bi, DIB_RGB_COLORS);

    // 释放资源
    delete[] pPixels;
    CloseHandle(hFile);

    return hBitmap;
}

bool CaptureWindowAndSave(HWND hwnd, const std::wstring& filename, const RECT* captureRect = nullptr) {
    HBITMAP hBitmap = CaptureWindow(hwnd, captureRect);
    auto ret = SaveBitmapToBmp(hBitmap, filename);
    DeleteObject(hBitmap); // 释放hBitmap
    return ret;
}

bool IslegalPoint(StructDefs::Point point) {
    return point.x != Constants::PointIllegal.x && point.y != Constants::PointIllegal.y;
}

// 时间复杂度很高，hBitmapA尽量不要是整个窗口的区域，而是先确定一个大概的范围，缩小hBitmapA，在特定区域去找hBitmapB Todo：importantRect参数貌似会导致dataA不对，以后解决，先不用这个参数
StructDefs::Point FindPic(HBITMAP hBitmapA, HBITMAP hBitmapB, float similarityThreshold, int maxColorDiff = 30, const RECT* importantRect = nullptr) {
    BITMAP bmpA, bmpB;
    GetObject(hBitmapA, sizeof(BITMAP), &bmpA);
    GetObject(hBitmapB, sizeof(BITMAP), &bmpB);

    // 确保位图类型相同
    if (bmpA.bmBitsPixel != bmpB.bmBitsPixel) {
        return Constants::PointIllegal; // 不同类型的位图，不支持查找
    }

    int widthA = bmpA.bmWidth;
    int heightA = bmpA.bmHeight;
    int widthB = bmpB.bmWidth;
    int heightB = bmpB.bmHeight;

    // 确保B的尺寸小于A
    if (widthB > widthA || heightB > heightA) {
        return Constants::PointIllegal;
    }

    int bytesPerPixel = bmpA.bmBitsPixel / 8; // 根据位深度计算每个像素的字节数

    // 创建位图信息头
    BITMAPINFO biA, biB;
    ZeroMemory(&biA, sizeof(BITMAPINFO));
    ZeroMemory(&biB, sizeof(BITMAPINFO));

    biA.bmiHeader.biSize = sizeof(BITMAPINFOHEADER);
    biA.bmiHeader.biWidth = bmpA.bmWidth;
    biA.bmiHeader.biHeight = bmpA.bmHeight;
    biA.bmiHeader.biPlanes = bmpA.bmPlanes;
    biA.bmiHeader.biBitCount = bmpA.bmBitsPixel;
    biA.bmiHeader.biCompression = BI_RGB;

    biB.bmiHeader.biSize = sizeof(BITMAPINFOHEADER);
    biB.bmiHeader.biWidth = bmpB.bmWidth;
    biB.bmiHeader.biHeight = bmpB.bmHeight;
    biB.bmiHeader.biPlanes = bmpB.bmPlanes;
    biB.bmiHeader.biBitCount = bmpB.bmBitsPixel;
    biB.bmiHeader.biCompression = BI_RGB;

    BYTE* dataA = new BYTE[bmpA.bmWidthBytes * bmpA.bmHeight];
    GetDIBits(GetDC(NULL), hBitmapA, 0, bmpA.bmHeight, dataA, &biA, DIB_RGB_COLORS);

    BYTE* dataB = new BYTE[bmpB.bmWidthBytes * bmpB.bmHeight];
    GetDIBits(GetDC(NULL), hBitmapB, 0, bmpB.bmHeight, dataB, &biA, DIB_RGB_COLORS);

    // 确定遍历范围
    int startX = 0, startY = 0;
    int endX = widthA - widthB;
    int endY = heightA - heightB;

    if (importantRect) {
        startX = importantRect->left;
        startY = importantRect->top;
        endX = std::min(endX, static_cast<int>(importantRect->right) - widthB);
        endY = std::min(endY, static_cast<int>(importantRect->bottom) - heightB);
    }

    // 遍历A图像的指定区域
    for (int y = startY; y <= endY; ++y) {
        for (int x = startX; x <= endX; ++x) {
            int matchCount = 0;
            int totalPixels = widthB * heightB;

            // 比较B图像在A中的对应区域
            for (int j = 0; j < heightB; ++j) {
                for (int i = 0; i < widthB; ++i) {
                    int indexA = ((y + j) * widthA + (x + i)) * bytesPerPixel;
                    int indexB = (j * widthB + i) * bytesPerPixel;

                    BYTE rA = dataA[indexA + 2];
                    BYTE gA = dataA[indexA + 1];
                    BYTE bA = dataA[indexA + 0];

                    BYTE rB = dataB[indexB + 2];
                    BYTE gB = dataB[indexB + 1];
                    BYTE bB = dataB[indexB + 0];

                    // 计算颜色差异
                    int diff = abs(rA - rB) + abs(gA - gB) + abs(bA - bB);
                    if (diff < maxColorDiff) { // 颜色相似度的阈值
                        matchCount++;
                    }
                }
            }

            // 计算相似度
            float similarity = (float)matchCount / totalPixels;
            if (similarity >= similarityThreshold) {
                return { x + widthB / 2, y + heightB / 2 }; // 返回中心坐标
            }
        }
    }

    delete[] dataA;
    delete[] dataB;

    return Constants::PointIllegal; // 未找到
}

// 左键点击
void ClickLeft(HWND hwnd, int32_t x, int32_t y, std::chrono::milliseconds duration = std::chrono::milliseconds(0)) {
    LPARAM lParam = (y << 16) | (x & 0xFFFF);
    PostMessage(hwnd, WM_LBUTTONDOWN, MK_LBUTTON, lParam);
    std::this_thread::sleep_for(duration); // 线程睡眠(毫秒)
    PostMessage(hwnd, WM_LBUTTONUP, 0, lParam);
}

// 右键点击
void ClickRight(HWND hwnd, int32_t x, int32_t y, std::chrono::milliseconds duration = std::chrono::milliseconds(0)) {
    LPARAM lParam = (y << 16) | (x & 0xFFFF);
    PostMessage(hwnd, WM_RBUTTONDOWN, MK_RBUTTON, lParam);
    std::this_thread::sleep_for(duration); // 线程睡眠(毫秒)
    PostMessage(hwnd, WM_RBUTTONUP, 0, lParam);
}

// 模拟按键
void KeyBoard(HWND hwnd, uintptr_t key, std::chrono::milliseconds duration = std::chrono::milliseconds(0)) {
    PostMessage(hwnd, WM_KEYDOWN, key, 0);
    std::this_thread::sleep_for(duration); // 线程睡眠(毫秒)
    PostMessage(hwnd, WM_KEYUP, key, 0);
}
#endif // UTILS_H
