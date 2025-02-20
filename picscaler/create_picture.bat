@ECHO OFF
PUSHD %~dp0

call env.bat
IF ERRORLEVEL 1 EXIT /b %ERRORLEVEL%

echo PIC_IN_FOLDER: %PIC_IN_FOLDER%
go run main.go -in=%PIC_IN_FOLDER% -out=../mode/gen_picture -pixelCount=58 -asByte

POPD