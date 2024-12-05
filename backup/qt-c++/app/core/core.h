#ifndef CORE_H
#define CORE_H

#include <iostream>
#include <thread>
#include <mutex>
#include <map>
#include <stdexcept>
#include <utils/utils_ddt.h>


struct Statistics {
    int totalCount;             // 挑战次数
    int winCount;               // 胜利次数
    int failCount;              // 失败次数
};

struct RoundCache {
    int roundCount;                           // 当前回合
    Constants::FubenInitPostion initPosition; // 初始位置
};


struct PositionNode {
    Constants::FubenType type; // 副本类型
    int page;                  // 处于第几页
    int count;                 // 处于当前页的第几个
};

class ChildNode {
public:
    ChildNode(HWND hwnd) {
        this->hwnd = hwnd;
        this->readyState = Constants::ReadyNo;
    }

private:
    HWND hwnd; // 窗口句柄
    Constants::ReadyState readyState;
};

class ScriptCtrl {
public:
    // 构造函数，使用初始化列表进行初始化
    ScriptCtrl(HWND hwnd, Constants::FunctionID functionID, Constants::FubenLv lv, bool isChild, HWND masterHWND) {
        this->hwnd = hwnd;
        this->functionID = functionID;
        this->lv = lv;
        this->isChild = isChild;
        this->masterHWND = masterHWND;
    }

    void Release() {
        this->childs.clear();
    }

    void Stop() {
        this->running = false;
    }

    void Run() {
        std::thread([this]() {
            try {
                if (functionID > Constants::FunctionFubenStart && functionID < Constants::FunctionFubenEnd) {
                    runFuben();
                } else if (functionID > Constants::FunctionJinjiStart && functionID < Constants::FunctionJinjiEnd) {
                    runJinji();
                }
            } catch (const std::exception& e) {
                std::cerr << "hwnd: " << hwnd << " exit. reason: " << e.what() << std::endl;
            } catch (...) {
                std::cerr << "hwnd: " << hwnd << " exit. reason: unknown exception" << std::endl;
            }
        }).detach(); // 创建并分离子线程
    }

private:
    HWND hwnd;                        // 窗口句柄
    Constants::FunctionID functionID; // 脚本功能ID
    Constants::FubenLv lv;            // 副本难度
    bool isChild;                     // 是否为小号
    HWND masterHWND;                  // 所属主号句柄

    std::mutex mu;              // 互斥锁用于保护共享资源，std::mutex 的默认构造函数会在创建对象时自动调用
    std::atomic<bool> running;  // 脚本是否运行

    Statistics statistics; // 统计数据
    RoundCache roundCache; // 回合数据

    std::map<HWND, ChildNode*> childs; // 小号列表

    void runFuben() {

    }

    void runJinji() {

    }

};


class ScriptCtrlMgr {
public:
    static ScriptCtrlMgr* Instance() {
        if (instance == nullptr) {
            instance = new ScriptCtrlMgr(); // 创建单例实例
            instance->init();
        }
        return instance;
    }

    // 设置副本位置配置
    void SetFubenPosition(const std::map<Constants::FunctionID, PositionNode> &position) {
        this->positionDefs = position;
    }

    // 获取当前屏幕的角度
    int GetAgnle(HWND hwnd) {
        auto captureRect = ToWinRECT(Constants::RectAngle);
        auto bitmap = CaptureWindowLight(hwnd, &captureRect);
        ConvertBitmapColorWithNormalization(bitmap);
        for (const auto& pair : this->agnleDefs) {
            if (IslegalPoint( FindPic(bitmap, pair.second, 0.9))) {
                return pair.first;
            }
        }
        return 0;
    }

