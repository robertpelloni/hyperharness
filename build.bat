@echo off
setlocal

:: HyperHarness Build Script for Windows
set /p VERSION=<VERSION 2>nul || set VERSION=0.5.0-alpha.1
set BINARY=hyperharness.exe
set LDFLAGS=-s -w -X internal/buildinfo.Version=%VERSION%

echo.
echo ============================================================
echo    HyperHarness Build System v%VERSION%
echo ============================================================
echo.

if "%1"=="clean" (
    echo Cleaning...
    del /q %BINARY% 2>nul
    go clean -cache 2>nul
)

echo Building...
go build -buildvcs=false -ldflags="%LDFLAGS%" -o %BINARY% .

if %ERRORLEVEL%==0 (
    echo.
    echo Build successful: %BINARY%
    echo.
    echo Run: %BINARY%
    echo    or: %BINARY% tui
    echo    or: %BINARY% serve
) else (
    echo.
    echo Build failed!
    exit /b 1
)

endlocal
