package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/handlers"
	configs "github.com/sisoputnfrba/tp-golang/utils/config"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	//// Configuración ////
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	globals.MConfig = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.Config{}).(*globals.Config)

	if globals.MConfig == nil {
		slog.Debug(fmt.Sprintf("Error al cargar la configuración"))
	}

	//// Logger ////
	configs.ConfigurarLogger("memoria", globals.MConfig.LogLevel)

	//// Inicialización ////
	globals.InicializarMemoriaUsuario()
	globals.InicializarMemoria()

	//// Conexión ////
	mux := http.NewServeMux()
	mux.HandleFunc("/contexto_de_ejecucion", handlers.HandleDevolverContexto)
	mux.HandleFunc("/actualizar_contexto", handlers.HandleActualizarContexto)
	mux.HandleFunc("/obtener_instruccion", handlers.HandleEnviarInstruccion)
	mux.HandleFunc("/read_mem", handlers.HandleReadMemory)
	mux.HandleFunc("/write_mem", handlers.HandleWriteMemory)
	mux.HandleFunc("/crear_proceso", handlers.HandleCrearProceso)
	mux.HandleFunc("/finalizar_proceso", handlers.HandleFinalizarProceso)
	mux.HandleFunc("/crear_hilo", handlers.HandleCrearHilo)
	mux.HandleFunc("/finalizar_hilo", handlers.HandleFinalizarHilo)
	mux.HandleFunc("/memory_dump", handlers.HandleMemoryDump)

	port := fmt.Sprintf(":%d", globals.MConfig.Port)

	slog.Info(fmt.Sprintf("El módulo memoria está a la escucha en el puerto %s", port))

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}
