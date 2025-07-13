# GOD (Go Direction) - Simulación de un Sistema Operativo

Este proyecto fue desarrollado como trabajo práctico cuatrimestral para la materia **Sistemas Operativos** (UTN FRBA, 2° cuatrimestre 2024).

El objetivo principal fue simular un sistema operativo distribuido compuesto por distintos módulos comunicados vía HTTP:

- **Kernel**: planificación y gestión de procesos e hilos  
- **Memoria**: asignación de memoria con particiones fijas/dinámicas  
- ⚙**CPU**: ciclo de instrucción con registros, interrupciones y MMU  
- **FileSystem**: sistema de archivos con asignación indexada y persistencia en disco

El proyecto incluye soporte para planificación por FIFO, Prioridades y Colas Multinivel, manejo de mutex, syscalls, dump de memoria y testing distribuido.

## Tecnologías utilizadas
- Golang
- APIs RESTful
- Linux CLI
- Git
- Testing distribuido
- Logs estructurados con `slog`

## Estructura general del proyecto
```
/kernel
/memoria
/cpu
/filesystem
/scripts
```

## Cómo ejecutar
1. Configurar los archivos `.config` en cada módulo.
2. Ejecutar los módulos con `go run` o `make`.
3. Usar los scripts de prueba para simular ejecuciones distribuidas.

> Implementado y probado en entorno Linux. Compatible con múltiples computadoras.

## Créditos
Trabajo realizado por equipo de estudiantes de UTN FRBA como parte del TP final de la materia Sistemas Operativos (2C2024).
