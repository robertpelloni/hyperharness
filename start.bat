@echo off
setlocal enabledelayedexpansion

:: ═══════════════════════════════════════════════════════════════
:: HyperHarness — AI Coding Agent (Start Script)
:: ═══════════════════════════════════════════════════════════════

echo.
echo   ╔══════════════════════════════════════════════════════╗
echo   ║          🧠  HyperHarness  —  AI Control Plane      ║
echo   ╚══════════════════════════════════════════════════════╝
echo.

:: ─── Configuration ──────────────────────────────────────────
set "BINARY=hyperharness.exe"
set "PORT=8080"

:: ─── Check for existing binary ──────────────────────────────
if not exist "%BINARY%" (
    echo   [Build] Binary not found. Building HyperHarness...
    echo.
    go build -buildvcs=false -o %BINARY% .
    if errorlevel 1 (
        echo.
        echo   [Error] Build failed! Check the errors above.
        pause
        exit /b 1
    )
    echo.
    echo   [Build] Build successful!
    echo.
)

:: ─── Kill old processes on port ─────────────────────────────
echo   [Port] Checking port %PORT%...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":%PORT% " ^| findstr "LISTENING"') do (
    echo   [Port] Killing old server PID %%a on port %PORT%...
    taskkill /PID %%a /F >nul 2>&1
    timeout /t 1 /nobreak >nul
)

:: ─── Parse command ──────────────────────────────────────────
set "CMD=%1"
if "%CMD%"=="" set "CMD=tui"

if "%CMD%"=="tui" goto :run_tui
if "%CMD%"=="serve" goto :run_serve
if "%CMD%"=="pipe" goto :run_pipe
if "%CMD%"=="build" goto :run_build
if "%CMD%"=="test" goto :run_test
if "%CMD%"=="clean" goto :run_clean
if "%CMD%"=="help" goto :show_help

echo   Unknown command: %CMD%
goto :show_help

:: ─── TUI Mode ──────────────────────────────────────────────
:run_tui
echo   [Mode] Interactive TUI
echo.
%BINARY%
goto :end

:: ─── Serve Mode ─────────────────────────────────────────────
:run_serve
echo   [Mode] API Server on port %PORT%
echo.
%BINARY% serve
goto :end

:: ─── Pipe Mode ──────────────────────────────────────────────
:run_pipe
set "QUERY=%2"
if "%QUERY%"=="" (
    echo   [Error] Pipe mode requires a query. Usage: start.bat pipe "your query"
    goto :end
)
echo   [Mode] Pipe: %QUERY%
echo.
echo %QUERY% | %BINARY% pipe
goto :end

:: ─── Build ──────────────────────────────────────────────────
:run_build
echo   [Mode] Rebuilding...
echo.
go build -buildvcs=false -o %BINARY% .
if errorlevel 1 (
    echo   [Error] Build failed!
    pause
    exit /b 1
)
echo.
echo   [Build] Successful!
goto :end

:: ─── Test ───────────────────────────────────────────────────
:run_test
echo   [Mode] Running tests...
echo.
go test ./... 2>&1 | findstr /V "no test files"
goto :end

:: ─── Clean ──────────────────────────────────────────────────
:run_clean
echo   [Mode] Cleaning...
del /q %BINARY% 2>nul
go clean
echo   [Clean] Done.
goto :end

:: ─── Help ───────────────────────────────────────────────────
:show_help
echo.
echo   Usage: start.bat [command]
echo.
echo   Commands:
echo     tui     Start interactive TUI (default)
echo     serve   Start API server on port %PORT%
echo     pipe    One-shot pipe mode (start.bat pipe "query")
echo     build   Rebuild the binary
echo     test    Run all tests
echo     clean   Remove build artifacts
echo     help    Show this help message
echo.
echo   Keyboard Shortcuts (TUI):
echo     Enter       Send message
echo     Tab         Autocomplete slash commands
echo     Ctrl+C      Cancel / quit
echo     Ctrl+L      Toggle file tree pane
echo     Ctrl+D      Toggle dashboard
echo     Ctrl+Y      Accept shell proposal
echo     /help       Show all slash commands
echo     /hotkeys    Show all keybindings
echo     ??query     Shell command proposal
echo.

:end
endlocal
