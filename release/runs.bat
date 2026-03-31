@echo off
setlocal

echo ================================================
echo WebSupervisor Microservice Release Version
echo Starting all services...
echo ================================================

cd "%~dp0"

REM 设置工作目录
set BASE_DIR=%~dp0..
set CONFIG_DIR=%BASE_DIR%\microservice

echo Starting cache-service...
start "cache-service" %BASE_DIR%\release\cache-service.exe --config %CONFIG_DIR%\cache-service\config.json

echo Starting crawler-service...
start "crawler-service" %BASE_DIR%\release\crawler-service.exe --config %CONFIG_DIR%\crawler-service\config.json

echo Starting monitor-service...
start "monitor-service" %BASE_DIR%\release\monitor-service.exe --config %CONFIG_DIR%\monitor-service\config.json

echo Starting notifier-service...
start "notifier-service" %BASE_DIR%\release\notifier-service.exe --config %CONFIG_DIR%\notifier-service\config.json

echo Starting parser-service...
start "parser-service" %BASE_DIR%\release\parser-service.exe --config %CONFIG_DIR%\parser-service\config.json

echo Starting template-service...
start "template-service" %BASE_DIR%\release\template-service.exe --config %CONFIG_DIR%\template-service\config.json

echo.
echo ================================================
echo All services started successfully!
echo You can monitor each service in its own window.
echo Press any key to exit this window...
echo ================================================

pause >nul

endlocal