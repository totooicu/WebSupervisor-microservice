@echo off
setlocal

echo ================================================
echo Starting Template Service Python Version
echo ================================================

REM 检查Python是否安装
python --version >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: Python is not installed or not in PATH
    echo Please install Python 3.8 or higher
    pause
    exit /b 1
)

REM 检查是否已安装依赖
echo Checking dependencies...
pip list | findstr "redis" >nul 2>&1
if %errorlevel% neq 0 (
    echo Installing dependencies...
    pip install -r requirements.txt
    if %errorlevel% neq 0 (
        echo Error: Failed to install dependencies
        pause
        exit /b 1
    )
)

REM 启动服务
echo Starting service...
python main.py --config ./config.json --debug

endlocal