    // 执行选择副本的动作
    void DoSelectFuben(HWND hwnd, Constants::FunctionID id, Constants::FubenLv lv, bool isBossFight = true) {
        try {
            auto node = positionDefs.at(id); // 如果 "key" 不存在，会抛出异常
            // 先取消遮挡
            FocusDDTHwnd(hwnd, true);
            // 点击选择副本按钮
            ClickPointSleepM(hwnd, Constants::ClickPointFubenSelect);
            // 选择副本类型
            switch (node.type) {
            case Constants::TpSpecial:
                ClickPointSleepM(hwnd, Constants::ClickPointFubenTypeSpecial);
                break;
            default:
                break;
            }
            // 跳转到副本指定页面
            for (int i = 1; i < node.page; ++i) {
                int nextPageClickTimes = 8;
                if (i%2 == 0 ){
                    nextPageClickTimes = 9;
                }
                for (int j = 0; j < nextPageClickTimes; ++j) {
                    ClickPointSleep(hwnd, Constants::ClickPointFubenPageDown);
                }
            }
            // 选择副本
            ClickRectSleepM(hwnd, Constants::RectFubenNameMap.at(node.count));
            // 选择难度
            auto bitmapB = fubenLvDefs.at(lv);
            for (const auto& rect : Constants::RectFubenLvList) {
                auto winRect = ToWinRECT(rect);
                auto bitmapA = CaptureWindow(hwnd, &winRect);
                auto point = FindPic(bitmapA, bitmapB, 0.9);
                DeleteObject(bitmapA);
                if (IslegalPoint(point)) {
                    ClickPointSleepM(hwnd, {rect.x+point.x, rect.y+point.y});
                    break;
                }
            }
            DeleteObject(bitmapB);
            // 确认选择
            if (isBossFight) {
                ClickPointSleepM(hwnd, Constants::ClickPointFubenBossFight);
                ClickPointSleepM(hwnd, Constants::ClickPointFubenSelectAck);
                ClickPointSleepM(hwnd, Constants::ClickPointFubenBossFightAck);
            } else {
                ClickPointSleepM(hwnd, Constants::ClickPointFubenSelectAck);
            }
        } catch (const std::out_of_range& e) {
            std::cerr << "Key not found: " << e.what() << std::endl;
        }
    }

    void AddScriptCtrl(HWND hwnd, ScriptCtrl* ctrl) {
        items[hwnd] = ctrl;
    }

    ScriptCtrl* GetScriptCtrl(HWND hwnd) {
        auto it = items.find(hwnd);
        return (it != items.end()) ? it->second : nullptr; // 查找并返回
    }

private:
    ScriptCtrlMgr() {}                 // 私有构造函数
    static ScriptCtrlMgr* instance;    // 单例实例
    std::map<HWND, ScriptCtrl*> items; // 存储 ScriptCtrl 实例的容器

    std::map<Constants::FunctionID, PositionNode> positionDefs; // 副本位置配置
    std::map<int, HBITMAP> agnleDefs;                  // 角度识别
    std::map<Constants::FubenLv, HBITMAP> fubenLvDefs; // 难度识别

    void init () {
        // 加载角度
        for (int i = -31; i <= 100; ++i) {
            std::string filename = format("asserts/{}.bmp", i);
            HBITMAP bmp = LoadBitmapFromBmp(StringToWString(filename));
            if (bmp == NULL) {
                continue;
            }
            agnleDefs[i] = bmp;
            std::cout << "load defs agnle success. filename: " << filename << std::endl;
        }
        // 加载副本等级
        for (int i = 1; i <= 4; ++i) {
            std::string filename = format("asserts/f{}.bmp", i);
            HBITMAP bmp = LoadBitmapFromBmp(StringToWString(filename));
            if (bmp == NULL) {
                continue;
            }
            fubenLvDefs[static_cast<Constants::FubenLv>(i)] = bmp;
            std::cout << "load defs fubenlv success. filename: " << filename << std::endl;
        }
    }
};

// 初始化静态成员变量，静态成员变量在类的定义中声明，但它们不是类的每个实例的一部分，而是类本身的一个成员，必须在类外定义和初始化该静态成员变量
// 访问方式： ScriptCtrlMgr::Instance()->AddScriptCtrl(hwnd, ctrl);
ScriptCtrlMgr* ScriptCtrlMgr::instance = nullptr;


#endif // CORE_H
