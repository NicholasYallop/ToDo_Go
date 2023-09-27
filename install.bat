@ECHO OFF

REM build executable
go build -C todo -o ../todo.exe main.go
if %errorlevel% NEQ 0 (
  ECHO error: Build failed
  EXIT /b %errorlevel%
)

REM create target directory
mkdir %LocalAppData%\todo

REM move executable to taregt directory
move todo.exe %LocalAppData%\todo\todo.exe
if %errorlevel% NEQ 0 (
  ECHO error: Could not move executable to local app data directory
  EXIT /b %errorlevel%
)

REM add target directory to path, if not already present
path|find /i "%LocalAppData%\todo" >nul || set path=%path%;%LocalAppData%\todo
if %errorlevel% NEQ 0 (
  ECHO error: Failed to append LocalAppData/todo to path variable
  EXIT /b %errorlevel%
)
