package main

import (
	"fmt"
	"os"

	"github.com/ozzhugoramirez/Evon/internal/cli"
)

func main() {
	args := os.Args[1:]

	// Si el usuario escribe solo "evon"
	if len(args) == 0 {
		if cli.IsEvonProject() {
			// Si detecta que estamos dentro de un proyecto, abre el REPL automáticamente
			cli.StartREPL(os.Stdin, os.Stdout)
		} else {
			// Si está en cualquier otra carpeta, muestra la ayuda
			cli.ShowGlobalHelp()
		}
		return
	}

	// Enrutador de comandos externos
	switch args[0] {
	case "new":
		cli.HandleNewProject(os.Stdout, args)
	case "init":
		cli.HandleInitProject(os.Stdout)
	case "shell":
		if cli.IsEvonProject() {
			cli.StartREPL(os.Stdin, os.Stdout)
		} else {
			fmt.Println("[Error] No estás dentro de un proyecto Evon. Usa 'evon init' o 'evon new <nombre>'.")
		}
	case "build":

		fmt.Println("[Build] Funcionalidad de Forja (LLVM/C++) en desarrollo.")
		fmt.Println("        Próximamente compilará el .axm a binario nativo.")
	case "--version", "-v":
		fmt.Println("Evon Engine v0.1.0-alpha (Core: Go | Backend: LLVM)")
	case "help", "--help", "-h":
		cli.ShowGlobalHelp()
	default:
		fmt.Printf("[Error] Comando desconocido: '%s'\n", args[0])
		cli.ShowGlobalHelp()
	}
}
