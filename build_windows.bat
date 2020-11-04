set GOARCH=amd64
set GOOS=windows

echo building for: %GOOS%(%GOARCH%)...

go build -o build/%GOOS%/LEDean.sh
echo done