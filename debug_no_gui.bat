@REM Pay attention to row size: this doubles the number of LEDs
go build -o ledean.exe && ledean.exe -gpio_button=GPIO17 -led_count=58 -led_rows=2 -reverse_rows=0,1 -path2frontend="" -log_level="debug" -direct_start -no_gui
@REM go build -o ledean.exe && ledean.exe -gpio_button=GPIO17 -led_count=36 -path2frontend="" -log_level="debug" -direct_start -show_mode