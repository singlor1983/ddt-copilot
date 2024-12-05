#ifndef CONSTANTS_H
#define CONSTANTS_H

#include <QString>
#include <defs/struct.h>
#include <map>

namespace Constants {
    const QString AppName ="3.6经典弹弹堂通用脚本";
    const int MaxMainWindowSizeW = 1150;
    const int MaxMainWindowSizeH = 800;
    const std::wstring ProcessNameTGWeb = L"TangoWeb.exe";

    const StructDefs::Point PointNoEffect = {1, 1};  // 无效点击点
    const StructDefs::Point PointIllegal = {-1, -1}; // 非法的坐标

    const StructDefs::Point ClickPointFubenEnter = {722, 479};         // 从副本大厅进入副本房间
    const StructDefs::Point ClickPointFubenHall = {869, 501};          // 进入副本大厅的按钮
    const StructDefs::Point ClickPointFubenSelect = {597, 234};        // 副本内-选择副本
    const StructDefs::Point ClickPointFubenPageDown = {768, 440};      // 副本选择页面-下拉条滑动
    const StructDefs::Point ClickPointFubenTypeNormal = {287, 305};    // 副本选择页面-普通副本
    const StructDefs::Point ClickPointFubenTypeSpecial = {427, 305};   // 副本选择页面-特殊副本
    const StructDefs::Point ClickPointFubenBossFight = {320, 565};     // 副本选择页面-BOSS战勾选框
    const StructDefs::Point ClickPointFubenBossFightAck = {433, 340};  // 副本选择页面-BOSS战确认按钮
    const StructDefs::Point ClickPointFubenSelectAck = {500, 565};     // 副本选择页面-副本选择确认按钮
    const StructDefs::Point ClickPointFightEquipItem1 = {810, 140};    // 房间内-道具装备区1
    const StructDefs::Point ClickPointFightEquipItem2 = {880, 140};    // 房间内-道具装备区2
    const StructDefs::Point ClickPointFightEquipItem3 = {950, 140};    // 房间内-道具装备区3
    const StructDefs::Point ClickPointFightSelectItem1 = {807, 234};   // 房间内-道具选择区1-1
    const StructDefs::Point ClickPointFightSelectItem2 = {857, 234};   // 房间内-道具选择区1-2
    const StructDefs::Point ClickPointFightSelectItem3 = {907, 234};   // 房间内-道具选择区1-3
    const StructDefs::Point ClickPointFightSelectItem4 = {957, 234};   // 房间内-道具选择区1-4
    const StructDefs::Point ClickPointFightSelectItem5 = {807, 289};   // 房间内-道具选择区2-1
    const StructDefs::Point ClickPointFightSelectItem6 = {857, 289};   // 房间内-道具选择区2-2
    const StructDefs::Point ClickPointFightSelectItem7 = {907, 289};   // 房间内-道具选择区2-3
    const StructDefs::Point ClickPointFightSelectItem8 = {957, 289};   // 房间内-道具选择区2-4
    const StructDefs::Point ClickPointFightStart = {940, 500};         // 战斗开始按钮，副本和竞技是同一个位置的按钮
    const StructDefs::Point ClickPointFubenFightStartAck = {413, 339}; // 副本战斗开始确认按钮
    const StructDefs::Point ClickPointBackAndExit = {965, 570};        // 导航栏-返回&退出按钮


