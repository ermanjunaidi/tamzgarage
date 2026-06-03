# Product Requirements Document (PRD)
# Sistem Manajemen Bengkel — "BengkelPro"

---

## Daftar Isi
1. [Ringkasan Eksekutif](#1-ringkasan-eksekutif)
2. [Tujuan Produk](#2-tujuan-produk)
3. [Persona & Role Pengguna](#3-persona--role-pengguna)
4. [Cakupan & Modul](#4-cakupan--modul)
5. [Functional Requirements](#5-functional-requirements)
   - [5.1 Autentikasi & Otorisasi](#51-autentikasi--otorisasi)
   - [5.2 Dashboard](#52-dashboard)
   - [5.3 Manajemen Pelanggan](#53-manajemen-pelanggan)
   - [5.4 Manajemen Kendaraan](#54-manajemen-kendaraan)
   - [5.5 Manajemen Order Servis (Work Order)](#55-manajemen-order-servis-work-order)
   - [5.6 Manajemen Antrian](#56-manajemen-antrian)
   - [5.7 Manajemen Inventori & Suku Cadang](#57-manajemen-inventori--suku-cadang)
   - [5.8 Manajemen Mekanik & Karyawan](#58-manajemen-mekanik--karyawan)
   - [5.9 Point of Sale (POS) / Kasir](#59-point-of-sale-pos--kasir)
   - [5.10 Manajemen Supplier](#510-manajemen-supplier)
   - [5.11 Pembayaran & Invoice](#511-pembayaran--invoice)
   - [5.12 Laporan & Analytics](#512-laporan--analytics)
   - [5.13 Notifikasi](#513-notifikasi)
   - [5.14 Multi-Cabang](#514-multi-cabang)
   - [5.15 Integrasi](#515-integrasi)
6. [Non-Functional Requirements](#6-non-functional-requirements)
7. [Arsitektur & Tech Stack (Rekomendasi)](#7-arsitektur--tech-stack-rekomendasi)
8. [UI/UX Requirements](#8-uiux-requirements)
9. [Milestone & Fase Pengembangan](#9-milestone--fase-pengembangan)
10. [Risk & Mitigasi](#10-risk--mitigasi)

---

## 1. Ringkasan Eksekutif

**BengkelPro** adalah sistem manajemen bengkel terintegrasi berbasis web yang dirancang untuk bengkel kendaraan (mobil & motor) skala kecil hingga menengah. Sistem ini mencakup seluruh siklus operasional bengkel: mulai dari pendaftaran pelanggan, pencatatan kendaraan, manajemen antrian servis, alokasi mekanik, manajemen inventori suku cadang, kasir/POS, hingga pelaporan keuangan dan operasional.

Produk hadir sebagai solusi **SaaS multi-tenant** dengan kemampuan multi-cabang, memungkinkan pemilik bengkel mengelola beberapa cabang dari satu dashboard terpusat.

---

## 2. Tujuan Produk

| # | Tujuan | Key Metric |
|---|--------|------------|
| 1 | Mendigitalisasi operasional bengkel dari manual (kertas) ke digital | Pengurangan waktu administrasi 60% |
| 2 | Meningkatkan transparansi harga & servis ke pelanggan | Kepuasan pelanggan > 4.5/5 |
| 3 | Optimalisasi inventori suku cadang | Stok mati berkurang 30% |
| 4 | Mempercepat proses servis dari antrian → selesai | Throughput kendaraan +25% |
| 5 | Menyediakan laporan keuangan real-time untuk pemilik | Laporan tersedia <5 detik |

---

## 3. Persona & Role Pengguna

### Role & Permission Matrix

| Role | Deskripsi | Akses Kunci |
|------|-----------|-------------|
| **Super Admin** | Pemilik bengkel / owner | Semua akses, multi-cabang, laporan keuangan |
| **Admin Cabang** | Kepala cabang | Kelola cabang sendiri, semua modul |
| **Service Advisor (SA)** | Penerima pelanggan | Buat order, cek stok, atur jadwal |
| **Mekanik / Teknisi** | Pelaksana servis | Lihat work order, update progress, catat suku cadang terpakai |
| **Kasir** | Penanganan pembayaran | Invoice, pembayaran, tutup kas |
| **Gudang / Inventory** | Pengelola stok | Stok masuk/keluar, PO, stok opname |
| **Pelanggan** | End-user | Lihat history servis, booking online, lacak progress |

---

## 4. Cakupan & Modul

```
BengkelPro
├── 01. Autentikasi & Manajemen User
├── 02. Dashboard (Overview & KPI)
├── 03. Manajemen Pelanggan (CRM)
├── 04. Manajemen Kendaraan
├── 05. Work Order / Order Servis
├── 06. Manajemen Antrian (Queue)
├── 07. Inventori & Suku Cadang
├── 08. Manajemen Mekanik & Karyawan
├── 09. Point of Sale (POS) / Kasir
├── 10. Manajemen Supplier
├── 11. Invoice & Pembayaran
├── 12. Laporan & Analytics
├── 13. Notifikasi (WhatsApp/Email)
├── 14. Multi-Cabang
└── 15. Integrasi Eksternal
```

---

## 5. Functional Requirements

### 5.1 Autentikasi & Otorisasi

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| AUTH-01 | Login dengan username/email + password | P0 | — |
| AUTH-02 | Role-based access control (RBAC) dengan 7 role di atas | P0 | — |
| AUTH-03 | Multi-tenant: setiap bengkel terisolasi datanya | P0 | — |
| AUTH-04 | Reset password via email | P1 | — |
| AUTH-05 | 2FA (WhatsApp OTP) untuk role admin ke atas | P2 | — |
| AUTH-06 | Session timeout otomatis 30 menit | P1 | — |
| AUTH-07 | Audit log semua aktivitas user (create, update, delete) | P1 | — |

### 5.2 Dashboard

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| DASH-01 | Ringkasan kendaraan masuk/hari ini | P0 | — |
| DASH-02 | KPI: total revenue hari ini, minggu ini, bulan ini | P0 | — |
| DASH-03 | Grafik revenue 30 hari terakhir | P1 | — |
| DASH-04 | Status antrian real-time (menunggu, dikerjakan, selesai) | P0 | — |
| DASH-05 | Alert stok suku cadang menipis | P1 | — |
| DASH-06 | Performa mekanik (job selesai/hari) | P2 | — |
| DASH-07 | Quick action: buat order baru, cari pelanggan | P0 | — |

### 5.3 Manajemen Pelanggan

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| CRM-01 | Tambah/edit/hapus/arsip pelanggan | P0 | — |
| CRM-02 | Pencarian berdasarkan nama, no HP, no plat | P0 | — |
| CRM-03 | Riwayat servis lengkap per pelanggan | P0 | — |
| CRM-04 | Kategori pelanggan: regular, member, fleet (perusahaan) | P1 | — |
| CRM-05 | Program loyalty: poin per transaksi → diskon | P2 | — |
| CRM-06 | Import/export data pelanggan (CSV/Excel) | P2 | — |
| CRM-07 | Tagging & notes khusus per pelanggan | P1 | — |
| CRM-08 | Blacklist pelanggan (riwayat tidak bayar) | P2 | — |

### 5.4 Manajemen Kendaraan

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| VEH-01 | Tambah kendaraan: no plat, merk, tipe, tahun, warna | P0 | — |
| VEH-02 | Relasi kendaraan ↔ pelanggan (one-to-many) | P0 | — |
| VEH-03 | Catat nomor rangka (VIN) & nomor mesin | P1 | — |
| VEH-04 | KM terakhir & KM servis berikutnya | P1 | — |
| VEH-05 | Riwayat servis per kendaraan (spare part, mekanik, keluhan) | P0 | — |
| VEH-06 | Reminder servis berkala via WhatsApp/Email | P2 | — |
| VEH-07 | Upload foto kondisi kendaraan (pre/post servis) | P2 | — |

### 5.5 Manajemen Order Servis (Work Order)

**Core Module** — ini adalah workflow utama bengkel.

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| WO-01 | Buat Work Order (WO) baru: pelanggan, kendaraan, keluhan | P0 | — |
| WO-02 | Status WO: **Menunggu → Diagnosis → Menunggu Persetujuan → Dikerjakan → QC → Selesai → Diambil** | P0 | — |
| WO-03 | Estimasi biaya & waktu pengerjaan | P0 | — |
| WO-04 | Persetujuan pelanggan (approval) untuk servis berbayar | P0 | — |
| WO-05 | Alokasi mekanik ke WO (bisa multiple) | P0 | — |
| WO-06 | Catat suku cadang yang digunakan per WO | P0 | — |
| WO-07 | Catat jasa / labor cost per WO | P0 | — |
| WO-08 | Checklist inspeksi kendaraan (pre-service inspection) | P1 | — |
| WO-09 | Print WO & checklist fisik | P1 | — |
| WO-10 | Upsell rekomendasi servis tambahan (berdasarkan temuan) | P2 | — |
| WO-11 | Barcode/QR Code untuk tracking WO | P2 | — |
| WO-12 | Integrasi dengan antrian digital (TV display) | P1 | — |

**Detail Workflow WO:**

```
Pelanggan Datang
       │
       ▼
SA Buat WO ──► Catat keluhan, estimasi awal
       │
       ▼
Antrian Menunggu ──► Tampil di display antrian
       │
       ▼
Mekanik Diagnosis ──► Catat temuan, suku cadang dibutuhkan
       │
       ▼
Menunggu Persetujuan ──► Pelanggan approve/reject estimasi biaya
       │
       ▼
Dikerjakan ──► Mekanik kerjakan, catat spare part terpakai
       │
       ▼
Quality Control ──► QC oleh kepala mekanik / admin
       │
       ▼
Selesai ──► Notifikasi ke pelanggan
       │
       ▼
Diambil ──► Pembayaran, kendaraan diambil
```

### 5.6 Manajemen Antrian

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| QUE-01 | Nomor antrian otomatis per hari | P0 | — |
| QUE-02 | Display antrian real-time (TV/monitor) | P1 | — |
| QUE-03 | Estimasi waktu tunggu | P1 | — |
| QUE-04 | Prioritas antrian: booking vs walk-in | P1 | — |
| QUE-05 | Pelanggan bisa lihat status via link/WhatsApp | P2 | — |
| QUE-06 | Panggil antrian via speaker/TV | P2 | — |

### 5.7 Manajemen Inventori & Suku Cadang

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| INV-01 | Master data suku cadang: kode, nama, kategori, merk, harga beli, harga jual | P0 | — |
| INV-02 | Kategori: oli, rem, ban, aki, filter, kelistrikan, mesin, body, dll | P0 | — |
| INV-03 | Multi-satuan: pcs, liter, set, box | P1 | — |
| INV-04 | Stok masuk (penerimaan dari supplier/PO) | P0 | — |
| INV-05 | Stok keluar (terpakai di WO, rusak, retur) | P0 | — |
| INV-06 | Stok minimum alert (reorder point) | P0 | — |
| INV-07 | Stok opname periodik | P1 | — |
| INV-08 | Serial number & batch number (opsional) | P2 | — |
| INV-09 | Multi-gudang per cabang | P2 | — |
| INV-10 | Harga bertingkat: eceran, grosir, member | P1 | — |
| INV-11 | Barcode scanning untuk stok opname | P2 | — |

### 5.8 Manajemen Mekanik & Karyawan

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| EMP-01 | Data karyawan: nama, alamat, kontak, role | P0 | — |
| EMP-02 | Keahlian/skill mekanik (mesin, kaki-kaki, AC, kelistrikan, dll) | P1 | — |
| EMP-03 | Status ketersediaan mekanik (available/busy) | P0 | — |
| EMP-04 | Beban kerja mekanik (jumlah WO aktif) | P1 | — |
| EMP-05 | Produktivitas mekanik (WO selesai/hari, revenue generated) | P1 | — |
| EMP-06 | Komisi mekanik per WO (opsional) | P2 | — |
| EMP-07 | Absensi (check-in/check-out) | P2 | — |
| EMP-08 | Gaji & slip gaji | P3 | — |

### 5.9 Point of Sale (POS) / Kasir

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| POS-01 | Pembayaran WO (servis + spare part) | P0 | — |
| POS-02 | Penjualan langsung suku cadang (tanpa servis) | P1 | — |
| POS-03 | Multi-metode pembayaran: tunai, debit, kredit, QRIS, transfer | P0 | — |
| POS-04 | Diskon per item atau per transaksi (%) | P1 | — |
| POS-05 | Pajak: PPN 11% (opsional per transaksi) | P1 | — |
| POS-06 | Print struk thermal 58mm/80mm | P1 | — |
| POS-07 | Void / pembatalan transaksi (dengan alasan) | P1 | — |
| POS-08 | Tutup kas (cash register close) per shift | P1 | — |
| POS-09 | Refund / pengembalian dana | P2 | — |

### 5.10 Manajemen Supplier

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| SUP-01 | Data supplier: nama, kontak, alamat, NPWP | P0 | — |
| SUP-02 | Katalog suku cadang per supplier | P1 | — |
| SUP-03 | Purchase Order (PO) ke supplier | P1 | — |
| SUP-04 | Penerimaan PO → update stok otomatis | P1 | — |
| SUP-05 | Retur ke supplier | P2 | — |
| SUP-06 | Hutang supplier (account payable) | P2 | — |
| SUP-07 | Histori pembelian per supplier | P1 | — |

### 5.11 Pembayaran & Invoice

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| PAY-01 | Generate invoice otomatis dari WO | P0 | — |
| PAY-02 | Nomor invoice sequential per cabang | P0 | — |
| PAY-03 | Invoice itemized: jasa + spare part + diskon + pajak | P0 | — |
| PAY-04 | Status pembayaran: lunas, DP, belum bayar | P0 | — |
| PAY-05 | Cicilan / pembayaran bertahap | P2 | — |
| PAY-06 | Piutang pelanggan (account receivable) | P2 | — |
| PAY-07 | Kirim invoice via WhatsApp/Email | P1 | — |
| PAY-08 | Print invoice A4 | P1 | — |

### 5.12 Laporan & Analytics

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| RPT-01 | Laporan pendapatan harian/mingguan/bulanan/tahunan | P0 | — |
| RPT-02 | Laporan revenue by kategori (jasa vs spare part) | P1 | — |
| RPT-03 | Laporan stok: stok tersedia, stok menipis, stok mati | P0 | — |
| RPT-04 | Laporan stok masuk/keluar per periode | P1 | — |
| RPT-05 | Laporan WO: total, selesai, batal per periode | P1 | — |
| RPT-06 | Laporan produktivitas mekanik | P2 | — |
| RPT-07 | Laporan pelanggan: top customer, customer baru | P2 | — |
| RPT-08 | Laporan laba rugi sederhana | P1 | — |
| RPT-09 | Laporan kas & bank | P2 | — |
| RPT-10 | Laporan PPN / pajak | P2 | — |
| RPT-11 | Export semua laporan ke Excel/PDF | P1 | — |
| RPT-12 | Dashboard analytics dengan filter tanggal | P0 | — |

### 5.13 Notifikasi

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| NOT-01 | Notifikasi status WO ke pelanggan via WhatsApp | P1 | — |
| NOT-02 | Notifikasi stok menipis ke admin gudang | P1 | — |
| NOT-03 | Reminder servis berkala ke pelanggan | P2 | — |
| NOT-04 | Broadcast WhatsApp ke pelanggan (promo) | P3 | — |
| NOT-05 | In-app notification untuk internal user | P1 | — |

### 5.14 Multi-Cabang

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| MC-01 | Satu database, banyak cabang (tenant isolation) | P0 | — |
| MC-02 | Dashboard konsolidasi multi-cabang untuk Super Admin | P1 | — |
| MC-03 | Transfer stok antar cabang | P2 | — |
| MC-04 | Perbandingan performa antar cabang | P2 | — |
| MC-05 | Konfigurasi independen per cabang (pajak, diskon, harga) | P1 | — |

### 5.15 Integrasi

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| INT-01 | WhatsApp Business API (notifikasi) | P1 | — |
| INT-02 | Midtrans / Xendit / payment gateway (QRIS, VA) | P1 | — |
| INT-03 | Thermal printer ESC/POS | P1 | — |
| INT-04 | Barcode scanner untuk inventori | P2 | — |
| INT-05 | API publik untuk integrasi pihak ketiga | P3 | — |

---

## 6. Non-Functional Requirements

| ID | Kategori | Requirement |
|----|----------|-------------|
| NFR-01 | **Performa** | Halaman utama load < 2 detik, API response < 500ms |
| NFR-02 | **Skalabilitas** | Mendukung 50+ cabang, 500+ user bersamaan |
| NFR-03 | **Keamanan** | HTTPS, enkripsi password (bcrypt), SQL injection prevention, XSS protection, CSRF token |
| NFR-04 | **Ketersediaan** | Uptime 99.5% (SLA) |
| NFR-05 | **Kompatibilitas** | Browser: Chrome, Firefox, Edge (2 versi terakhir). Mobile-responsive |
| NFR-06 | **Backup** | Database backup harian otomatis, retensi 30 hari |
| NFR-07 | **Audit Trail** | Semua perubahan data kritis tercatat (who, when, what) |
| NFR-08 | **Offline Mode** | POS/kasir bisa beroperasi saat internet down (P2) |
| NFR-09 | **Lokalisasi** | Bahasa Indonesia, format Rupiah (IDR), timezone WIB |
| NFR-10 | **Aksesibilitas** | Kontras warna cukup, font readable (WCAG AA) |

---

## 7. Arsitektur & Tech Stack (Rekomendasi)

### Arsitektur
```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Frontend   │────▶│   API Layer   │────▶│  Database    │
│  (Next.js)   │     │  (Node/Go)    │     │ (PostgreSQL) │
└─────────────┘     └──────────────┘     └─────────────┘
                            │
              ┌─────────────┼─────────────┐
              ▼             ▼             ▼
        ┌──────────┐ ┌──────────┐ ┌──────────────┐
        │  Redis    │ │  Queue   │ │ Object Store  │
        │  (Cache)  │ │ (Rabbit) │ │  (S3/MinIO)   │
        └──────────┘ └──────────┘ └──────────────┘
```

### Rekomendasi Tech Stack

| Layer | Teknologi | Alasan |
|-------|-----------|--------|
| **Frontend** | Next.js 14 + React + TypeScript + TailwindCSS + Shadcn/ui | SSR, SEO, ekosistem mature |
| **Backend API** | Node.js (NestJS) atau Go (Echo/Fiber) | Performa tinggi, strongly-typed |
| **Database Utama** | PostgreSQL 16 | ACID, JSON support, full-text search |
| **Cache** | Redis | Session store, caching dashboard |
| **Message Queue** | Redis BullMQ / RabbitMQ | Notifikasi, background job |
| **File Storage** | MinIO / AWS S3 | Upload foto kendaraan, export laporan |
| **Real-time** | WebSocket (Socket.io) | Update dashboard & antrian real-time |
| **Mobile** | PWA (Progressive Web App) | Tidak perlu native app untuk MVP |

---

## 8. UI/UX Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| UX-01 | Sidebar navigasi dengan grouping modul yang jelas | P0 |
| UX-02 | Tabel data dengan pagination, sorting, filtering, dan search | P0 |
| UX-03 | Dark mode / light mode | P2 |
| UX-04 | Mobile-first responsive design | P0 |
| UX-05 | Form validation real-time dengan error message jelas | P0 |
| UX-06 | Toast notification untuk aksi sukses/error | P0 |
| UX-07 | Loading skeleton untuk data fetching | P1 |
| UX-08 | Keyboard shortcuts untuk aksi cepat (Ctrl+K search) | P2 |
| UX-09 | Confirm dialog sebelum aksi destruktif (delete, void) | P0 |
| UX-10 | Role-based menu visibility (menu yang tidak diizinkan disembunyikan) | P0 |
| UX-11 | Breadcrumb navigasi | P1 |
| UX-12 | Offline indicator (saat koneksi terputus) | P2 |

---

## 9. Milestone & Fase Pengembangan

### Fase 1: MVP (4-6 minggu)
**Core workflow bengkel berjalan end-to-end**

- [ ] Autentikasi & role dasar (Super Admin, Admin, SA, Mekanik, Kasir)
- [ ] Multi-tenant & multi-cabang dasar
- [ ] Manajemen pelanggan + kendaraan
- [ ] Work Order (WO) dengan status workflow penuh
- [ ] Inventori dasar (stok masuk/keluar, stok minimum alert)
- [ ] Kasir/POS dasar (tunai, debit, kredit)
- [ ] Invoice & pembayaran
- [ ] Dashboard sederhana (ringkasan harian)
- [ ] Laporan pendapatan & stok
- [ ] Manajemen mekanik dasar

### Fase 2: Enhancement (4-6 minggu)
**Fitur pendukung & pengalaman lebih baik**

- [ ] Antrian digital dengan TV display
- [ ] Notifikasi WhatsApp ke pelanggan
- [ ] Manajemen supplier + Purchase Order
- [ ] Payment gateway (QRIS, VA) via Midtrans
- [ ] Print thermal & A4
- [ ] Checklist inspeksi kendaraan
- [ ] Program loyalty pelanggan
- [ ] Laporan advanced (produktivitas, laba rugi, top customer)
- [ ] Audit log lengkap

### Fase 3: Scale (4-6 minggu)
**Enterprise-ready & optimasi**

- [ ] Dashboard konsolidasi multi-cabang
- [ ] Transfer stok antar cabang
- [ ] Piutang & hutang
- [ ] Absensi & penggajian
- [ ] Barcode scanning untuk stok opname
- [ ] Reminder servis berkala
- [ ] Broadcast WhatsApp
- [ ] API publik untuk integrasi
- [ ] Offline mode POS
- [ ] 2FA

---

## 10. Risk & Mitigasi

| Risk | Level | Mitigasi |
|------|-------|----------|
| Adopsi user rendah (terbiasa kertas) | High | UI sesederhana mungkin, onboarding in-app, training langsung |
| Data migration dari sistem lama | Medium | Sediakan template import Excel, jasa migrasi |
| Keamanan data pelanggan | High | Enkripsi data sensitif, backup rutin, penetration testing |
| Koneksi internet tidak stabil | Medium | Optimasi bundle size, offline mode untuk POS (Fase 3) |
| Perubahan regulasi perpajakan | Low | Desain pajak fleksibel (configurable tax rate) |
| Konflik multi-tenant | Low | Row-level security di database, testing isolasi ketat |

---

## Appendix A: Glosarium

| Istilah | Definisi |
|---------|----------|
| **WO** | Work Order — dokumen perintah kerja servis |
| **SA** | Service Advisor — petugas penerima pelanggan |
| **POS** | Point of Sale — sistem kasir |
| **PO** | Purchase Order — pesanan pembelian ke supplier |
| **QC** | Quality Control — pemeriksaan kualitas hasil servis |
| **Stok Opname** | Perhitungan fisik stok untuk mencocokkan dengan sistem |
| **RBAC** | Role-Based Access Control — kontrol akses berdasarkan peran |

---

## Appendix B: Entity Relationship (High Level)

```
Tenant (Bengkel)
  └── Cabang (Branch)
        ├── User (Karyawan)
        │     └── Role
        ├── Pelanggan (Customer)
        │     └── Kendaraan (Vehicle)
        │           └── WorkOrder
        │                 ├── WO_Item (Jasa)
        │                 ├── WO_SparePart (Suku cadang terpakai)
        │                 └── WO_Mekanik (Mekanik assigned)
        ├── SparePart (Master suku cadang)
        │     └── StockMutation (Stok masuk/keluar)
        ├── Supplier
        │     └── PurchaseOrder
        ├── Invoice
        │     └── Payment
        └── Queue (Antrian)
```

---

_Dokumen ini adalah living document. Review dan revisi dilakukan setiap awal sprint._

**Versi:** 1.0  
**Tanggal:** 2026-06-03  
**Author:** Tim Produk  
**Status:** Draft — Menunggu Review
