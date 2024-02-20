@echo off

set GOOS=linux

@REM Pi 4
set GOARCH=arm64
set GOARM=7

@REM Pi Zero
@REM set GOARCH=arm
@REM set GOARM=6

echo building for: %GOOS%(%GOARCH%)...

go build -o build/%GOOS%/ledean -tags "pi"

echo Exit Code is %errorlevel%
echo.

if %ERRORLEVEL% EQU 0 (
   echo Success 
) else (
  echo Fail
  exit /b %errorlevel%
)