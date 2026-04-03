#include <QApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
// Include the bobui OmniTerminal header
#include "../../bobui/OmniUI/omnicore/include/OmniTerminal.h"

int main(int argc, char *argv[])
{
    QApplication app(argc, argv);

    // Register the bobui OmniTerminal to the QML engine
    // This allows native terminal rendering with full PTY support directly in Qt6
    qmlRegisterType<OmniTerminal>("Maestro.Components", 1, 0, "OmniTerminal");

    QQmlApplicationEngine engine;
    
    const QUrl url(QStringLiteral("qrc:/Maestro/assets/main.qml"));
    QObject::connect(&engine, &QQmlApplicationEngine::objectCreated,
                     &app, [url](QObject *obj, const QUrl &objUrl) {
        if (!obj && url == objUrl)
            QCoreApplication::exit(-1);
    }, Qt::QueuedConnection);
    engine.load(url);

    return app.exec();
}
