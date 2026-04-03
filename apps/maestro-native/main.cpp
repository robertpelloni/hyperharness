#include <QApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>

int main(int argc, char *argv[])
{
    QApplication app(argc, argv);

    QQmlApplicationEngine engine;
    
    // In a real implementation, we would register OmniTerminal from bobui here
    // qmlRegisterType<OmniTerminal>("Maestro", 1, 0, "OmniTerminal");

    const QUrl url(QStringLiteral("qrc:/Maestro/assets/main.qml"));
    QObject::connect(&engine, &QQmlApplicationEngine::objectCreated,
                     &app, [url](QObject *obj, const QUrl &objUrl) {
        if (!obj && url == objUrl)
            QCoreApplication::exit(-1);
    }, Qt::QueuedConnection);
    engine.load(url);

    return app.exec();
}
