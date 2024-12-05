#include "mainwindow.h"
#include "ui_mainwindow.h"
#include <vector>
#include <windows.h>
#include <defs/constants.h>
#include <utils/utils_ddt.h>
#include <app/core/core.h>

HWND GetFirstHwnd() {
    std::vector<HWND> hwnds = GetDDTHwnds();
    if (hwnds.empty()) {
        return NULL;
    }
    return hwnds[0];
}

void Test() {
    for (HWND hwnd : GetDDTHwnds()) {
        std::cout << "Found ddt hwnds: " << reinterpret_cast<std::uintptr_t>(hwnd) << std::endl;
        std::string result = format("screenshot_{}.bmp", reinterpret_cast<std::uintptr_t>(hwnd));
        RECT captureRect = {1, 421, 484, 600};
        // auto ret = CaptureWindowAndSave(hwnd, StringToWString(result), &captureRect);
        auto bitmap = CaptureWindowLight(hwnd, &captureRect);
        auto ret = SaveBitmapToBmp(bitmap, StringToWString(result));
        std::cout << format("capture flag: {}", ret) << std::endl;

        ScriptCtrl ctrl(hwnd, Constants::FunctionFubenEnd, Constants::LvEasy, false, nullptr);
        ScriptCtrlMgr::Instance()->AddScriptCtrl(hwnd, &ctrl);
        ctrl.Run();
    }
}

void LoadTest() {
    auto bitmap = LoadBitmapFromBmp(L"screenshot_20975466.bmp");
    PrintBitmapInfo(bitmap);
    PrintBitmapInfo(bitmap);
    auto ret = SaveBitmapToBmp(bitmap, L"loadaftersave.bmp");
    DeleteObject(bitmap); // 释放位图
    std::cout << format("capture flag: {}", ret) << std::endl;
}

void CaputureTest() {
    HWND hwnd = GetFirstHwnd();
    if (hwnd == NULL) {
        std::cout << "not found hwnd" << std::endl;
        return;
    }
    std::cout << "found hwnd:" << reinterpret_cast<std::uintptr_t>(hwnd) << std::endl;
    std::string result = format("screenshot_{}.bmp", reinterpret_cast<std::uintptr_t>(hwnd));
    auto captureRect = ToWinRECT(Constants::RectAngle);
    auto bitmap = CaptureWindowLight(hwnd, &captureRect);
    ConvertBitmapColorWithNormalization(bitmap);
    PrintBitmapInfo(bitmap);
    auto ret = SaveBitmapToBmp(bitmap, StringToWString(result));
    DeleteObject(bitmap); // 释放位图
    std::cout << format("capture flag: {}", ret) << std::endl;
}

void FindPicTest() {
    HWND hwnd = GetFirstHwnd();
    if (hwnd == NULL) {
        std::cout << "not found hwnd" << std::endl;
        return;
    }
    std::cout << "found hwnd:" << reinterpret_cast<std::uintptr_t>(hwnd) << std::endl;
    std::string result = format("screenshot_{}.bmp", reinterpret_cast<std::uintptr_t>(hwnd));
    auto captureRect = ToWinRECT(Constants::RectAngle);
    auto bitmapA = CaptureWindowLight(hwnd, &captureRect);
    ConvertBitmapColorWithNormalization(bitmapA);
    auto ret1 = SaveBitmapToBmp(bitmapA, StringToWString(result));

    auto bitmapB = LoadBitmapFromBmp(L"1.bmp");

    StructDefs::Point point = FindPic(bitmapA, bitmapB, 0.9);
    std::cout << "point:" << point.x << " | " << point.y << std::endl;
}

void GenAngle() {
    HWND hwnd = GetFirstHwnd();
    if (hwnd == NULL) {
        std::cout << "not found hwnd" << std::endl;
        return;
    }
    std::cout << "found hwnd:" << reinterpret_cast<std::uintptr_t>(hwnd) << std::endl;

    int i = 1;
    while (true) {
        std::string result = format("angle_{}.bmp", i);
        auto captureRect = ToWinRECT(Constants::RectAngle);
        auto bitmapA = CaptureWindowLight(hwnd, &captureRect);
        ConvertBitmapColorWithNormalization(bitmapA);
        auto ret = SaveBitmapToBmp(bitmapA, StringToWString(result));
        std::cout << format("capture flag: {}-{}", i, ret) << std::endl;
        i++;
        std::this_thread::sleep_for(std::chrono::seconds(1));
    }
}

void LoadContans() {
    HWND hwnd = GetFirstHwnd();
    if (hwnd == NULL) {
        std::cout << "not found hwnd" << std::endl;
        return;
    }
    std::cout << "found hwnd:" << reinterpret_cast<std::uintptr_t>(hwnd) << std::endl;

    const std::map<Constants::FunctionID, PositionNode> node = {
        {Constants::FunctionFubenMaYi, {Constants::TpNormal, 1, 1}},
        {Constants::FunctionFubenXiaoji, {Constants::TpNormal, 1, 2}},
        {Constants::FunctionFubenBoGu, {Constants::TpNormal, 1, 3}},
        {Constants::FunctionFubenXieShen, {Constants::TpNormal, 1, 4}},
        {Constants::FunctionFubenBaoLei, {Constants::TpNormal, 1, 5}},
        {Constants::FunctionFubenLongChao, {Constants::TpNormal, 1, 6}},
        {Constants::FunctionFubenYunDongHui, {Constants::TpNormal, 1, 7}},
        {Constants::FunctionFubenJinjiChang, {Constants::TpNormal, 1, 8}},
    };

    ScriptCtrlMgr::Instance()->SetFubenPosition(node);
    ScriptCtrlMgr::Instance()->DoSelectFuben(hwnd, Constants::FunctionFubenBoGu, Constants::LvNormal, false);
}

void GenFubenLv() {
    HWND hwnd = GetFirstHwnd();
    if (hwnd == NULL) {
        std::cout << "not found hwnd" << std::endl;
        return;
    }
    std::cout << "found hwnd:" << reinterpret_cast<std::uintptr_t>(hwnd) << std::endl;

    for (size_t i = 0; i <  Constants::RectFubenLvList.size(); ++i) {
        const auto rect =  Constants::RectFubenLvList[i];
        std::string result = format("fubenLv_{}.bmp", i);
        auto captureRect = ToWinRECT(rect);
        auto bitmapA = CaptureWindowLight(hwnd, &captureRect, false);
        auto ret = SaveBitmapToBmp(bitmapA, StringToWString(result));
        std::cout << format("capture flag: {}-{}", i, ret) << std::endl;
    }
}

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::MainWindow)
{
    ui->setupUi(this);

    setWindowTitle(Constants::AppName); // 设置标题名字
    setFixedSize(Constants::MaxMainWindowSizeW, Constants::MaxMainWindowSizeH); // 固定窗口大小

    LoadContans();
}

MainWindow::~MainWindow()
{
    delete ui;
}