    // 以下为DDT元素出现的区域，相对于窗口左上角
    const StructDefs::Rect RectPassBtn = {477, 159, 47, 16};                   // 战斗内-PASS按钮
    const StructDefs::Rect RectAngle = {32, 558, 35, 15};                      // 力度识别区
    const StructDefs::Rect RectFubenSelectText = {540, 110, 100, 20};          // 副本房间-选择副本
    const StructDefs::Rect RectFubenInviteAndChangeTeam = {761, 453, 125, 65}; // 副本房间特征值-邀请&换队
    const StructDefs::Rect RectJinjiInviteAndChangeArea = {761, 453, 125, 65}; // 竞技房间特征值-邀请&换区
    const StructDefs::Rect RectFubenHall = {30, 83, 58, 16};                   // 副本大厅特征值-筛选副本
    const StructDefs::Rect RectJinjiHall = {25, 60, 188, 25};                  // 竞技大厅特征值-所有模式
    const StructDefs::Rect RectFightRightTop = {946, 4, 50, 16};               // 在竞技战斗或副本内的特征元素【右上角的设置和退出按钮】
    const StructDefs::Rect RectFightResult = {678, 34, 125, 30};               // 战斗结算特征元素【我的成绩，竞技战和副本战都用到、胜利和失败都有】
    const StructDefs::Rect RectFightLoading = {396, 301, 190, 62};             // 副本战斗加载特征元素【副本战、自由战】
    const StructDefs::Rect RectFubenFightSettle = {159, 17, 149, 38};          // 副本战斗结算特征元素【游戏结算，左上角】也是boss关翻牌画面
    const StructDefs::Rect RectJinjiFightSettle = {749, 32, 149, 38};          // 竞技战斗结算特征元素【游戏结算，右上角】也是小关翻牌画面
    const StructDefs::Rect RectBackAndExit = {945, 580, 38, 18};               // 导航栏-退出&返回
    const StructDefs::Rect RectMiniMap = {788, 24, 211, 97};                   // 战斗内-小地图
    const StructDefs::Rect RectIsYourTurn = {460, 225, 100, 25};               // 战斗内-轮到你的回合了
    const StructDefs::Rect RectWinOrFail = {821, 13, 155, 104};                // 战斗内-结算 胜利还是失败识别区
    const StructDefs::Rect RectReadyState = {893, 507, 90, 35};                // 房间内-准备状态识别区
    const StructDefs::Rect RectFubenNameBtn1 = {219, 332, 130, 49};                                        // 副本选择-副本名字识别区1-1
    const StructDefs::Rect RectFubenNameBtn2 = {RectFubenNameBtn1.x+130+6, 332, 130, 49};                  // 副本选择-副本名字识别区1-2
    const StructDefs::Rect RectFubenNameBtn3 = {RectFubenNameBtn2.x+130+6, 332, 130, 49};                  // 副本选择-副本名字识别区1-3
    const StructDefs::Rect RectFubenNameBtn4 = {RectFubenNameBtn3.x+130+6, 332, 130, 49};                  // 副本选择-副本名字识别区1-4
    const StructDefs::Rect RectFubenNameBtn5 = {219, 332+49+15, 130, 49};                                  // 副本选择-副本名字识别区2-1
    const StructDefs::Rect RectFubenNameBtn6 = {RectFubenNameBtn5.x+130+6, RectFubenNameBtn5.y, 130, 49};  // 副本选择-副本名字识别区2-2
    const StructDefs::Rect RectFubenNameBtn7 = {RectFubenNameBtn6.x+130+6, RectFubenNameBtn5.y, 130, 49};  // 副本选择-副本名字识别区2-3
    const StructDefs::Rect RectFubenNameBtn8 = {RectFubenNameBtn7.x+130+6, RectFubenNameBtn5.y, 130, 49};  // 副本选择-副本名字识别区2-4
    const StructDefs::Rect RectFubenLvBtn1 = {437, 489, 110, 30};  // 副本选择-副本难度识别区1
    const StructDefs::Rect RectFubenLvBtn2 = {325, 489, 110, 30};  // 副本选择-副本难度识别区2
    const StructDefs::Rect RectFubenLvBtn3 = {260, 489, 110, 30};  // 副本选择-副本难度识别区3
    const StructDefs::Rect RectFubenLvBtn4 = {232, 489, 110, 30};  // 副本选择-副本难度识别区4
    const StructDefs::Rect RectFubenLvBtn5 = {565, 489, 110, 30};  // 副本选择-副本难度识别区5
    const StructDefs::Rect RectFubenLvBtn6 = {440, 489, 110, 30};  // 副本选择-副本难度识别区6
    const StructDefs::Rect RectFubenLvBtn7 = {366, 489, 110, 30};  // 副本选择-副本难度识别区7
    const StructDefs::Rect RectFubenLvBtn8 = {620, 489, 110, 30};  // 副本选择-副本难度识别区8
    const StructDefs::Rect RectFubenLvBtn9 = {496, 489, 110, 30};  // 副本选择-副本难度识别区9
    const StructDefs::Rect RectFubenLvBtn10 = {632, 489, 110, 30}; // 副本选择-副本难度识别区10

