package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ozzhugoramirez/Evon/internal/memory"
)

func RunEditor(
	in io.Reader,
	out io.Writer,
	nodeName string,
	store *memory.MemoryStore,
	ledger *memory.Ledger,
) {
	if store == nil {
		fmt.Fprintln(out, ColorRed+"[ERR] El almacén de memoria no está disponible."+ColorReset)
		return
	}

	node, exists := store.Nodes[nodeName]
	if !exists || node == nil {
		fmt.Fprintf(
			out,
			ColorRed+"[ERR] El nodo '%s' no existe."+ColorReset+"\n",
			nodeName,
		)
		return
	}

	tmpFile := filepath.Join(".evon", "holograma.evon")

	// Recuperar el código anterior.
	oldCode := ""
	if node.Props != nil {
		oldCode = node.Props["raw_logic"]
	}

	// Crear el archivo temporal con el contenido anterior.
	if err := os.WriteFile(tmpFile, []byte(oldCode), 0644); err != nil {
		fmt.Fprintf(
			out,
			ColorRed+"[ERR] No se pudo crear el archivo temporal: %v\n"+ColorReset,
			err,
		)
		return
	}

	// Seleccionar el editor del sistema.
	editor := os.Getenv("EDITOR")

	if editor == "" {
		if runtime.GOOS == "windows" {
			editor = "notepad"
		} else {
			editor = "vim"
		}
	}

	cmd := exec.Command(editor, tmpFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(
			out,
			ColorRed+"[ERR] El editor nativo falló: %v\n"+ColorReset,
			err,
		)
		_ = os.Remove(tmpFile)
		return
	}

	// Leer el contenido escrito.
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		fmt.Fprintf(
			out,
			ColorRed+"[ERR] No se pudo leer el código: %v\n"+ColorReset,
			err,
		)
		_ = os.Remove(tmpFile)
		return
	}

	_ = os.Remove(tmpFile)

	rawCode := string(data)

	// Un documento vacío no reemplaza el código anterior.
	if strings.TrimSpace(rawCode) == "" {
		fmt.Fprint(out, "\033[H\033[2J")
		fmt.Fprintln(
			out,
			ColorDim+"[SYS] Lienzo vacío. Código anterior conservado."+ColorReset,
		)
		return
	}

	// Garantizar que Props exista.
	if node.Props == nil {
		node.Props = make(map[string]string)
	}

	// Guardar solamente la lógica.
	// No reemplazamos el NodeRecord.
	previousCode, hadPreviousCode := node.Props["raw_logic"]
	node.Props["raw_logic"] = rawCode

	// Persistir el nodo completo conservando su estructura.
	if err := store.Save(); err != nil {
		// Restaurar el estado anterior si falla la persistencia.
		if hadPreviousCode {
			node.Props["raw_logic"] = previousCode
		} else {
			delete(node.Props, "raw_logic")
		}

		fmt.Fprintf(
			out,
			ColorRed+"[ERR] No se pudo guardar el código: %v\n"+ColorReset,
			err,
		)
		return
	}

	// Registrar la edición.
	if ledger != nil {
		payload := fmt.Sprintf(
			"Código almacenado en '%s'. Validación sintáctica pendiente.",
			nodeName,
		)

		if _, err := ledger.Record("EDIT", payload); err != nil {
			fmt.Fprintf(
				out,
				ColorRed+"[WARN] El código se guardó, pero no se pudo registrar en el ledger: %v\n"+ColorReset,
				err,
			)
		}
	}

	// Confirmación visual.
	width := getTerminalWidth()
	divider := strings.Repeat("-", width)

	fmt.Fprint(out, "\033[H\033[2J")
	fmt.Fprintln(out, ColorBlue+divider+ColorReset)
	fmt.Fprintln(out, ColorWhite+"[ CÓDIGO ALMACENADO ]"+ColorReset)
	fmt.Fprintf(out, " Nodo: %s\n", node.Name)
	fmt.Fprintf(out, " Tipo: %s\n", node.Type)
	fmt.Fprintln(out, ColorDim+" Estado: guardado sin validación sintáctica."+ColorReset)
	fmt.Fprintln(out, ColorDim+" El compilador evaluará el código cuando ejecutes run."+ColorReset)
	fmt.Fprintln(out, ColorBlue+divider+ColorReset)
}
