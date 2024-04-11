# 1º Checkpoint

## Ramas

Salen desde `main`:
- `conexion-inicial/kernel` (marcos, santiago)
- `conexion-inicial/cpu` (juan)
- `conexion-inicial/memoria` (elías)
- `conexion-inicial/entradasalida` (matías)

## Estructura del proyecto

Para cada módulo (`kernel`/`cpu`/`memoria`/`entradasalida`):
- _globals_: structs y variables globales de cada módulo
- _handlers_: funciones que se realizan al recibir un mensaje específico (handles the http request)
- _utils_: struct y funciones globales usadas por los 4 módulos
- _utils/logs/logs.go_: archivo dedicado al manejo y configuración de logs

## Requerimientos

- APIs (mensajes/endpoints de la interfaz)
  - `kernel`: definidas por enunciado
  - demás módulos: al menos un endpoint que retorne un _status 200_
- Logs: mínimos definidos por enunciado (en _handlers_)
- Configs: mínimas properties definidas por enunciado en un archivo json (en _globals_)

## Server

El archivo `{modulo}.go` actúa como server (se configura, expone una interfaz y se inicia)

```go
func main() {
	// =============
	// Configuración
	// =============
	logs.ConfigurarLogger("{ruta-log}") // ruta a raíz del módulo

	globals.Config = IniciarConfiguracion("{ruta-config}") // ruta a raíz del módulo
	if globals.Config == nil {
		log.Fatalln("Error al cargar la configuración")
	}

	// ========
	// Interfaz
	// ========
	mux := http.NewServeMux()
	mux.HandleFunc("/process", handlers.IniciarProceso)
	// ... demás mensajes

	// ======
	// Inicio
	// ======
	err := http.ListenAndServe(fmt.Sprintf(":%d", globals.Config.Port), mux)
	if err != nil {
		panic(err)
	}

	log.Printf("El módulo {modulo} está a la escucha en el puerto %d", globals.Config.Port)
}
```
