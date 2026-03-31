@echo off
echo Starting Template Service...
go run main.go service.go handlers.go example_usage.go --config ./config.json --debug
pause