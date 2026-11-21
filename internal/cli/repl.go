package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ozzhugoramirez/Evon/internal/memory"

	"github.com/chzyer/readline"
)

type REPLState struct {
	Path         []string
	PreviousPath []string
	Directory    string
	GitBranch    string
	Store        *memory.MemoryStore
	Ledger       *memory.Ledger
}

// ---------------------------------------------------------
// PROMPT
// ---------------------------------------------------------

// GetPrompt devuelve únicamente el prompt de entrada.
// No debe contener saltos de línea ni el encabezado del proyecto,
// porque readline redibuja este texto mientras el usuario escribe.
func (s *REPLState) GetPrompt() string {
	return "> "
}

// GetHeader devuelve la información contextual que se imprime
// una sola vez antes de cada lectura.
func (s *REPLState) GetHeader() string {
	directory := s.Directory
	if directory == "" {
		directory = "evon"
	}

	branch := s.GitBranch
	if branch == "" {
		branch = "sin-git"
	}

	location := "/"
	if len(s.Path) > 0 {
		location += strings.Join(s.Path, "/")
	}

	return fmt.Sprintf(
		"%s%s (%s) | %s%s",
		ColorDim,
		directory,
		branch,
		location,
		ColorReset,
	)
}

// Obtiene el nombre del directorio del proyecto y su rama Git.
func getREPLLocation() (string, string) {
	cwd, err := os.Getwd()
	if err != nil {
		return "evon", "sin-git"
	}

	directory := filepath.Base(filepath.Clean(cwd))
	if directory == "." || directory == string(filepath.Separator) {
		directory = "evon"
	}

	branch := "sin-git"

	cmd := exec.Command("git", "-C", cwd, "branch", "--show-current")
	if output, err := cmd.Output(); err == nil {
		value := strings.TrimSpace(string(output))
		if value != "" {
			branch = value
		}
	}

	return directory, branch
}

// ---------------------------------------------------------
// INICIO DEL REPL
// ---------------------------------------------------------

