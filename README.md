# BiasharaTrack

BiasharaTrack is a full-stack application with a Go backend and a frontend development server, featuring M-Pesa integration for payments.

## Prerequisites

Before running the application, make sure you have the following installed:

- [Go](https://golang.org/doc/install)
- [Node.js and npm](https://nodejs.org/)
- [ngrok](https://ngrok.com/download)
- [Make](https://www.gnu.org/software/make/)

## Quick Start

1. Clone the repository:
```bash
git clone <repository-url>
cd biasharatrack
```

2. Install dependencies:
```bash
make install
```

3. Set up environment variables:
```bash
make create-env
```

4. Fill in the required environment variables in `backend/.env`:
```env
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_ENDPOINT=your_db_host
DB_NAME=your_db_name
DB_PORT=your_db_port
MPESA_CONSUMER_KEY=your_mpesa_key
MPESA_CONSUMER_SECRET=your_mpesa_secret
MPESA_PASSKEY=your_mpesa_passkey
MPESA_BUSINESS_SHORTCODE=your_shortcode
MPESA_ENVIRONMENT=sandbox
CALLBACK_URL=http://localhost:8080
```

5. Start the application:
```bash
make run
```

## Available Commands

```bash
# Display all available commands
make help

# Install all dependencies
make install

# Create .env file in backend directory
make create-env

# Verify if all environment variables are set
make check-env

# Start all servers (frontend, backend, and ngrok)
make start-servers

# Build both backend and frontend
make build

# Run both backend and frontend
make run

# Clean up build artifacts
make clean

# Individual Components
make build-backend    # Build the Go backend
make run-backend     # Run the Go backend server
make build-frontend  # Build the frontend
make run-frontend   # Run the frontend development server
make start-ngrok    # Start ngrok and update callback URL
```
