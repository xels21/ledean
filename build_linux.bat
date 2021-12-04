@echo off

@REM set GOARCH=arm64
set GOARCH=arm
set GOOS=linux
@REM Pi 4
@REM set GOARM=7
@REM Pi Zero
set GOARM=6

echo building for: %GOOS%(%GOARCH%)...

go build -o build/%GOOS%/ledean

echo Exit Code is %errorlevel%
echo.

if %ERRORLEVEL% EQU 0 (
   echo Success 
) else (
  echo Fail
  exit /b %errorlevel%
)