func StartREPL(in io.Reader, out io.Writer) {
	store := memory.NewMemoryStore(".evon/state.axm")

	if err := store.LoadMmap(); err != nil {
		fmt.Fprintf(
			out,
			ColorRed+"[ERR] No se pudo cargar la memoria: %v\n"+ColorReset,
			err,
		)
		return
	}

	ledger, err := memory.NewLedger(".evon/ledger.wal")
	if err != nil {
		fmt.Fprintf(
			out,
			ColorRed+"[ERR] No se pudo cargar el ledger: %v\n"+ColorReset,
			err,
		)
		return
	}

	directory, branch := getREPLLocation()

	state := &REPLState{
		Path:         []string{},
		PreviousPath: nil,
		Directory:    directory,
		GitBranch:    branch,
		Store:        store,
		Ledger:       ledger,
	}

	fmt.Fprint(out, "\033[H\033[2J")

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          state.GetPrompt(),
		HistoryFile:     ".evon/history.tmp",
		AutoComplete:    nil,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})

	if err != nil {
		fmt.Fprintln(
			out,
			ColorRed+"[ERR] No se pudo inicializar la terminal."+ColorReset,
		)
		return
	}

	defer rl.Close()

	ShowGlobalHelp()

	for {
		// El encabezado se imprime fuera de readline.
		// Así no se vuelve a dibujar por cada tecla ingresada.
		fmt.Fprintln(out, state.GetHeader())

		// readline administra solamente esta línea.
		rl.SetPrompt(state.GetPrompt())

		line, err := rl.Readline()

		if err != nil {
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					break
				}
				continue
			}

			break
		}

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		cmd := strings.ToLower(parts[0])

		switch cmd {

		// -------------------------------------------------
		// SALIDA
		// -------------------------------------------------

		case "exit", "quit":
			fmt.Fprintln(
				out,
				ColorDim+"[SYS] Apagando motor Evon. Memoria consolidada."+ColorReset,
			)
			return

		// -------------------------------------------------
		// TERMINAL
		// -------------------------------------------------

		case "clear":
			fmt.Fprint(out, "\033[H\033[2J")

		// -------------------------------------------------
		// NAVEGACIÓN
		// -------------------------------------------------

		case "go":
			if len(parts) < 2 {
				HandleGo(out, state, "")
				continue
			}

			HandleGo(out, state, parts[1])

		// -------------------------------------------------
		// ESTADO
		// -------------------------------------------------

		case "status":
			fmt.Fprintln(out, ColorWhite+"[Status] Estado del Proyecto:"+ColorReset)
			fmt.Fprintf(
				out,
				"  -> Nodos Activos en Bóveda: %d\n",
				len(state.Store.Nodes),
			)
			fmt.Fprintf(
				out,
				"  -> Eventos en Ledger (WAL): %d\n",
				len(state.Ledger.GetEvents()),
			)

			current := "/"
			if len(state.Path) > 0 {
				current += strings.Join(state.Path, "/")
			}

			fmt.Fprintf(out, "  -> Contexto Actual: %s\n", current)

			previous := "No disponible"
			if state.PreviousPath != nil {
				previous = "/"
				if len(state.PreviousPath) > 0 {
					previous += strings.Join(state.PreviousPath, "/")
				}
			}

			fmt.Fprintf(out, "  -> Contexto Anterior: %s\n", previous)

		// -------------------------------------------------
		// HISTORIAL
		// -------------------------------------------------

		case "log":
			fmt.Fprintln(
				out,
				ColorWhite+"[Log] Historial de Eventos (.wal):"+ColorReset,
			)

			for _, ev := range state.Ledger.GetEvents() {
				fmt.Fprintf(
					out,
					"  [%02d] %s | %s -> %s (Hash: %s)\n",
					ev.ID,
					ev.Timestamp.Format("15:04:05"),
					ev.Action,
					ev.Payload,
					ev.Hash,
				)
			}

		// -------------------------------------------------
		// UNDO
		// -------------------------------------------------

		case "undo":
			ev, err := state.Ledger.Undo()

			if err != nil {
				fmt.Fprintf(out, ColorRed+"[ERR] %v\n"+ColorReset, err)
			} else {
				fmt.Fprintf(
					out,
					ColorWhite+"[Undo] Evento #%d [%s] revertido.\n"+ColorReset,
					ev.ID,
					ev.Action,
				)
			}

		// -------------------------------------------------
		// EDICIÓN
		// -------------------------------------------------

		case "edit":
			var target string

			if len(parts) > 1 {
				target = parts[1]
			} else if len(state.Path) > 0 {
				target = state.Path[len(state.Path)-1]
			} else {
				fmt.Fprintln(
					out,
					ColorRed+"[ERR] Estás en Root. Usa 'go <nodo>' o 'edit <nodo>'."+ColorReset,
				)
				continue
			}

			RunEditor(in, out, target, state.Store, state.Ledger)

		// -------------------------------------------------
		// CREACIÓN
		// -------------------------------------------------

		case "create", "new":
			HandleCreate(out, state, parts)

			if len(parts) > 1 {
				_, _ = state.Ledger.Record(
					"CREATE",
					strings.Join(parts[1:], " "),
				)
			}

		// -------------------------------------------------
		// VISUALIZACIÓN
		// -------------------------------------------------

		case "show":
			HandleShow(out, state)

		// -------------------------------------------------
		// EJECUCIÓN
		// -------------------------------------------------

		case "run", "sim":
			HandleRunSim(out, state, cmd)

		// -------------------------------------------------
		// AYUDA
		// -------------------------------------------------

		case "help":
			ShowGlobalHelp()

		// -------------------------------------------------
		// CÓDIGO ESTRUCTURAL
		// -------------------------------------------------

		default:
			if isStructuralCode(line) {
				HandleDirectInjection(out, state, line)
			} else {
				fmt.Fprintf(
					out,
					ColorRed+"[ERR] Comando desconocido: '%s'. Escribe 'help' para ver los comandos.\n"+ColorReset,
					cmd,
				)
			}
		}
	}
}

// ---------------------------------------------------------
// DETECCIÓN DE CÓDIGO ESTRUCTURAL
// ---------------------------------------------------------

func isStructuralCode(line string) bool {
	return strings.Contains(line, "=") ||
		strings.Contains(line, ":") ||
		strings.HasPrefix(line, "let ") ||
		strings.HasPrefix(line, "mut ") ||
		strings.HasPrefix(line, "deny ") ||
		strings.HasPrefix(line, "ensure ") ||
		strings.HasPrefix(line, "puts") ||
		strings.Contains(line, "(")
}
