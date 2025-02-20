@ECHO OFF

if "%COMPUTERNAME%" EQU "DESKTOP-DEAN" (
  SET PIC_IN_FOLDER=D:\GDrive\Poi\pictures\active
) else if "%COMPUTERNAME%" EQU "ASUSTUF" (
  SET PIC_IN_FOLDER=X:\GDrive\Poi\pictures\active
) else (
  echo ERROR: unknown computer %COMPUTERNAME%
  EXIT /b 1
)
EXIT /b 0