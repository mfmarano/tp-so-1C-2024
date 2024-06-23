@echo off

echo Running kernel.go in a new console window...
start cmd /k "cd ../kernel && go run kernel.go configs/config_PRUEBA_DEADLOCK.json"

echo Waiting for kernel.go to load (10 seconds)...
timeout /t 10 /nobreak >nul

echo Running entradasalida.go with SLP1.json in a new console window...
start cmd /k "cd ../entradasalida && go run entradasalida.go configs/ESPERA.json"

echo Running memoria.go in a new console window...
start cmd /k "cd ../memoria && go run memoria.go configs/config_PRUEBA_DEADLOCK.json"

echo Waiting for memoria.go to load (10 seconds)...
timeout /t 10 /nobreak >nul

echo Running cpu.go in a new console window...
start cmd /k "cd ../cpu && go run cpu.go configs/config_PRUEBA_DEADLOCK.json"

echo.
echo Script execution initiated.
pause