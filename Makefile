# Variables
BINARY_NAME=biasharatrack
BACKEND_DIR=backend
FRONTEND_DIR=.

# Go commands
GOBUILD=go build
GORUN=go run
GOCLEAN=go clean

# NPM commands
NPM=npm

# Ngrok command
NGROK=ngrok

# Default make target
.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build-backend   - Build the Go backend"
	@echo "  build-frontend  - Build the frontend"
	@echo "  build          - Build both backend and frontend"
	@echo "  run-backend    - Run the Go backend"
	@echo "  run-frontend   - Run the frontend development server"
	@echo "  run            - Run both backend and frontend"
	@echo "  clean          - Clean up build artifacts"
	@echo "  install        - Install all dependencies"
	@echo "  start-ngrok    - Start ngrok and update callback URL"

.PHONY: build-backend
build-backend:
	cd $(BACKEND_DIR) && $(GOBUILD) -o $(BINARY_NAME)

.PHONY: build-frontend
build-frontend:
	cd $(FRONTEND_DIR) && $(NPM) run build

.PHONY: build
build: build-backend build-frontend

.PHONY: run-backend
run-backend:
	cd $(BACKEND_DIR) && $(GORUN) main.go

.PHONY: run-frontend
run-frontend:
	cd $(FRONTEND_DIR) && $(NPM) run dev

.PHONY: run
run:
	make -j 2 start-ngrok run-backend run-frontend 

.PHONY: clean
clean:
	cd $(BACKEND_DIR) && $(GOCLEAN)
	cd $(FRONTEND_DIR) && rm -rf dist node_modules

.PHONY: install
install:
	cd $(BACKEND_DIR) && go mod download
	cd $(FRONTEND_DIR) && $(NPM) install

.PHONY: start-ngrok
start-ngrok:
	@if ! pgrep -x "ngrok" > /dev/null; then \
		echo "Starting ngrok..."; \
		$(NGROK) http 8080 > /dev/null & \
		sleep 5; \
		NGROK_URL=$$(curl -s http://localhost:4040/api/tunnels | jq -r '.tunnels[0].public_url'); \
		echo "Ngrok URL: $$NGROK_URL"; \
		sed -i.bak "s|CALLBACK_URL=.*|CALLBACK_URL=$$NGROK_URL|" $(BACKEND_DIR)/.env; \
		echo "Updated CALLBACK_URL in .env"; \
	else \
		echo "Ngrok is already running."; \
	fi