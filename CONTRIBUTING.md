

# Contribuir a Axiom Engine

¡Gracias por tu interés en mejorar Axiom! Al ser un motor arquitectónico, somos muy estrictos con la organización del código y el mantenimiento de la simplicidad.

## Reglas de Arquitectura (Standard Go Layout)

Axiom respeta la separación de responsabilidades absoluta. Si vas a enviar un Pull Request (PR), asegurate de que tu código respete esta estructura:

* **`cmd/axiom/`**: Es el punto de entrada. **Prohibido** agregar lógica de negocio acá. Solo debe inicializar el sistema.
* **`internal/engine/`**: Maneja el estado en la RAM y la persistencia en el archivo JSON. No interactúa con la terminal del usuario.
* **`internal/repl/`**: Exclusivamente para el ruteo de comandos (CLI) y la impresión de la consola. No debe contener lógica de transpilación ni escritura a disco.
* **`internal/editor/`**: Gestiona la entrada multilínea de texto.
* **`internal/transpiler/`**: Analiza el AST de Axiom y genera código binario/nativo.

## Guía de Estilo (Go)

1. **Usá `gofmt` siempre:** No aceptamos PRs con código mal indentado. Corré `go fmt ./...` antes de hacer un commit.
2. **Nombres explícitos:** Evitá variables de una letra excepto en bucles cortos. Preferimos `contextoActual` antes que `ctx`.
3. **Manejo de Errores:** No ignores los errores con `_`. En Go, los errores se manejan y se reportan de forma transparente.
4. **Cero dependencias innecesarias:** Axiom busca ser liviano. Usá la librería estándar de Go (`std`) tanto como sea posible. Consultá antes de agregar un paquete externo en `go.mod`.

## ¿Cómo enviar tus cambios?

1. Hacé un fork del repositorio.
2. Creá una rama descriptiva: `git checkout -b feature/nuevo-comando`.
3. Escribí tu código respetando los pilares de arriba.
4. Subí tus cambios y abrí un Pull Request explicando claramente qué problema resuelve tu código.