    const std::map<int, StructDefs::Rect> RectFubenNameMap = {
        {1, RectFubenNameBtn1},
        {2, RectFubenNameBtn2},
        {3, RectFubenNameBtn3},
        {4, RectFubenNameBtn4},
        {5, RectFubenNameBtn5},
        {6, RectFubenNameBtn6},
        {7, RectFubenNameBtn7},
        {8, RectFubenNameBtn8},
    };

    const std::vector<StructDefs::Rect> RectFubenLvList = {
        RectFubenLvBtn1,
        RectFubenLvBtn2,
        RectFubenLvBtn3,
        RectFubenLvBtn4,
        RectFubenLvBtn5,
        RectFubenLvBtn6,
        RectFubenLvBtn7,
        RectFubenLvBtn8,
        RectFubenLvBtn9,
        RectFubenLvBtn10,
    };

    enum FunctionID {
        FunctionFubenStart = 0,
        FunctionFubenMaYi = 1,
        FunctionFubenXiaoji = 2,
        FunctionFubenBoGu = 3,
        FunctionFubenXieShen = 4,
        FunctionFubenBaoLei = 5,
        FunctionFubenLongChao = 6,
        FunctionFubenYunDongHui = 7,
        FunctionFubenJinjiChang = 8,
        FunctionFubenEnd = 1000,
        FunctionJinjiStart = 1000,
        FunctionJinjiEnd = 2000,
    };

    enum FubenLv {
        LvEasy = 1,       // 简单
        LvNormal = 2,     // 普通
        LvDifficulty = 3, // 困难
        LvHero = 4,       // 英雄
        LvNightmare = 5,  // 噩梦
    };

    enum FubenType {
        TpNormal = 1,  // 普通副本
        TpSpecial = 2, // 特殊副本
    };

    enum FubenInitPostion {
        Position1 = 1, // 初始1号位
        Position2 = 2, // 初始2号位
        Position3 = 3, // 初始3号位
        Position4 = 4, // 初始4号位
    };

    enum VKCode {
        VKCode_0       = 48,
        VKCode_1       = 49,
        VKCode_2       = 50,
        VKCode_3       = 51,
        VKCode_4       = 52,
        VKCode_5       = 53,
        VKCode_6       = 54,
        VKCode_7       = 55,
        VKCode_8       = 56,
        VKCode_9       = 57,
        VKCode_B       = 66,
        VKCode_Q       = 81,
        VKCode_E       = 69,
        VKCode_T       = 84,
        VKCode_Y       = 89,
        VKCode_U       = 85,
        VKCode_P       = 80,
        VKCode_F       = 70,
        VKCode_SPACE   = 32,
        VKCode_LEFT    = 37,
        VKCode_UP      = 38,
        VKCode_RIGHT   = 39,
        VKCode_DOWN    = 40,
        VKCode_ESCAPE  = 27,
    };

    const std::map<std::string, VKCode> VKMap = {
        {"0", VKCode_0},
        {"1", VKCode_1},
        {"2", VKCode_2},
        {"3", VKCode_3},
        {"4", VKCode_4},
        {"5", VKCode_5},
        {"6", VKCode_6},
        {"7", VKCode_7},
        {"8", VKCode_8},
        {"9", VKCode_9},
        {"B", VKCode_B},
        {"Q", VKCode_Q},
        {"E", VKCode_E},
        {"T", VKCode_T},
        {"Y", VKCode_Y},
        {"U", VKCode_U},
        {"P", VKCode_P},
        {"F", VKCode_F},
        {"SPACE", VKCode_SPACE},
        {"LEFT", VKCode_LEFT},
        {"UP", VKCode_UP},
        {"RIGHT", VKCode_RIGHT},
        {"DOWN", VKCode_DOWN},
        {"ESCAPE", VKCode_ESCAPE},
    };

    enum ReadyState {
        ReadyNo = 0,
        ReadyOk = 1,
    };

    enum Direction {
        DirectionLeft = 1,
        DirectionRight = 2,
        DirectionUp = 3,
        DirectionDown = 4,
    };

    const int ClickWaitShort = 150;
    const int ClickWaitMid = 700;
    const int ClickWaitLong = 2000;
}

#endif // CONSTANTS_H
