# BengkelPro

**Sistem Manajemen Bengkel Terintegrasi** — Mencakup seluruh siklus operasional bengkel: pelanggan, kendaraan, work order, inventori, kasir, supplier, hingga laporan keuangan.

## Tech Stack

| Layer       | Teknologi                                       |
|-------------|-------------------------------------------------|
| Frontend    | React 19 + TypeScript + Vite 8 + TailwindCSS 4 + shadcn/ui |
| Backend     | Go 1.23 + Fiber v2                              |
| Database    | PostgreSQL 16                                   |
| Dev Server  | Vite (hot reload) + Go (DIY reload)             |
| Production  | Single binary (Go serves frontend + API)        |

## Prerequisites

- [Go](https://go.dev/) 1.23+
- [Node.js](https://nodejs.org/) 20+ & npm
- [Docker](https://www.docker.com/) + [Docker Compose](https://docs.docker.com/compose/)
- [Cloudflared](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/downloads/) (opsional, untuk tunnel)

## Perintah Makefile

### 🚀 Mode Development (hot reload)

| Perintah           | Fungsi                                                                                   |
|--------------------|------------------------------------------------------------------------------------------|
| `make dev`         | Jalankan semuanya lokal: PostgreSQL di Docker, backend Go, frontend Vite (hot reload)    |
| `make up`          | Build & start semua service di Docker — frontend di `:5173`, API di `:8080`              |
| `make start`       | Start service Docker tanpa rebuild                                                       |
| `make dev-backend` | Jalankan backend Go saja (lokal)                                                         |
| `make dev-frontend`| Jalankan frontend Vite saja (lokal)                                                      |

```
make dev
  → Frontend : http://localhost:5173
  → Backend  : http://localhost:8080
```

### 🏭 Mode Production (stabil, single port)

| Perintah           | Fungsi                                                   |
|--------------------|----------------------------------------------------------|
| `make prod`        | Build & start production — satu port `:8080` melayani frontend + API |
| `make build-prod`  | Build image Docker production saja                       |

```
make prod
  → URL : http://localhost:8080
```

### 🌐 Cloudflare Tunnel (akses dari internet)

| Perintah           | Fungsi                                              |
|--------------------|-----------------------------------------------------|
| `make tunnel`      | Jalankan tunnel saja (attach ke terminal)           |
| `make tunnel-url`  | Tampilkan URL tunnel dari logs                      |
| `make prod-tunnel` | Production + tunnel (gabungan)                      |

```
make prod-tunnel
  → Tunnel URL: https://xxx.trycloudflare.com
```

### 🛠 Utilitas

| Perintah    | Fungsi                                            |
|-------------|---------------------------------------------------|
| `make logs` | Streaming logs semua container                    |
| `make down` | Hentikan semua service                            |
| `make build`| Build image Docker untuk development              |
| `make reset`| Hentikan service + hapus volume database (data hilang!) |

## Alur Kerja

### Development (hot reload)

```bash
# Cepat — semua berjalan di host (rekomendasi)
make dev

# Atau via Docker
make up
```

### Production

```bash
# Build & jalankan production
make prod

# Atau dengan akses internet via Cloudflare
make prod-tunnel
```

## Struktur Project

```
tamzgarage/
├── frontend/               # React + Vite + Tailwind + shadcn/ui
│   ├── src/
│   │   ├── pages/          # Dashboard, Pelanggan, Kendaraan, WO, POS, dll.
│   │   ├── components/     # UI (shadcn), Layout, Shared
│   │   ├── context/        # Auth context
│   │   ├── api/            # API client
│   │   └── types/          # TypeScript types
│   ├── package.json
│   └── Dockerfile.dev
├── backend/                # Go + Fiber API
│   ├── main.go
│   ├── handlers/           # 10 handler: auth, dashboard, customer, vehicle, dll.
│   ├── middleware/         # JWT auth & RBAC
│   ├── models/             # Data models
│   ├── routes/             # API routes dengan role-based access
│   ├── database/           # Database connection
│   ├── Dockerfile.dev
│   └── .air.toml           # Hot reload config
├── database/
│   └── init.sql            # Schema + seed data
├── docker-compose.yml      # PostgreSQL + backend + frontend (dev)
├── Dockerfile.prod         # Production single-binary
├── Makefile
└── .env
```

## Login Admin

Panel login di `http://localhost:5173/login`.

| Username    | Password  | Role            |
|-------------|-----------|-----------------|
| `admin`     | admin123  | Super Admin     |
| `sa01`      | admin123  | Service Advisor |
| `mekanik01` | admin123  | Mekanik         |
| `kasir01`   | admin123  | Kasir           |
| `gudang01`  | admin123  | Gudang          |

## 10 Modul Utama

| Modul             | Deskripsi                                                |
|-------------------|----------------------------------------------------------|
| Dashboard         | KPI: kendaraan masuk, revenue, antrian, stok menipis     |
| Pelanggan         | CRUD, search, kategori (regular/member/fleet), loyalitas |
| Kendaraan         | Plat, merk, KM, riwayat servis                           |
| Work Order        | Workflow 7 status: menunggu → diagnosis → ... → diambil  |
| Inventori         | Stok masuk/keluar, alert minimum, mutasi stok            |
| Karyawan          | Data mekanik, keahlian, beban kerja                      |
| POS / Kasir       | Invoice, multi-payment (tunai, debit, kredit, QRIS)      |
| Supplier          | Data supplier & Purchase Order                           |
| Laporan           | Revenue, stok, WO, top customers                         |
| Multi-Cabang      | Tenant isolation, dashboard konsolidasi (coming soon)    |
