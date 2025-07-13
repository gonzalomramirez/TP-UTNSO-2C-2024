package globals

import (
	"errors"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log/slog"
	"sync"
)

type Config struct {
	Port            int    `json:"port"`
	MemorySize      int    `json:"memory_size"`
	InstructionPath string `json:"instruction_path"`
	ResponseDelay   int    `json:"response_delay"`
	IpKernel        string `json:"ip_kernel"`
	PortKernel      int    `json:"port_kernel"`
	IpCpu           string `json:"ip_cpu"`
	PortCpu         int    `json:"port_cpu"`
	IpFileSystem    string `json:"ip_filesystem"`
	PortFileSystem  int    `json:"port_filesystem"`
	Scheme          string `json:"scheme"`
	SearchAlgorithm string `json:"search_algorithm"`
	Partitions      []int  `json:"partitions"`
	LogLevel        string `json:"log_level"`
}

var MConfig *Config

type ContextoProceso struct {
	Base   int // Dirección base de la memoria del proceso
	Limite int // Límite (tamaño de memoria asignada)
}

type InstruccionesHilo struct {
	Instrucciones []string // Instrucciones leídas del pseudocódigo
}

type MemSistema struct {
	TablaProcesos map[int]*ContextoProceso           // Tabla de procesos (PID -> Contexto de proceso)
	TablaHilos    map[int]map[int]*commons.Registros // Tabla de hilos (PID -> TID -> Contexto de hilo)
	Pseudocodigos map[int]map[int]*InstruccionesHilo // Pseudocódigos (PID -> TID -> Código)
}

var MemoriaSistema MemSistema

type MemUsuario struct {
	Datos       []byte
	Particiones []*Particion
	MutexDatos  sync.Mutex
}

var MemoriaUsuario MemUsuario

func InicializarMemoriaUsuario() {
	if MConfig == nil {
		panic("MConfig no está inicializado")
	}

	MemoriaUsuario = MemUsuario{
		Datos:       make([]byte, MConfig.MemorySize),
		Particiones: []*Particion{},
	}
}

type Particion struct {
	Base   int
	Limite int
	Libre  bool
	Pid    int
}

func InicializarMemoria() {
	MemoriaSistema = MemSistema{
		TablaProcesos: make(map[int]*ContextoProceso),
		TablaHilos:    make(map[int]map[int]*commons.Registros),
		Pseudocodigos: make(map[int]map[int]*InstruccionesHilo),
	}

	if MConfig.Scheme == "FIJAS" {
		base := 0
		for _, tamanio := range MConfig.Partitions {
			if base+tamanio > MConfig.MemorySize {
				_ = errors.New("error: Particiones fijas exceden el tamaño total de memoria")
			}

			// Crear una nueva partición y agregarla a la lista
			nuevaParticion := Particion{
				Base:   base,
				Limite: base + tamanio,
				Libre:  true,
				Pid:    -1,
			}
			MemoriaUsuario.Particiones = append(MemoriaUsuario.Particiones, &nuevaParticion)
			base += tamanio
		}

		if base != MConfig.MemorySize {
			fmt.Println("Advertencia: No se ha utilizado la memoria completa en particiones fijas.")
		}
	} else {
		nuevaParticion := Particion{
			Base:   0,
			Limite: MConfig.MemorySize - 1,
			Libre:  true,
			Pid:    -1,
		}
		MemoriaUsuario.Particiones = append(MemoriaUsuario.Particiones, &nuevaParticion)
	}

	slog.Debug(fmt.Sprintf("Memoria inicializada con éxito."))
}
