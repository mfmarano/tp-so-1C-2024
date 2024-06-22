@echo off

echo Running kernel.go in a new console window...
start cmd /k "cd kernel && go run kernel.go"

echo Waiting for kernel.go to load (10 seconds)...
timeout /t 10 /nobreak >nul

echo Running entradasalida.go with TECLADO.json in a new console window...
start cmd /k "cd entradasalida && go run entradasalida.go TECLADO.json"

echo Running entradasalida.go with GENERICA.json in a new console window...
start cmd /k "cd entradasalida && go run entradasalida.go GENERICA.json"

echo Running entradasalida.go with MONITOR.json in a new console window...
start cmd /k "cd entradasalida && go run entradasalida.go MONITOR.json"

echo Running memoria.go in a new console window...
start cmd /k "cd memoria && go run memoria.go"

echo Waiting for memoria.go to load (10 seconds)...
timeout /t 10 /nobreak >nul

echo Running cpu.go in a new console window...
start cmd /k "cd cpu && go run cpu.go"

echo.
echo Script execution initiated.
pause