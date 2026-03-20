@echo off
echo Starting Borg...

where pnpm >nul 2>nul
if %errorlevel% neq 0 (
    echo pnpm could not be found. Installing...
    npm install -g pnpm
)

echo Installing dependencies...
call pnpm install
if %errorlevel% neq 0 exit /b %errorlevel%

echo Building...
call pnpm run build
if %errorlevel% neq 0 exit /b %errorlevel%

echo Starting Hub...
call pnpm start
if %errorlevel% neq 0 exit /b %errorlevel%
