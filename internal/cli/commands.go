package cli

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ozzhugoramirez/Evon/internal/compiler"
	"github.com/ozzhugoramirez/Evon/internal/memory"

	"golang.org/x/term"
)

// Códigos ANSI.
const (
	ColorRed   = "\033[31m"
	ColorWhite = "\033[1;37m"
	ColorDim   = "\033[2;37m"
	ColorBlue  = "\033[36m"
	ColorReset = "\033[0m"
)

func generateMockUUID(prefix string) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("axm_%s_%03d", prefix, rand.Intn(999))
}

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width < 40 {
		return 80
	}
	return width
}

// ---------------------------------------------------------
// GESTOR DE PROYECTOS Y CONTEXTO
// ---------------------------------------------------------

func IsEvonProject() bool {
	_, err := os.Stat(".evon/state.axm")
	return err == nil
}

func ShowGlobalHelp() {
	width := getTerminalWidth()
	divider := strings.Repeat("-", width)

	fmt.Fprintln(os.Stdout, ColorRed+divider+ColorReset)
	fmt.Fprintln(os.Stdout, ColorBlue+" EVON ENGINE (STRUCTURAL PROGRAMMING)"+ColorReset)
	fmt.Fprintln(os.Stdout, ColorRed+divider+ColorReset)

	fmt.Fprintln(os.Stdout, "\n  Navegación:")
	fmt.Fprintln(os.Stdout, ColorWhite+"  go /                         "+ColorDim+"Volver a la raíz."+ColorReset)
	fmt.Fprintln(os.Stdout, ColorWhite+"  go core                      "+ColorDim+"Ir a un nodo por nombre."+ColorReset)
	fmt.Fprintln(os.Stdout, ColorWhite+"  go core.entidades.usuario    "+ColorDim+"Ir a una ruta completa."+ColorReset)
	fmt.Fprintln(os.Stdout, ColorWhite+"  go usuario                   "+ColorDim+"Buscar dentro del contexto actual."+ColorReset)
	fmt.Fprintln(os.Stdout, ColorWhite+"  go @8f31                    "+ColorDim+"Buscar por prefijo de UUID."+ColorReset)
	fmt.Fprintln(os.Stdout, ColorWhite+"  go -                         "+ColorDim+"Volver al destino anterior."+ColorReset)

	fmt.Fprintln(os.Stdout, "\n  Comandos:")
	fmt.Fprintln(os.Stdout, ColorWhite+"  new <nombre>  "+ColorDim+" Crea un nuevo proyecto Evon."+ColorReset)
	fmt.Fprintln(os.Stdout, ColorWhite+"  init          "+ColorDim+" Inicializa Evon en el directorio actual."+ColorReset)
	fmt.Fprintln(os.Stdout, ColorWhite+"  shell         "+ColorDim+" Abre el entorno interactivo."+ColorReset)
	fmt.Fprintln(os.Stdout, ColorWhite+"  create entity <nombre> "+ColorDim+" Crea una entidad."+ColorReset)
	fmt.Fprintln(os.Stdout, ColorWhite+"  show          "+ColorDim+" Muestra la estructura."+ColorReset)
	fmt.Fprintln(os.Stdout, ColorWhite+"  run           "+ColorDim+" Evalúa la lógica del nodo actual."+ColorReset)
	fmt.Fprintln(os.Stdout, ColorWhite+"  edit          "+ColorDim+" Edita la lógica del nodo actual."+ColorReset)
	fmt.Fprintln(os.Stdout, ColorWhite+"  status        "+ColorDim+" Muestra el estado del proyecto."+ColorReset)
	fmt.Fprintln(os.Stdout, ColorWhite+"  log           "+ColorDim+" Muestra el historial."+ColorReset)
	fmt.Fprintln(os.Stdout, ColorWhite+"  undo          "+ColorDim+" Revierte el último evento del ledger."+ColorReset)
	fmt.Fprintln(os.Stdout, ColorWhite+"  exit          "+ColorDim+" Cierra Evon."+ColorReset)
}

