@echo off
setlocal EnableDelayedExpansion
echo Starting HyperCode 1.0.0-alpha.32 (Go Native Core)...

cd go
where go >nul 2>nul
if errorlevel 1 (
    echo Go could not be found. Please install Go 1.22+.
    exit /b 1
)

set VER=1.0.0-alpha.32
echo Building HyperCode Go Control Plane...
go build -ldflags "-X internal/buildinfo.Version=%VER%" -buildvcs=false -o ../bin/hypercode.exe ./cmd/hypercode
if errorlevel 1 (
    echo Go build failed.
    exit /b 1
)

cd ..
echo Launching HyperCode...
bin\hypercode.exe %*
exit /b %errorlevel%
