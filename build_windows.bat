@echo off

set GOARCH=amd64
set GOOS=windows

echo building for: %GOOS%(%GOARCH%)...

go build -o build/%GOOS%/LEDean.sh

echo Exit Code is %errorlevel%
echo.

if %ERRORLEVEL% EQU 0 (
   echo Success 
) else (
  echo Fail
  exit /b %errorlevel%
)