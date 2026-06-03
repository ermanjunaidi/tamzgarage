# BengkelPro Makefile
# Gunakan: make <target>

.PHONY: dev up start dev-backend dev-frontend prod down logs build build-prod reset tunnel tunnel-url

# ==========================================================================
# 🚀 Mode Development (hot reload)
# ==========================================================================

# Jalankan semuanya lokal: PostgreSQL di Docker, backend Go, frontend Vite (hot reload)
dev:
	@echo "🔧 Membersihkan port & memastikan PostgreSQL berjalan..."
	@lsof -ti:8082 | xargs kill -9 2>/dev/null || true
	@lsof -ti:5173 | xargs kill -9 2>/dev/null || true
	@docker compose up -d postgres 2>/dev/null || true
	@echo "⏳ Menunggu PostgreSQL siap..."
	@for i in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15; do \
		docker compose exec -T postgres pg_isready -U bengkelpro 2>/dev/null && break; \
		sleep 1; \
	done
	@echo ""
	@echo "═══════════════════════════════════════════"
	@echo "  BengkelPro Development Mode"
	@echo "═══════════════════════════════════════════"
	@echo "  → Frontend : http://localhost:5173"
	@echo "  → Backend  : http://localhost:8082"
	@echo "  → Database : localhost:5432"
	@echo ""
	@echo "  Default login: admin / admin123"
	@echo "═══════════════════════════════════════════"
	@echo ""
	@trap 'kill 0' EXIT; \
	(cd backend   && [ -f go.sum ] || go mod tidy && go run .) & \
	(cd frontend  && [ -d node_modules ] || npm install && npm run dev) & \
	wait

# Jalankan backend Go saja (lokal)
dev-backend:
	@lsof -ti:8082 | xargs kill -9 2>/dev/null || true
	@docker compose up -d postgres 2>/dev/null || true
	@sleep 2
	@cd backend && ([ -f go.sum ] || go mod tidy) && go run .

# Jalankan frontend Vite saja (lokal)
dev-frontend:
	@cd frontend && ([ -d node_modules ] || npm install) && npm run dev

# Build & start semua service di Docker
up:
	@docker compose up -d --build
	@echo ""
	@echo "═══════════════════════════════════════════"
	@echo "  Semua service berjalan di Docker"
	@echo "  → Frontend : http://localhost:5173"
	@echo "  → Backend  : http://localhost:8082"
	@echo "═══════════════════════════════════════════"

# Start service Docker tanpa rebuild
start:
	@docker compose up -d

# ==========================================================================
# 🏭 Mode Production
# ==========================================================================

# Build & start production — satu port :8080
prod:
	@docker compose -f docker-compose.prod.yml up -d --build
	@echo ""
	@echo "═══════════════════════════════════════════"
	@echo "  Production Mode"
	@echo "  → URL : http://localhost:8082"
	@echo "═══════════════════════════════════════════"

# Build image production saja
build-prod:
	@docker build -f Dockerfile.prod -t bengkelpro:latest .
	@echo "✅ Image bengkelpro:latest siap"

# ==========================================================================
# 🌐 Cloudflare Tunnel
# ==========================================================================

# Jalankan tunnel saja (attach ke terminal, lihat URL langsung)
tunnel:
	@cloudflared tunnel --url http://localhost:8082

# Tampilkan URL tunnel dari container
# Gunakan setelah: make prod-tunnel
tunnel-url:
	@docker logs bengkelpro-tunnel 2>/dev/null | grep -o 'https://[^ ]*\.trycloudflare\.com' | tail -1 || echo "Tunnel tidak aktif. Jalankan: make prod-tunnel"

# Production + tunnel (semua background, ambil URL setelah 5 detik)
prod-tunnel:
	@docker compose -f docker-compose.prod.yml --profile tunnel up -d --build
	@sleep 5
	@echo ""
	@echo "═══════════════════════════════════════════"
	@echo "  Production + Tunnel"
	@echo "  → Local  : http://localhost:8082"
	@echo "  → Tunnel : $$(docker logs bengkelpro-tunnel 2>/dev/null | grep -o 'https://[^ ]*\.trycloudflare\.com' | tail -1 || echo 'menunggu...')"
	@echo "═══════════════════════════════════════════"

# ==========================================================================
# 🛠 Utilitas
# ==========================================================================

# Streaming logs semua container
logs:
	@docker compose logs -f

# Hentikan semua service (dev + prod)
down:
	@docker compose down 2>/dev/null || true
	@docker compose -f docker-compose.prod.yml down 2>/dev/null || true
	@echo "✅ Semua service dihentikan"

# Build image Docker untuk development
build:
	@docker compose build

# Hentikan service + hapus volume database (data hilang!)
reset:
	@docker compose down -v 2>/dev/null || true
	@docker compose -f docker-compose.prod.yml down -v 2>/dev/null || true
	@echo "⚠️  Database direset — semua data dihapus"
