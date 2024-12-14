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
	@echo "  create-env     - Create .env file in backend directory"
	@echo "  check-env      - Check if .env is properly filled"
	@echo "  start-servers  - Start both servers after env setup"

.PHONY: create-env
create-env:
	@echo "Creating .env file in backend directory..."
	@if [ ! -f "$(BACKEND_DIR)/.env" ]; then \
		echo "DB_USER=" > $(BACKEND_DIR)/.env; \
		echo "DB_PASSWORD=" >> $(BACKEND_DIR)/.env; \
		echo "DB_ENDPOINT=" >> $(BACKEND_DIR)/.env; \
		echo "DB_NAME=" >> $(BACKEND_DIR)/.env; \
		echo "DB_PORT=" >> $(BACKEND_DIR)/.env; \
		echo "MPESA_CONSUMER_KEY=" >> $(BACKEND_DIR)/.env; \
		echo "MPESA_CONSUMER_SECRET=" >> $(BACKEND_DIR)/.env; \
		echo "MPESA_PASSKEY=" >> $(BACKEND_DIR)/.env; \
		echo "MPESA_BUSINESS_SHORTCODE=" >> $(BACKEND_DIR)/.env; \
		echo "MPESA_ENVIRONMENT=sandbox" >> $(BACKEND_DIR)/.env; \
		echo "CALLBACK_URL=http://localhost:8080" >> $(BACKEND_DIR)/.env; \
		echo ".env file created successfully"; \
		echo "Please fill in the environment variables in $(BACKEND_DIR)/.env before running the servers"; \
	else \
		echo ".env file already exists"; \
	fi

.PHONY: check-env
check-env:
	@echo "Checking if .env file is properly filled..."
	@if [ ! -f "$(BACKEND_DIR)/.env" ]; then \
		echo "Error: .env file does not exist. Run 'make create-env' first."; \
		exit 1; \
	fi
	@for var in DB_USER DB_PASSWORD DB_ENDPOINT DB_NAME DB_PORT MPESA_CONSUMER_KEY MPESA_CONSUMER_SECRET MPESA_PASSKEY MPESA_BUSINESS_SHORTCODE MPESA_ENVIRONMENT CALLBACK_URL; do \
		if ! grep -q "^$$var=.\+" "$(BACKEND_DIR)/.env"; then \
			echo "Error: $$var is not set in .env file"; \
			exit 1; \
		fi \
	done
	@echo "Environment variables are properly set."

.PHONY: check-ngrok
check-ngrok:
	@if ! command -v ngrok >/dev/null 2>&1; then \
		echo "Error: ngrok is not installed."; \
		echo "Please install ngrok from https://ngrok.com/download"; \
		echo "After installation, run 'make start-servers' again."; \
		exit 1; \
	fi

.PHONY: start-ngrok
start-ngrok: check-ngrok
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

.PHONY: start-servers
start-servers: check-env check-ngrok
	@echo "Starting all servers..."
	make -j 3 start-ngrok run-backend run-frontend

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
run: create-env check-env
	make -j 2 start-ngrok run-backend run-frontend

.PHONY: clean
clean:
	cd $(BACKEND_DIR) && $(GOCLEAN)
	cd $(FRONTEND_DIR) && rm -rf dist node_modules

.PHONY: install
install: create-env
	cd $(BACKEND_DIR) && go mod download
	cd $(FRONTEND_DIR) && $(NPM) install