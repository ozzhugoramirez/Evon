.PHONY: run build clean install-path

# Variables
EVON_DIR = $(HOME)/.evon
BIN_DIR = $(EVON_DIR)/bin
STD_LIB = $(EVON_DIR)/lib/std

run:
	go run ./cmd/Evon/main.go

build:
	@echo "🔨 Compilando el motor Evon..."
	go build -o evon ./cmd/Evon/main.go
	@echo "📂 Preparando el entorno oficial en $(EVON_DIR)..."
	mkdir -p $(BIN_DIR)
	mkdir -p $(STD_LIB)
	@echo "🚀 Inyectando el nuevo binario en el sistema..."
	mv evon $(BIN_DIR)/
	@echo "Evon actualizado con éxito. Ya puedes usar el comando 'evon' en cualquier terminal."

clean:
	@echo "🧹 Limpiando binarios temporales..."
	rm -f evon
	rm -rf bin/

# Una regla de ayuda por si necesitas recordar cómo agregar Evon a tu PATH en el futuro
install-path:
	@echo "Para que 'evon' funcione globalmente, asegúrate de que tu PATH incluya ~/.evon/bin."
	@echo "Si usas Zsh (Mac/Linux moderno), corre este comando una sola vez:"
	@echo "echo 'export PATH=\"\$$HOME/.evon/bin:\$$PATH\"' >> ~/.zshrc && source ~/.zshrc"