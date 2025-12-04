@echo off
echo ===================================
echo   Quick IP Switcher
echo   Building...
echo ===================================
echo.

REM 设置版本号(可选,默认为dev)
set VERSION=dev
if not "%1"=="" set VERSION=%1

echo Building version: %VERSION%
echo.

REM 清理旧文件
if exist quick-ip-switcher.exe del quick-ip-switcher.exe

REM 编译程序 (windowsgui模式,不显示控制台窗口,注入版本号)
go build -ldflags="-H windowsgui -X main.Version=%VERSION%" -o quick-ip-switcher.exe

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ===================================
    echo   Build Success!
    echo   Output: quick-ip-switcher.exe
    echo ===================================
) else (
    echo.
    echo ===================================
    echo   Build Failed!
    echo ===================================
    pause
    exit /b 1
)

echo.
pause
