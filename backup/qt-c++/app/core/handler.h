#ifndef HANDLER_H
#define HANDLER_H

#include <iostream>
#include <functional>
#include <app/core/core.h>

using FnBeofreStart = std::function<void(ScriptCtrl*)>;
using FnBattle = std::function<void(ScriptCtrl*)>;
using FnFirstRoundInit = std::function<void(ScriptCtrl*)>;

struct Handle {
    Constants::FunctionID id;

    FnBeofreStart cbBeforeStart; // 开始前的回调
    FnBattle cbBattle; // 战斗回调
    FnFirstRoundInit cbFirstRoundInit; // 第一轮初始化回调

    // 结构体的构造函数
    Handle(Constants::FunctionID id, FnBeofreStart cbBeforeStart, FnBattle cbBattle, FnFirstRoundInit cbFirstRoundInit)
        : id(id), cbBeforeStart(cbBeforeStart), cbBattle(cbBattle), cbFirstRoundInit(cbFirstRoundInit) {}
};

class HandlerMgr {
public:
    static HandlerMgr* Instance() {
        if (instance == nullptr) {
            instance = new HandlerMgr(); // 创建单例实例
            instance->init();
        }
        return instance;
    }

    void init () {

    }

    /*
     *  operator[]：用于访问或插入元素，若不存在会插入默认值。
        at：用于安全访问元素，若不存在会抛出异常。
        find：用于查找元素，返回迭代器，适合检查键的存在性而不修改容器。
     */
    void ResgiterHandler(Constants::FunctionID id, FnBattle cbBattle, FnBeofreStart cbBeforeStart = nullptr, FnFirstRoundInit cbFirstRoundInit = nullptr) {
        if (handlers.find(id) != handlers.end()) {
            throw std::runtime_error("handle: " + std::to_string(id) + " already exists");
        }

        Handle handler(id, cbBeforeStart, cbFirstRoundInit, cbBattle);
        handlers[id] = handler;
    }

private:
    HandlerMgr() {}                 // 私有构造函数
    static HandlerMgr* instance;    // 单例实例

    std::map<Constants::FunctionID, Handle> handlers;

};

HandlerMgr* HandlerMgr::instance = nullptr;

#endif // HANDLER_H
