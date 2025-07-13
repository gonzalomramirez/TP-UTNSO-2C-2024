package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sisoputnfrba/tp-golang/filesystem/globals"
	"github.com/sisoputnfrba/tp-golang/filesystem/handlers"
	"github.com/sisoputnfrba/tp-golang/filesystem/inicializacion"

	configs "github.com/sisoputnfrba/tp-golang/utils/config"
)

func main() {
	//// Configuracion ////
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	globals.FSConfig = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.Config{}).(*globals.Config)

	if globals.FSConfig == nil {
		slog.Debug(fmt.Sprintf("Error al cargar la configuración"))
	}

	//// Logger ////
	configs.ConfigurarLogger("filesystem", globals.FSConfig.LogLevel)

	//// Inicialización ////

	err = inicializacion.IniciarFS(globals.FSConfig.MountDir)

	if err != nil {
		log.Fatalf("Error al inicializar el File System: %v", err)
	}
	slog.Debug(fmt.Sprintf("Inicialización del File System completada."))

	//// Conexión ////
	mux := http.NewServeMux()
	mux.HandleFunc("/memory_dump", handlers.CrearArchivo)

	port := fmt.Sprintf(":%d", globals.FSConfig.Port)

	slog.Info(fmt.Sprintf("El módulo filesystem está a la escucha en el puerto %s", port))

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}