func HandleNewProject(out io.Writer, args []string) {
	if len(args) < 2 {
		fmt.Fprintln(out, ColorRed+"[ERR] Falta el nombre del proyecto."+ColorReset)
		return
	}

	projectName := args[1]

	if _, err := os.Stat(projectName); !os.IsNotExist(err) {
		fmt.Fprintf(out, ColorRed+"[ERR] El directorio '%s' ya existe.\n"+ColorReset, projectName)
		return
	}

	fmt.Fprintf(out, ColorDim+"[SYS] Construyendo andamiaje para '%s'...\n"+ColorReset, projectName)

	if err := os.Mkdir(projectName, 0755); err != nil {
		fmt.Fprintf(out, ColorRed+"[ERR] No se pudo crear el directorio: %v\n"+ColorReset, err)
		return
	}

	setupEvonEnvironment(out, projectName)

	fmt.Fprintf(out, ColorWhite+"[OK] Proyecto '%s' inicializado. cd %s && evon\n"+ColorReset, projectName, projectName)
}

func HandleInitProject(out io.Writer) {
	if IsEvonProject() {
		fmt.Fprintln(out, ColorDim+"[SYS] Este directorio ya es un proyecto Evon."+ColorReset)
		return
	}

	setupEvonEnvironment(out, ".")
	fmt.Fprintln(out, ColorWhite+"[OK] Proyecto inicializado. Ejecuta 'evon' para entrar."+ColorReset)
}

