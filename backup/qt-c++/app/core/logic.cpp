
//#include <string>
//#include <vector>
//#include <defs/constants.h>
//#include <utils/utils_ddt.h>
//#include <app/core/core.h>

//const std::vector<std::string> DefaultAttackList = {
//    "E", "2", "3",
//    "4", "4", "4", "4", "4", "4", "4", "4",
//    "5", "5", "5", "5", "5", "5", "5", "5",
//    "6", "6", "6", "6", "6", "6", "6", "6",
//    "7", "7", "7", "7", "7", "7", "7", "7",
//    "8", "8", "8", "8", "8", "8", "8", "8",
//};

///*
//    范围-based for 循环：简洁易用，适合大多数情况。
//    传统 for 循环：在需要索引或修改元素时使用。
//    迭代器：提供灵活的遍历方式，适合更复杂的场景。
//    std::for_each 和 Lambda：适合需要复杂操作的情况。
// */
//std::vector<uintptr_t> GetAttackCMD() {
//    std::vector<uintptr_t> vec;
//    // Todo：从全局配置中读取攻击指令
//    auto attackCMD = DefaultAttackList;
//    for (const auto& k : attackCMD) {
//        auto it = Constants::VKMap.find(k);
//        if (it == Constants::VKMap.end()) {
//            continue;
//        }
//        vec.push_back(it->second);
//    }
//    return vec;
//}

//void UseSkillByConfig(HWND hwnd) {
//    auto cmds = GetAttackCMD();
//    for (const auto& k : cmds) {
//        KeyBoard(hwnd, k);
//    }
//}

//void LaunchWithAP(HWND hwnd, int angle, int power) {
//    auto num = ScriptCtrlMgr::Instance()->GetAgnle(hwnd);
//    auto diff = angle - num;
//    UpdateAngle(hwnd, diff);
//    Launch(hwnd, power);
//}
