@echo off
chdir /d %~dp0\..\..
echo Starting Production Service...
go run cmd/production/main.go -config=config/production.yaml
pause

