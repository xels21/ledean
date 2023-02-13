@REM set Path=%Path%;X:\Software\avr8-gnu-toolchain\bin
@REM set Path=%Path%;X:\Software\avr8-gnu-toolchain\avr\bin
@REM set Path=%Path%;C:\Program Files (x86)\GnuWin32
@REM set Path=%Path%;X:\Software\avrdude


@REM set Path=X:\Software\tinygo\avr8-gnu-toolchain-win32_x86\bin;%Path%

@REM tinygo flash -target=arduino-nano
tinygo build -target=arduino-nano