@echo off

echo Running kernel.go in a new console window...
start cmd /k "cd ../kernel && go run kernel.go configs/config_PRUEBA_SALVATIONS_EDGE.json"

echo Waiting for kernel.go to load (10 seconds)...
timeout /t 10 /nobreak >nul

echo Running entradasalida.go with GENERICA.json in a new console window...
start cmd /k "cd ../entradasalida && go run entradasalida.go configs/GENERICA.json"

echo Running entradasalida.go with SLP1.json in a new console window...
start cmd /k "cd ../entradasalida && go run entradasalida.go configs_salvations_edge/SLP1.json"

echo Running entradasalida.go with ESPERA.json in a new console window...
start cmd /k "cd ../entradasalida && go run entradasalida.go configs_salvations_edge/ESPERA.json"

echo Running entradasalida.go with TECLADO.json in a new console window...
start cmd /k "cd ../entradasalida && go run entradasalida.go configs_salvations_edge/TECLADO.json"

echo Running entradasalida.go with MONITOR.json in a new console window...
start cmd /k "cd ../entradasalida && go run entradasalida.go configs_salvations_edge/MONITOR.json"

echo Running memoria.go in a new console window...
start cmd /k "cd ../memoria && go run memoria.go configs/config_PRUEBA_SALVATIONS_EDGE.json"

echo Waiting for memoria.go to load (10 seconds)...
timeout /t 10 /nobreak >nul

echo Running cpu.go in a new console window...
start cmd /k "cd ../cpu && go run cpu.go configs/config_PRUEBA_SALVATIONS_EDGE.json"

echo.
echo Script execution initiated.
pause
