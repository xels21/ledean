@REM set GOARCH=arm64
set GOARCH=arm
set GOOS=linux
set GOARM=7

echo building for: %GOOS%(%GOARCH%)...

go build -o build/%GOOS%/ledean
echo done