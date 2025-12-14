.PHONY: run build swagger clean

# Variables
MAIN_FILE=cmd/api/main.go

# Generar documentación Swagger usando go run
swagger:
	@echo "Generando documentación Swagger..."
	go run github.com/swaggo/swag/cmd/swag@latest init -g $(MAIN_FILE) -o docs
	@echo "Documentación generada en docs/"

# Ejecutar servidor
run: swagger
	@echo "Iniciando servidor..."
	go run $(MAIN_FILE)

# Limpiar
clean:
	@echo "Limpiando..."
	go clean
	if exist docs rmdir /s /q docs

# Ayuda
help:
	@echo "Comandos disponibles:"
	@echo "  make swagger  - Generar documentación Swagger"
	@echo "  make run      - Ejecutar servidor"
	@echo "  make clean    - Limpiar archivos generados"