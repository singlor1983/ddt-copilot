QT       += core gui

greaterThan(QT_MAJOR_VERSION, 4): QT += widgets

CONFIG += c++17

# You can make your code fail to compile if it uses deprecated APIs.
# In order to do so, uncomment the following line.
#DEFINES += QT_DISABLE_DEPRECATED_BEFORE=0x060000    # disables all the APIs deprecated before Qt 6.0.0

SOURCES += \
    app/core/logic.cpp \
    main.cpp \
    mainwindow.cpp \
    utils/utils.cpp \
    utils/utils_ddt.cpp

HEADERS += \
    app/core/config.h \
    app/core/core.h \
    app/core/handler.h \
    app/core/logic.h \
    defs/constants.h \
    defs/struct.h \
    mainwindow.h \
    utils/utils.h \
    utils/utils_ddt.h

FORMS += \
    mainwindow.ui

LIBS += -lgdi32 # user32已自动链接

# Default rules for deployment.
qnx: target.path = /tmp/$${TARGET}/bin
else: unix:!android: target.path = /opt/$${TARGET}/bin
!isEmpty(target.path): INSTALLS += target