func setupEvonEnvironment(out io.Writer, targetDir string) {
	evonDir := filepath.Join(targetDir, ".evon")

	if err := os.MkdirAll(evonDir, 0755); err != nil {
		fmt.Fprintf(out, ColorRed+"[ERR] No se pudo crear .evon: %v\n"+ColorReset, err)
		return
	}

	ledgerPath := filepath.Join(evonDir, "ledger.wal")
	if err := os.WriteFile(ledgerPath, []byte(""), 0644); err != nil {
		fmt.Fprintf(out, ColorRed+"[ERR] No se pudo crear ledger: %v\n"+ColorReset, err)
		return
	}

	storePath := filepath.Join(evonDir, "state.axm")
	store := memory.NewMemoryStore(storePath)
	store.Nodes = make(map[string]*memory.NodeRecord)

	store.Nodes["core"] = &memory.NodeRecord{
		UUID:     generateMockUUID("mod"),
		Type:     "MODULE",
		Name:     "core",
		Parent:   "root",
		Children: []string{"Ejemplo"},
		Props:    make(map[string]string),
	}

	store.Nodes["Ejemplo"] = &memory.NodeRecord{
		UUID:     generateMockUUID("ent"),
		Type:     "ENTITY",
		Name:     "Ejemplo",
		Parent:   "core",
		Children: []string{},
		Props:    map[string]string{"mensaje": "String"},
	}

	store.Nodes["Inicio"] = &memory.NodeRecord{
		UUID:     generateMockUUID("flw"),
		Type:     "FLOW",
		Name:     "Inicio",
		Parent:   "root",
		Children: []string{},
		Props:    map[string]string{"raw_logic": `puts("¡Hola desde Evon Engine!")`},
	}

	if err := store.Save(); err != nil {
		fmt.Fprintf(out, ColorRed+"[ERR] No se pudo guardar el estado inicial: %v\n"+ColorReset, err)
		return
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = targetDir
	_ = cmd.Run()

	gitignore := "# Evon Files\n.evon/\nhistory.tmp\nholograma.evon\n"
	if err := os.WriteFile(filepath.Join(targetDir, ".gitignore"), []byte(gitignore), 0644); err != nil {
		fmt.Fprintf(out, ColorRed+"[WARN] No se pudo crear .gitignore: %v\n"+ColorReset, err)
	}
}

// ---------------------------------------------------------
// ACCESO AL NODO ACTUAL
// ---------------------------------------------------------

func GetCurrentNode(state *REPLState) *memory.NodeRecord {
	if state == nil || state.Store == nil || len(state.Path) == 0 {
		return nil
	}

	return state.Store.Nodes[state.Path[len(state.Path)-1]]
}

// ---------------------------------------------------------
// CREACIÓN DE NODOS
// ---------------------------------------------------------

func HandleCreate(out io.Writer, state *REPLState, parts []string) {
	if len(parts) < 3 {
		fmt.Fprintln(out, ColorRed+"[ERR] Sintaxis: new [module|entity|law|flow] [Nombre]"+ColorReset)
		return
	}

	nodeType := strings.ToLower(parts[1])
	nodeName := parts[2]

	validTypes := map[string]bool{
		"module": true,
		"entity": true,
		"law":    true,
		"flow":   true,
	}

	if !validTypes[nodeType] {
		fmt.Fprintf(out, ColorRed+"[ERR] Tipo de nodo desconocido: '%s'.\n"+ColorReset, nodeType)
		return
	}

	if state.Store.Nodes == nil {
		state.Store.Nodes = make(map[string]*memory.NodeRecord)
	}

	if _, exists := state.Store.Nodes[nodeName]; exists {
		fmt.Fprintf(out, ColorRed+"[IMMUNE] Violación Topológica: El nodo '%s' ya existe.\n"+ColorReset, nodeName)
		return
	}

	parentName := "root"
	if len(state.Path) > 0 {
		parentName = state.Path[len(state.Path)-1]
	}

	if parentName != "root" && nodeType == "module" {
		parentNode := GetCurrentNode(state)
		if parentNode != nil && parentNode.Type != "MODULE" {
			fmt.Fprintln(out, ColorRed+"[IMMUNE] No puedes crear un Módulo dentro de un Nodo no modular."+ColorReset)
			return
		}
	}

	uuid := generateMockUUID(nodeType[:3])

	newNode := &memory.NodeRecord{
		UUID:     uuid,
		Type:     strings.ToUpper(nodeType),
		Name:     nodeName,
		Parent:   parentName,
		Children: []string{},
		Props:    make(map[string]string),
	}

	state.Store.Nodes[nodeName] = newNode

	if parentNode := GetCurrentNode(state); parentNode != nil {
		parentNode.Children = append(parentNode.Children, nodeName)
	}

	if err := state.Store.Save(); err != nil {
		delete(state.Store.Nodes, nodeName)
		fmt.Fprintf(out, ColorRed+"[ERR] No se pudo guardar el nodo: %v\n"+ColorReset, err)
		return
	}

	fmt.Fprintf(out, ColorDim+"[MEM] Nodo [%s] '%s' materializado. UUID: %s\n"+ColorReset, newNode.Type, nodeName, uuid)
}

// ---------------------------------------------------------
// NAVEGACIÓN ESTRUCTURAL: GO
// ---------------------------------------------------------

func clonePath(path []string) []string {
	result := make([]string, len(path))
	copy(result, path)
	return result
}

func containsChild(parent *memory.NodeRecord, childName string) bool {
	if parent == nil {
		return false
	}

	for _, name := range parent.Children {
		if name == childName {
			return true
		}
	}

	return false
}

func isValidChild(state *REPLState, parentName, childName string) bool {
	child, exists := state.Store.Nodes[childName]
	if !exists || child == nil || child.Parent != parentName {
		return false
	}

	parent, exists := state.Store.Nodes[parentName]
	if !exists || parent == nil {
		return false
	}

	return containsChild(parent, childName)
}

// Construye la ruta real de un nodo siguiendo sus relaciones Parent.
// También verifica que la relación Children sea consistente.
func buildPathToNode(state *REPLState, targetName string) ([]string, error) {
	if state == nil || state.Store == nil {
		return nil, fmt.Errorf("el almacén de memoria no está disponible")
	}

	target, exists := state.Store.Nodes[targetName]
	if !exists || target == nil {
		return nil, fmt.Errorf("el nodo '%s' no existe", targetName)
	}

	reversed := []string{}
	visited := make(map[string]bool)

	current := target

	for {
		if current == nil {
			return nil, fmt.Errorf("se encontró una referencia nula en la topología")
		}

		if visited[current.Name] {
			return nil, fmt.Errorf("se detectó un ciclo en la topología")
		}
		visited[current.Name] = true

		reversed = append(reversed, current.Name)

		if current.Parent == "root" {
			break
		}

		parent, parentExists := state.Store.Nodes[current.Parent]
		if !parentExists || parent == nil {
			return nil, fmt.Errorf("el padre '%s' de '%s' no existe", current.Parent, current.Name)
		}

		if !containsChild(parent, current.Name) {
			return nil, fmt.Errorf("la relación entre '%s' y '%s' no es consistente", parent.Name, current.Name)
		}

		current = parent
	}

	path := make([]string, len(reversed))
	for i := range reversed {
		path[len(reversed)-1-i] = reversed[i]
	}

	return path, nil
}

// Resuelve una ruta completa desde la raíz.
// Ejemplo: core.entidades.usuario
func resolveAbsolutePath(state *REPLState, route string) ([]string, error) {
	segments := strings.Split(route, ".")

	if len(segments) == 0 {
		return nil, fmt.Errorf("la ruta está vacía")
	}

	for i := range segments {
		segments[i] = strings.TrimSpace(segments[i])

		if segments[i] == "" {
			return nil, fmt.Errorf("la ruta contiene un segmento vacío")
		}
	}

	first := segments[0]
	firstNode, exists := state.Store.Nodes[first]

	if !exists || firstNode == nil {
		return nil, fmt.Errorf("el nodo raíz '%s' no existe", first)
	}

	if firstNode.Parent != "root" {
		return nil, fmt.Errorf("'%s' no es un nodo de la raíz", first)
	}

	path := []string{first}

	for i := 1; i < len(segments); i++ {
		parentName := segments[i-1]
		childName := segments[i]

		if !isValidChild(state, parentName, childName) {
			return nil, fmt.Errorf("'%s' no pertenece directamente a '%s'", childName, parentName)
		}

		path = append(path, childName)
	}

	return path, nil
}

// Busca recursivamente dentro del contexto actual.
// En la raíz, busca dentro de todos los árboles del proyecto.
func findContextualMatches(state *REPLState, targetName string) []string {
	if state == nil || state.Store == nil {
		return nil
	}

	matches := []string{}
	visited := make(map[string]bool)

	var walk func(nodeName string)

	walk = func(nodeName string) {
		if visited[nodeName] {
			return
		}
		visited[nodeName] = true

		node, exists := state.Store.Nodes[nodeName]
		if !exists || node == nil {
			return
		}

		if strings.EqualFold(node.Name, targetName) {
			matches = append(matches, node.Name)
		}

		for _, childName := range node.Children {
			if isValidChild(state, nodeName, childName) {
				walk(childName)
			}
		}
	}

	if len(state.Path) == 0 {
		// En raíz: recorrer todos los nodos que cuelgan de root.
		for name, node := range state.Store.Nodes {
			if node != nil && node.Parent == "root" {
				walk(name)
			}
		}
	} else {
		// Dentro de un nodo: recorrer sus descendientes.
		currentName := state.Path[len(state.Path)-1]
		current := GetCurrentNode(state)

		if current == nil {
			return nil
		}

		// También permite encontrar el nodo actual.
		if strings.EqualFold(current.Name, targetName) {
			matches = append(matches, current.Name)
		}

		for _, childName := range current.Children {
			if isValidChild(state, currentName, childName) {
				walk(childName)
			}
		}
	}

	return matches
}

func findByUUIDPrefix(state *REPLState, prefix string) ([]string, error) {
	if prefix == "" {
		return nil, fmt.Errorf("debes indicar un prefijo UUID después de @")
	}

	matches := []string{}

	for _, node := range state.Store.Nodes {
		if node == nil {
			continue
		}

		if strings.HasPrefix(strings.ToLower(node.UUID), strings.ToLower(prefix)) {
			matches = append(matches, node.Name)
		}
	}

	return matches, nil
}

func printNavigationContext(out io.Writer, state *REPLState) {
	width := getTerminalWidth()
	divider := strings.Repeat("-", width)

	fmt.Fprintln(out, "")
	fmt.Fprintln(out, ColorBlue+divider+ColorReset)

	node := GetCurrentNode(state)

	if node == nil {
		fmt.Fprintln(out, ColorWhite+" CONTEXTO: ROOT"+ColorReset)
	} else {
		fmt.Fprintf(
			out,
			" "+ColorWhite+"CONTEXTO: %-16s"+ColorDim+" | "+ColorWhite+"TIPO: %-8s"+ColorDim+" | "+ColorWhite+"UUID: %s\n"+ColorReset,
			node.Name,
			node.Type,
			node.UUID,
		)
	}

	fmt.Fprintln(out, ColorBlue+divider+ColorReset)
}

// Cambia la ubicación actual y conserva la anterior.
// No registra una ubicación anterior si el destino es el mismo.
func navigateTo(state *REPLState, newPath []string) {
	if strings.Join(state.Path, "\x00") == strings.Join(newPath, "\x00") {
		return
	}

	state.PreviousPath = clonePath(state.Path)
	state.Path = clonePath(newPath)
}

func HandleGo(out io.Writer, state *REPLState, target string) {
	if state == nil || state.Store == nil {
		fmt.Fprintln(out, ColorRed+"[ERR] El estado de Evon no está disponible."+ColorReset)
		return
	}

	if state.Store.Nodes == nil {
		fmt.Fprintln(out, ColorRed+"[SYS] Bóveda vacía."+ColorReset)
		return
	}

	target = strings.TrimSpace(target)

	if target == "" {
		fmt.Fprintln(out, ColorRed+"[ERR] Especifica un destino. Ejemplo: go core"+ColorReset)
		return
	}

	// go / -> raíz
	if target == "/" {
		navigateTo(state, []string{})
		printNavigationContext(out, state)
		return
	}

	// go - -> destino anterior
	if target == "-" {
		if state.PreviousPath == nil {
			fmt.Fprintln(out, ColorRed+"[ERR] No existe una ubicación anterior."+ColorReset)
			return
		}

		oldCurrent := clonePath(state.Path)
		state.Path = clonePath(state.PreviousPath)
		state.PreviousPath = oldCurrent

		printNavigationContext(out, state)
		return
	}

	// go @uuid -> búsqueda por prefijo de UUID
	if strings.HasPrefix(target, "@") {
		matches, err := findByUUIDPrefix(state, strings.TrimPrefix(target, "@"))
		if err != nil {
			fmt.Fprintf(out, ColorRed+"[ERR] %v\n"+ColorReset, err)
			return
		}

		if len(matches) == 0 {
			fmt.Fprintf(out, ColorRed+"[ERR] No existe ningún nodo con UUID que comience con '%s'.\n"+ColorReset, target)
			return
		}

		if len(matches) > 1 {
			fmt.Fprintf(out, ColorRed+"[ERR] El prefijo '%s' es ambiguo. Coincidencias:\n"+ColorReset, target)
			for _, name := range matches {
				node := state.Store.Nodes[name]
				fmt.Fprintf(out, "  - %s (%s)\n", name, node.UUID)
			}
			return
		}

		path, err := buildPathToNode(state, matches[0])
		if err != nil {
			fmt.Fprintf(out, ColorRed+"[ERR] No se pudo construir la ruta: %v\n"+ColorReset, err)
			return
		}

		navigateTo(state, path)
		printNavigationContext(out, state)
		return
	}

	// Una ruta con puntos siempre se interpreta desde la raíz.
	if strings.Contains(target, ".") {
		path, err := resolveAbsolutePath(state, target)
		if err != nil {
			fmt.Fprintf(out, ColorRed+"[ERR] Ruta inválida: %v\n"+ColorReset, err)
			return
		}

		navigateTo(state, path)
		printNavigationContext(out, state)
		return
	}

	// Primero busca dentro del contexto actual.
	contextualMatches := findContextualMatches(state, target)

	if len(contextualMatches) > 1 {
		fmt.Fprintf(out, ColorRed+"[ERR] El nombre '%s' es ambiguo dentro del contexto actual:\n"+ColorReset, target)

		for _, name := range contextualMatches {
			node := state.Store.Nodes[name]
			path, _ := buildPathToNode(state, name)

			fmt.Fprintf(out, "  - %s | /%s | %s\n", name, strings.Join(path, "/"), node.UUID)
		}
		return
	}

	if len(contextualMatches) == 1 {
		path, err := buildPathToNode(state, contextualMatches[0])
		if err != nil {
			fmt.Fprintf(out, ColorRed+"[ERR] No se pudo construir la ruta: %v\n"+ColorReset, err)
			return
		}

		navigateTo(state, path)
		printNavigationContext(out, state)
		return
	}

	// Si no está en el contexto, busca por nombre global.
	var globalMatch *memory.NodeRecord

	for _, node := range state.Store.Nodes {
		if node != nil && strings.EqualFold(node.Name, target) {
			if globalMatch != nil {
				fmt.Fprintf(out, ColorRed+"[ERR] Hay más de un nodo llamado '%s'. Usa su UUID o ruta completa.\n"+ColorReset, target)
				return
			}
			globalMatch = node
		}
	}

	if globalMatch == nil {
		fmt.Fprintf(out, ColorRed+"[ERR] No se encontró el destino '%s'.\n"+ColorReset, target)
		return
	}

	path, err := buildPathToNode(state, globalMatch.Name)
	if err != nil {
		fmt.Fprintf(out, ColorRed+"[ERR] No se pudo construir la ruta: %v\n"+ColorReset, err)
		return
	}

	navigateTo(state, path)
	printNavigationContext(out, state)
}

// ---------------------------------------------------------
// VISUALIZACIÓN DE LA ESTRUCTURA
// ---------------------------------------------------------

func HandleShow(out io.Writer, state *REPLState) {
	width := getTerminalWidth()
	divider := strings.Repeat("-", width)

	fmt.Fprintln(out, ColorBlue+divider+ColorReset)
	fmt.Fprintln(out, " "+ColorRed+"SUBMISSION PLATFORM [AST]"+ColorReset)
	fmt.Fprintln(out, ColorBlue+divider+ColorReset)

	if state.Store.Nodes == nil || len(state.Store.Nodes) == 0 {
		fmt.Fprintln(out, ColorDim+"[SYS] No hay nodos en la bóveda."+ColorReset)
		return
	}

	for _, node := range state.Store.Nodes {
		if node == nil || node.Parent != "root" {
			continue
		}

		fmt.Fprintf(out, ColorWhite+" %-14s %-10s %s\n"+ColorReset, node.Name, "<"+node.Type+">", node.UUID)

		for _, childName := range node.Children {
			if child, exists := state.Store.Nodes[childName]; exists && child != nil {
				fmt.Fprintf(out, ColorDim+"   %-12s %-10s %s\n"+ColorReset, child.Name, "<"+child.Type+">", child.UUID)

				for _, grandchildName := range child.Children {
					if grandchild, ok := state.Store.Nodes[grandchildName]; ok && grandchild != nil {
						fmt.Fprintf(out, ColorDim+"      %-12s %-10s %s\n"+ColorReset, grandchild.Name, "<"+grandchild.Type+">", grandchild.UUID)
					}
				}
			}
		}
	}

	fmt.Fprintln(out, ColorBlue+divider+ColorReset)
}

// ---------------------------------------------------------
// EJECUCIÓN Y EVALUACIÓN DEL AST (RUN / SIM)
// ---------------------------------------------------------

func HandleRunSim(out io.Writer, state *REPLState, mode string) {
	currentNode := GetCurrentNode(state)

	if currentNode == nil {
		fmt.Fprintln(out, ColorRed+"[ERR] Estás en Root. Ingresá a un nodo (ej: go Inicio) para evaluar su lógica."+ColorReset)
		return
	}

	rawCode, exists := currentNode.Props["raw_logic"]
	if !exists || strings.TrimSpace(rawCode) == "" {
		fmt.Fprintf(out, ColorDim+"[SYS] El nodo '%s' está vacío. Nada que ejecutar.\n"+ColorReset, currentNode.Name)
		return
	}

	width := getTerminalWidth()
	divider := strings.Repeat("-", width)

	fmt.Fprintln(out, ColorBlue+divider+ColorReset)
	fmt.Fprintf(out, " "+ColorWhite+"EJECUTANDO NODO: %s"+ColorDim+" | MODO: %s\n"+ColorReset, currentNode.Name, strings.ToUpper(mode))
	fmt.Fprintln(out, ColorBlue+divider+ColorReset)

	// Análisis léxico.
	fmt.Fprintln(out, ColorDim+"--- [ ANÁLISIS LÉXICO (TOKENS) ] ---"+ColorReset)

	debugLexer := compiler.NewLexer(rawCode)

	for {
		tok := debugLexer.NextToken()
		if tok.Type == compiler.EOF {
			break
		}

		fmt.Fprintf(out, ColorDim+"Token: %-10s | Literal: '%s'\n"+ColorReset, tok.Type, tok.Literal)
	}

	fmt.Fprintln(out, "")

	// Parseo estructural.
	lexer := compiler.NewLexer(rawCode)
	parser := compiler.NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		fmt.Fprintln(out, ColorRed+"[SISTEMA INMUNE] Fallo estructural previo a ejecución:"+ColorReset)

		for _, err := range parser.Errors() {
			fmt.Fprintf(out, ColorWhite+"  ! %s\n"+ColorReset, err)
		}

		return
	}

	// Ejecución del AST.
	fmt.Fprintln(out, ColorDim+"--- [ RESULTADO DE EVALUACIÓN ] ---"+ColorReset)
	compiler.Eval(program)

	fmt.Fprintln(out, ColorBlue+divider+ColorReset)
}

