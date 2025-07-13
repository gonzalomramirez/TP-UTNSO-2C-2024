package processes

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/threads"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/request"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log/slog"
	"net/http"
	"strconv"
)

func ProcesoInicial(argumentos []string) {
	pseudocodigo := argumentos[1]
	tamanio, _ := strconv.Atoi(argumentos[2])
	prioridadHiloMain := 0

	CrearProceso(pseudocodigo, tamanio, prioridadHiloMain)

	return
}

func CrearProceso(pseudocodigo string, tamanioMemoria int, prioridad int) int {
	pcb := CrearPCB(pseudocodigo, tamanioMemoria, prioridad)

	slog.Info(fmt.Sprintf("## (%d:0) Se crea el proceso - Estado: NEW", pcb.Pid))

	if len(globals.Estructura.ColaNew) == 0 {
		slog.Debug(fmt.Sprintf("Cola NEW está vacía, solicitando memoria."))

		respuestaMemoria, err := SolicitarProcesoMemoria(pcb.Pid, pseudocodigo, tamanioMemoria)

		if err != nil {
			slog.Debug(fmt.Sprintf("Error al solicitar espacio en memoria."))
		}

		if respuestaMemoria.StatusCode == http.StatusOK {

			threads.CrearHilo(pcb.Pid, prioridad, pseudocodigo)

		} else {
			slog.Debug(fmt.Sprintf("Memoria no tiene suficiente espacio. Proceso en espera."))

			queues.AgregarProcesoACola(pcb, &globals.Estructura.ColaNew)

			return http.StatusBadRequest
		}
	} else {
		slog.Debug(fmt.Sprintf("Cola NEW no está vacía, proceso %d se encola en NEW.", pcb.Pid))

		queues.AgregarProcesoACola(pcb, &globals.Estructura.ColaNew)

		SolicitarProcesoMemoria(pcb.Pid, pseudocodigo, -1)

		return http.StatusBadRequest
	}
	return http.StatusOK
}

func CrearPCB(pseudocodigo string, tamanio int, prioridad int) *commons.PCB {
	pcb := commons.PCB{
		Pid:               globals.Estructura.ContadorPid,
		Estado:            "NEW",
		Tid:               []*commons.TCB{},
		ContadorHilos:     0,
		Tamanio:           tamanio,
		PseudoCodigoHilo0: pseudocodigo,
		PrioridadTID0:     prioridad,
		Mutex:             []commons.Mutex{},
	}

	globals.Estructura.ContadorPid++

	globals.Estructura.Procesos[pcb.Pid] = &pcb

	return &pcb
}

func FinalizarProceso(pid int) {
	req := request.RequestFinalizarProceso{
		Pid: pid,
	}

	solicitudCodificada, err := commons.CodificarJSON(req)

	if err != nil {
		slog.Debug(fmt.Sprintf("Error al codificar la solicitud de finalización de proceso"))
		return
	}

	pcb := queues.BuscarPCBEnColas(pid)

	for _, tcb := range pcb.Tid {
		threads.FinalizarHilo(pid, tcb.Tid)
	}

	response := cliente.Post(globals.KConfig.IpMemory, globals.KConfig.PortMemory, "finalizar_proceso", solicitudCodificada)

	if response.StatusCode == http.StatusOK {
		defer delete(globals.Estructura.Procesos, pid)

		pcb.Estado = "EXIT"

		slog.Info(fmt.Sprintf("## Finaliza el proceso %d", pid))

		for len(globals.Estructura.ColaNew) != 0 {
			procesoNuevo := globals.Estructura.ColaNew[0]
			response, _ := SolicitarProcesoMemoria(procesoNuevo.Pid, procesoNuevo.PseudoCodigoHilo0, procesoNuevo.Tamanio)
			if response.StatusCode == http.StatusOK {
				queues.SacarProcesoDeCola(procesoNuevo.Pid, &globals.Estructura.ColaNew)
				threads.CrearHilo(procesoNuevo.Pid, procesoNuevo.PrioridadTID0, procesoNuevo.PseudoCodigoHilo0)
			} else {
				break
			}
		}
	} else {
		slog.Debug(fmt.Sprintf("## Error al finalizar el proceso %d", pid))
	}
}

func SolicitarProcesoMemoria(pid int, pseudocodigo string, tamanio int) (*http.Response, error) {
	req := request.RequestProcessCreateMemoria{
		Pid:            pid,
		Pseudocodigo:   pseudocodigo,
		TamanioMemoria: tamanio,
	}

	solicitudCodificada, err := commons.CodificarJSON(req)

	if err != nil {
		return nil, err
	}

	return cliente.Post(globals.KConfig.IpMemory, globals.KConfig.PortMemory, "crear_proceso", solicitudCodificada), nil
}
