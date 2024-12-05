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
#include <defs/struct.h>

std::string format(const std::string& fmt);

template<typename T, typename... Args>
std::string format(const std::string& fmt, const T& value, const Args&... args) {
    std::ostringstream oss;
    size_t pos = fmt.find("{}");

    if (pos != std::string::npos) {
        oss << fmt.substr(0, pos) << value; // 插入第一个参数
        oss << format(fmt.substr(pos + 2), args...); // 递归处理剩余参数
    } else {
        return fmt; // 如果没有占位符，返回原始字符串
    }

    return oss.str();
}

std::wstring StringToWString(const std::string& str);

std::string WStringToString(const std::wstring& wstr);

std::vector<DWORD> GetProcessID(const std::wstring& processName);

std::vector<HWND> GetHwndsByPID(DWORD pid);

std::vector<HWND> GetChildHwnds(HWND parentHwnd);

HBITMAP CaptureWindow(HWND hwnd, const RECT* captureRect = nullptr);

void PrintBitmapInfo(HBITMAP hBitmap);

void ConvertBitmapColorWithNormalization(HBITMAP hBitmap);

bool SaveBitmapToBmp(HBITMAP hBitmap, const std::wstring& filename);

HBITMAP LoadBitmapFromBmp(const std::wstring& filename);

bool CaptureWindowAndSave(HWND hwnd, const std::wstring& filename, const RECT* captureRect = nullptr);

bool IslegalPoint(StructDefs::Point point);

StructDefs::Point FindPic(HBITMAP hBitmapA, HBITMAP hBitmapB, float similarityThreshold, int maxColorDiff = 30, const RECT* importantRect = nullptr);

void ClickLeft(HWND hwnd, int32_t x, int32_t y, std::chrono::milliseconds duration = std::chrono::milliseconds(0));

void ClickRight(HWND hwnd, int32_t x, int32_t y, std::chrono::milliseconds duration = std::chrono::milliseconds(0));

void KeyBoard(HWND hwnd, uintptr_t key, std::chrono::milliseconds duration = std::chrono::milliseconds(0));

#endif // UTILS_H
