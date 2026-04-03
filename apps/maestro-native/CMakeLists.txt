cmake_minimum_required(VERSION 3.16)
project(MaestroNative LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Point to bobui (Qt6 fork)
set(BOBUI_DIR "${CMAKE_CURRENT_SOURCE_DIR}/../../bobui")
list(APPEND CMAKE_PREFIX_PATH "${BOBUI_DIR}/build_output")

find_package(Qt6 REQUIRED COMPONENTS Core Gui Widgets Qml Quick)

qt_standard_project_setup()

add_executable(MaestroNative
    main.cpp
)

target_link_libraries(MaestroNative PRIVATE
    Qt6::Core
    Qt6::Gui
    Qt6::Widgets
    Qt6::Qml
    Qt6::Quick
)

qt_add_qml_module(MaestroNative
    URI Maestro
    VERSION 1.0
    QML_FILES
        assets/main.qml
)

set_target_properties(MaestroNative PROPERTIES
    WIN32_EXECUTABLE TRUE
    MACOSX_BUNDLE TRUE
)
