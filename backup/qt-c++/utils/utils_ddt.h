#ifndef UTILS_DDT_H
#define UTILS_DDT_H

#include <defs/constants.h>
#include <utils/utils.h>

std::vector<HWND> GetDDTHwnds();

void ClickEmpty(HWND hwnd);

RECT ToWinRECT(StructDefs::Rect rect);

StructDefs::Rect ToRECT(RECT rect);

void FocusDDTHwnd(HWND hwnd, bool dropBlock = true);

HBITMAP CaptureWindowLight(HWND hwnd, const RECT* captureRect = nullptr, bool dropBlock = true);

void ClickPointSleep(HWND hwnd,  const StructDefs::Point point, std::chrono::milliseconds duration = std::chrono::milliseconds(0));
void ClickPointSleepS(HWND hwnd,  const StructDefs::Point point, std::chrono::milliseconds duration = std::chrono::milliseconds(Constants::ClickWaitShort));
void ClickPointSleepM(HWND hwnd,  const StructDefs::Point point, std::chrono::milliseconds duration = std::chrono::milliseconds(Constants::ClickWaitMid));
void ClickPointSleepL(HWND hwnd,  const StructDefs::Point point, std::chrono::milliseconds duration = std::chrono::milliseconds(Constants::ClickWaitLong));

StructDefs::Point NewPointFromRectCenter(const StructDefs::Rect rect);

void ClickRectSleep(HWND hwnd,  const StructDefs::Rect rect, std::chrono::milliseconds duration = std::chrono::milliseconds(0));
void ClickRectSleepS(HWND hwnd,  const StructDefs::Rect rect, std::chrono::milliseconds duration = std::chrono::milliseconds(Constants::ClickWaitShort));
void ClickRectSleepM(HWND hwnd,  const StructDefs::Rect rect, std::chrono::milliseconds duration = std::chrono::milliseconds(Constants::ClickWaitMid));
void ClickRectSleepL(HWND hwnd,  const StructDefs::Rect rect, std::chrono::milliseconds duration = std::chrono::milliseconds(Constants::ClickWaitLong));

void UpdateAngle(HWND hwnd, int angle);
void Launch(HWND hwnd, int power);
#endif // UTILS_DDT_H
