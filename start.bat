@echo off
setlocal
title HyperHarness
cd /d "%~dp0"

echo [HyperHarness] Starting...
where go >nul 2>nul
if errorlevel 1 (
    echo [HyperHarness] go not found. Please install it.
    pause
    exit /b 1
)

go run . serve

if errorlevel 1 (
    echo [HyperHarness] Exited with error code %errorlevel%.
    pause
)
endlocal
