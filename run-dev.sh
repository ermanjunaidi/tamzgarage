#!/usr/bin/env bash
set -e

# BengkelPro Development Launcher
# Single command to run everything in development mode

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$ROOT_DIR/backend"
FRONTEND_DIR="$ROOT_DIR/frontend"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

cleanup() {
    echo -e "\n${YELLOW}Shutting down...${NC}"
    kill $BACKEND_PID 2>/dev/null || true
    kill $FRONTEND_PID 2>/dev/null || true
    docker compose stop 2>/dev/null || true
    exit 0
}

trap cleanup SIGINT SIGTERM EXIT

echo -e "${BLUE}==========================================${NC}"
echo -e "${BLUE}   BengkelPro Development Mode           ${NC}"
echo -e "${BLUE}==========================================${NC}"

# Step 1: Start PostgreSQL via Docker
echo -e "\n${GREEN}[1/3] Starting PostgreSQL...${NC}"
docker compose up -d postgres

# Wait for Postgres to be healthy
echo -e "${YELLOW}Waiting for PostgreSQL to be ready...${NC}"
for i in $(seq 1 30); do
    if docker compose exec -T postgres pg_isready -U bengkelpro 2>/dev/null; then
        echo -e "${GREEN}PostgreSQL is ready!${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}Failed to connect to PostgreSQL${NC}"
        exit 1
    fi
    sleep 1
done

# Step 2: Start Go Fiber Backend
echo -e "\n${GREEN}[2/3] Starting Go backend (port 8082)...${NC}"
cd "$BACKEND_DIR"

# Install dependencies if needed
if [ ! -f "go.sum" ]; then
    echo "Installing Go dependencies..."
    go mod tidy
fi

go run . &
BACKEND_PID=$!

# Wait for backend to start
sleep 2

# Step 3: Start Vite Frontend
echo -e "\n${GREEN}[3/3] Starting Vite frontend (port 5173)...${NC}"
cd "$FRONTEND_DIR"

if [ ! -d "node_modules" ]; then
    echo "Installing frontend dependencies..."
    npm install
fi

npm run dev &
FRONTEND_PID=$!

echo -e "\n${BLUE}==========================================${NC}"
echo -e "${GREEN}   All services running!${NC}"
echo -e "${BLUE}==========================================${NC}"
echo -e "Frontend:  ${GREEN}http://localhost:5173${NC}"
echo -e "Backend:   ${GREEN}http://localhost:8082${NC}"
echo -e "Database:  ${GREEN}localhost:5432${NC}"
echo -e ""
echo -e "${YELLOW}Default login:${NC}"
echo -e "  Username: admin"
echo -e "  Password: admin123"
echo -e ""
echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"
echo -e "${BLUE}==========================================${NC}"

wait
