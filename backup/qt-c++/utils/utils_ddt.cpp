#ifndef UTILS_DDT_H
#define UTILS_DDT_H

#include <defs/constants.h>
#include <utils/utils.h>

// 点击空白处，激活窗口
void ClickEmpty(HWND hwnd) {
    ClickLeft(hwnd, Constants::PointNoEffect.x,  Constants::PointNoEffect.y);
}

std::vector<HWND> GetDDTHwnds() {
    std::vector<HWND> hwnds;
    try {
        std::vector<DWORD> pids = GetProcessID(Constants::ProcessNameTGWeb);
        for (DWORD pid : pids) {
            std::vector<HWND> parentHwnds = GetHwndsByPID(pid);
            if (parentHwnds.size() < 1){
                continue;
            }
            std::vector<HWND> childs = GetChildHwnds(parentHwnds[0]); // 有效的窗口为第1个
            if (childs.size() < 5){
                continue;
            }
            hwnds.push_back(childs[4]); // 子窗口为第5个
        }
    } catch (const std::exception& e) {
        std::cerr << e.what() << std::endl;
    }
    return hwnds;
}

RECT ToWinRECT(StructDefs::Rect rect) {
    return {rect.x, rect.y, rect.x+rect.w, rect.y+rect.h};
}

StructDefs::Rect ToRECT(RECT rect) {
    return {rect.left, rect.top, rect.right-rect.left, rect.bottom-rect.top};
}

void FocusDDTHwnd(HWND hwnd, bool dropBlock = true){
    //SetActiveWindow(hwnd);
    ClickEmpty(hwnd); // 激活界面
    if (dropBlock) {  // 取出遮挡，只取主界面，最多5层遮挡
        for (int i = 0; i < 5; ++i) {
            KeyBoard(hwnd, Constants::VKCode_ESCAPE);
        }
    }
    std::this_thread::sleep_for(std::chrono::milliseconds(100));
}

HBITMAP CaptureWindowLight(HWND hwnd, const RECT* captureRect = nullptr, bool dropBlock = true) {
    FocusDDTHwnd(hwnd, dropBlock);
    return CaptureWindow(hwnd, captureRect);
}

void ClickPointSleep(HWND hwnd,  const StructDefs::Point point, std::chrono::milliseconds duration = std::chrono::milliseconds(0)) {
    ClickLeft(hwnd, point.x, point.y);
    std::this_thread::sleep_for(duration);
}

// 点击后等待短时间-0.15秒
void ClickPointSleepS(HWND hwnd,  const StructDefs::Point point, std::chrono::milliseconds duration = std::chrono::milliseconds(Constants::ClickWaitShort)) {
    ClickPointSleep(hwnd, point, duration);
}

// 点击后等待中时间-0.7秒
void ClickPointSleepM(HWND hwnd,  const StructDefs::Point point, std::chrono::milliseconds duration = std::chrono::milliseconds(Constants::ClickWaitMid)) {
    ClickPointSleep(hwnd, point, duration);
}

// 点击后等待长时间-2秒
void ClickPointSleepL(HWND hwnd,  const StructDefs::Point point, std::chrono::milliseconds duration = std::chrono::milliseconds(Constants::ClickWaitLong)) {
    ClickPointSleep(hwnd, point, duration);
}

// 返回rect的中心点坐标
StructDefs::Point NewPointFromRectCenter(const StructDefs::Rect rect) {
    StructDefs::Point point{rect.x+rect.w/2, rect.y+rect.h/2};
    return point;
}

void ClickRectSleep(HWND hwnd,  const StructDefs::Rect rect, std::chrono::milliseconds duration = std::chrono::milliseconds(0)) {
    ClickPointSleep(hwnd, NewPointFromRectCenter(rect), duration);
}

void ClickRectSleepS(HWND hwnd,  const StructDefs::Rect rect, std::chrono::milliseconds duration = std::chrono::milliseconds(Constants::ClickWaitShort)) {
    ClickPointSleep(hwnd, NewPointFromRectCenter(rect), duration);
}

void ClickRectSleepM(HWND hwnd,  const StructDefs::Rect rect, std::chrono::milliseconds duration = std::chrono::milliseconds(Constants::ClickWaitMid)) {
    ClickPointSleep(hwnd, NewPointFromRectCenter(rect), duration);
}

void ClickRectSleepL(HWND hwnd,  const StructDefs::Rect rect, std::chrono::milliseconds duration = std::chrono::milliseconds(Constants::ClickWaitLong)) {
    ClickPointSleep(hwnd, NewPointFromRectCenter(rect), duration);
}

void UpdateAngle(HWND hwnd, int angle) {
    auto direction = Constants::DirectionUp;
    if (angle < 0) {
        direction = Constants::DirectionDown;
        angle = -angle;
    }
    switch (direction) {
    case Constants::DirectionUp:
        for (int i = 0; i < angle; ++i) {
            KeyBoard(hwnd, Constants::VKCode_UP);
        }
        break;
    case Constants::DirectionDown:
        for (int i = 0; i < angle; ++i) {
            KeyBoard(hwnd, Constants::VKCode_DOWN);
        }
        break;
    default:
        break;
    }
}

void Launch(HWND hwnd, int power) {
    if (power < 0) {
        power = 0;
    }
    if (power > 100) {
        power = 100;
    }
    auto ts = std::chrono::milliseconds(power * 40);
    KeyBoard(hwnd, Constants::VKCode_SPACE, ts);
}

#endif // UTILS_DDT_H