// ---------------------------------------------------------
// INYECCIÓN DIRECTA DE CÓDIGO
// ---------------------------------------------------------

func HandleDirectInjection(out io.Writer, state *REPLState, rawLine string) {
	currentNode := GetCurrentNode(state)

	if currentNode == nil {
		fmt.Fprintln(out, ColorRed+"[ERR] Estás en Root. Ingresá a un nodo con 'go' para inyectar lógica."+ColorReset)
		return
	}

	oldLogic, _ := currentNode.Props["raw_logic"]

	newLogic := rawLine
	if oldLogic != "" {
		newLogic = oldLogic + "\n" + rawLine
	}

	lexer := compiler.NewLexer(newLogic)
	parser := compiler.NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		width := getTerminalWidth()
		divider := strings.Repeat("-", width)

		fmt.Fprintln(out, ColorRed+divider+ColorReset)
		fmt.Fprintln(out, " "+ColorWhite+"[ SISTEMA INMUNE: RECHAZO DE INYECCIÓN DIRECTA ]"+ColorReset)
		fmt.Fprintln(out, ColorRed+divider+ColorReset)

		for _, err := range parser.Errors() {
			fmt.Fprintf(out, ColorWhite+"  ! %s\n"+ColorReset, err)
		}

		fmt.Fprintln(out, ColorRed+divider+ColorReset)
		return
	}

	currentNode.Props["raw_logic"] = newLogic
	currentNode.Props["statements"] = fmt.Sprintf("%d", len(program.Statements))

	if err := state.Store.Save(); err != nil {
		fmt.Fprintf(out, ColorRed+"[ERR] No se pudo persistir la lógica: %v\n"+ColorReset, err)
		return
	}

	payload := fmt.Sprintf("Inyección Directa en '%s'. Oraciones: %d.", currentNode.Name, len(program.Statements))
	_, _ = state.Ledger.Record("INJECT", payload)

	fmt.Fprintf(out, ColorDim+"[INYECTADO] Lógica añadida a %s.\n"+ColorReset, currentNode.Name)
}
