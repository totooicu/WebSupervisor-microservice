@echo off
setlocal

echo ================================================
echo WebSupervisor Microservice Build Script
echo Release Version Build
echo ================================================

cd "%~dp0"

REM Set Go environment variables for release build
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
set LDFLAGS=-ldflags="-s -w"

echo Building release version with optimizations...
echo CGO_ENABLED=%CGO_ENABLED%
echo GOOS=%GOOS%
echo GOARCH=%GOARCH%
echo GOFLAGS=%GOFLAGS%
echo.

REM Create release directory
if not exist "release" mkdir release

echo Building cache-service...
go build %LDFLAGS% -o release/cache-service.exe ./microservice/cache-service
if %errorlevel% neq 0 goto error

echo Building crawler-service...
go build %LDFLAGS% -o release/crawler-service.exe ./microservice/crawler-service
if %errorlevel% neq 0 goto error

echo Building monitor-service...
go build %LDFLAGS% -o release/monitor-service.exe ./microservice/monitor-service
if %errorlevel% neq 0 goto error

echo Building notifier-service...
go build %LDFLAGS% -o release/notifier-service.exe ./microservice/notifier-service
if %errorlevel% neq 0 goto error

echo Building parser-service...
go build %LDFLAGS% -o release/parser-service.exe ./microservice/parser-service
if %errorlevel% neq 0 goto error

echo Building template-service...
go build %LDFLAGS% -o release/template-service.exe ./microservice/template-service
if %errorlevel% neq 0 goto error

echo.
echo ================================================
echo Build completed successfully!
echo Release files are located in: %~dp0release\
echo ================================================

goto end

:error
echo.
echo ================================================
echo Build failed!
echo ================================================
exit /b 1

:end
endlocal

