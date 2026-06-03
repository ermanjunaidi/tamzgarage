-- BengkelPro Database Schema
-- PostgreSQL 16

-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================
-- TENANT & BRANCH
-- ============================================================
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    address TEXT,
    phone VARCHAR(30),
    email VARCHAR(100),
    logo_url TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE branches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    address TEXT,
    phone VARCHAR(30),
    email VARCHAR(100),
    tax_rate NUMERIC(5,2) DEFAULT 11.00,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

-- ============================================================
-- USERS & ROLES
-- ============================================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id) ON DELETE SET NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('super_admin','admin_cabang','service_advisor','mekanik','kasir','gudang','pelanggan')),
    phone VARCHAR(30),
    address TEXT,
    skills TEXT[], -- for mekanik: array of skills
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(500) UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id),
    user_id UUID REFERENCES users(id),
    action VARCHAR(50) NOT NULL, -- CREATE, UPDATE, DELETE
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID,
    old_data JSONB,
    new_data JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- CUSTOMERS
-- ============================================================
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id),
    code VARCHAR(50),
    full_name VARCHAR(200) NOT NULL,
    phone VARCHAR(30) NOT NULL,
    email VARCHAR(100),
    address TEXT,
    id_number VARCHAR(50),
    category VARCHAR(20) DEFAULT 'regular' CHECK (category IN ('regular','member','fleet')),
    loyalty_points INT DEFAULT 0,
    tags TEXT[],
    notes TEXT,
    is_blacklisted BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- VEHICLES
-- ============================================================
CREATE TABLE vehicles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,
    plate_number VARCHAR(20) NOT NULL,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100),
    variant VARCHAR(100),
    year INT,
    color VARCHAR(50),
    vin VARCHAR(50),
    engine_number VARCHAR(50),
    last_km INT,
    next_service_km INT,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- WORK ORDERS
-- ============================================================
CREATE TABLE work_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id),
    wo_number VARCHAR(50) NOT NULL,
    customer_id UUID REFERENCES customers(id),
    vehicle_id UUID REFERENCES vehicles(id),
    status VARCHAR(30) DEFAULT 'menunggu' CHECK (status IN ('menunggu','diagnosis','menunggu_persetujuan','dikerjakan','qc','selesai','diambil','batal')),
    complaint TEXT NOT NULL,
    diagnosis TEXT,
    mechanic_notes TEXT,
    estimated_cost NUMERIC(15,2) DEFAULT 0,
    estimated_hours NUMERIC(5,1),
    actual_cost NUMERIC(15,2) DEFAULT 0,
    labor_cost NUMERIC(15,2) DEFAULT 0,
    total_sparepart_cost NUMERIC(15,2) DEFAULT 0,
    discount NUMERIC(15,2) DEFAULT 0,
    tax_amount NUMERIC(15,2) DEFAULT 0,
    grand_total NUMERIC(15,2) DEFAULT 0,
    customer_approved BOOLEAN DEFAULT false,
    queue_number INT,
    created_by UUID REFERENCES users(id),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Work Order Services (jasa)
CREATE TABLE wo_services (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    work_order_id UUID NOT NULL REFERENCES work_orders(id) ON DELETE CASCADE,
    service_name VARCHAR(200) NOT NULL,
    description TEXT,
    quantity NUMERIC(10,2) DEFAULT 1,
    unit_price NUMERIC(15,2) NOT NULL,
    total_price NUMERIC(15,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Work Order Spare Parts
CREATE TABLE wo_spareparts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    work_order_id UUID NOT NULL REFERENCES work_orders(id) ON DELETE CASCADE,
    sparepart_id UUID NOT NULL,
    sparepart_name VARCHAR(200) NOT NULL,
    quantity NUMERIC(10,2) NOT NULL,
    unit_price NUMERIC(15,2) NOT NULL,
    total_price NUMERIC(15,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Work Order Mechanics
CREATE TABLE wo_mechanics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    work_order_id UUID NOT NULL REFERENCES work_orders(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    commission NUMERIC(15,2) DEFAULT 0,
    assigned_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- INVENTORY / SPARE PARTS
-- ============================================================
CREATE TABLE spareparts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    category VARCHAR(50) NOT NULL,
    brand VARCHAR(100),
    unit VARCHAR(20) DEFAULT 'pcs',
    purchase_price NUMERIC(15,2) NOT NULL,
    selling_price NUMERIC(15,2) NOT NULL,
    wholesale_price NUMERIC(15,2),
    member_price NUMERIC(15,2),
    current_stock NUMERIC(10,2) DEFAULT 0,
    min_stock NUMERIC(10,2) DEFAULT 5,
    max_stock NUMERIC(10,2),
    barcode VARCHAR(100),
    location VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE stock_mutations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id),
    sparepart_id UUID NOT NULL REFERENCES spareparts(id),
    mutation_type VARCHAR(20) NOT NULL CHECK (mutation_type IN ('masuk','keluar','retur','opname','rusak')),
    quantity NUMERIC(10,2) NOT NULL,
    reference_type VARCHAR(50),
    reference_id UUID,
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- SUPPLIERS
-- ============================================================
CREATE TABLE suppliers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    code VARCHAR(50),
    name VARCHAR(200) NOT NULL,
    contact_person VARCHAR(200),
    phone VARCHAR(30),
    email VARCHAR(100),
    address TEXT,
    tax_id VARCHAR(50),
    bank_account VARCHAR(100),
    notes TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE purchase_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id),
    po_number VARCHAR(50) NOT NULL,
    supplier_id UUID NOT NULL REFERENCES suppliers(id),
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft','dikirim','diterima','sebagian','batal')),
    total_amount NUMERIC(15,2) DEFAULT 0,
    notes TEXT,
    created_by UUID REFERENCES users(id),
    ordered_at TIMESTAMPTZ,
    received_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE po_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    purchase_order_id UUID NOT NULL REFERENCES purchase_orders(id) ON DELETE CASCADE,
    sparepart_id UUID REFERENCES spareparts(id),
    item_name VARCHAR(200) NOT NULL,
    quantity NUMERIC(10,2) NOT NULL,
    unit_price NUMERIC(15,2) NOT NULL,
    total_price NUMERIC(15,2) NOT NULL,
    received_qty NUMERIC(10,2) DEFAULT 0
);

-- ============================================================
-- INVOICE & PAYMENT
-- ============================================================
CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id),
    invoice_number VARCHAR(50) NOT NULL,
    work_order_id UUID REFERENCES work_orders(id),
    customer_id UUID REFERENCES customers(id),
    subtotal NUMERIC(15,2) NOT NULL,
    discount NUMERIC(15,2) DEFAULT 0,
    tax_amount NUMERIC(15,2) DEFAULT 0,
    grand_total NUMERIC(15,2) NOT NULL,
    amount_paid NUMERIC(15,2) DEFAULT 0,
    balance_due NUMERIC(15,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'belum_bayar' CHECK (status IN ('belum_bayar','dp','lunas','batal')),
    notes TEXT,
    created_by UUID REFERENCES users(id),
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    payment_number VARCHAR(50),
    amount NUMERIC(15,2) NOT NULL,
    method VARCHAR(30) NOT NULL CHECK (method IN ('tunai','debit','kredit','qris','transfer')),
    reference VARCHAR(100),
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- QUEUE
-- ============================================================
CREATE TABLE queues (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id),
    queue_number INT NOT NULL,
    work_order_id UUID REFERENCES work_orders(id),
    status VARCHAR(20) DEFAULT 'menunggu' CHECK (status IN ('menunggu','dipanggil','dikerjakan','selesai')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- INDEXES
-- ============================================================
CREATE INDEX idx_users_tenant ON users(tenant_id);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_customers_tenant ON customers(tenant_id);
CREATE INDEX idx_customers_phone ON customers(phone);
CREATE INDEX idx_vehicles_plate ON vehicles(plate_number);
CREATE INDEX idx_vehicles_customer ON vehicles(customer_id);
CREATE INDEX idx_work_orders_tenant ON work_orders(tenant_id);
CREATE INDEX idx_work_orders_status ON work_orders(status);
CREATE INDEX idx_work_orders_vehicle ON work_orders(vehicle_id);
CREATE INDEX idx_work_orders_customer ON work_orders(customer_id);
CREATE INDEX idx_spareparts_tenant ON spareparts(tenant_id);
CREATE INDEX idx_spareparts_code ON spareparts(code);
CREATE INDEX idx_invoices_tenant ON invoices(tenant_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_stock_mutations_sparepart ON stock_mutations(sparepart_id);
CREATE INDEX idx_queues_tenant_date ON queues(tenant_id, created_at);
CREATE INDEX idx_audit_logs_tenant ON audit_logs(tenant_id);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);

-- ============================================================
-- SEED DATA
-- ============================================================
INSERT INTO tenants (id, name, code, address, phone, email) VALUES
    ('00000000-0000-0000-0000-000000000001', 'BengkelPro Pusat', 'BGP001', 'Jl. Raya No. 1, Jakarta', '021-1234567', 'pusat@bengkelpro.com');

INSERT INTO branches (id, tenant_id, name, code, address, phone) VALUES
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'Cabang Utama', 'CBG001', 'Jl. Raya No. 1, Jakarta', '021-1234567');

-- Password: admin123 (bcrypt hash)
INSERT INTO users (id, tenant_id, branch_id, username, email, password_hash, full_name, role, phone) VALUES
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'admin', 'admin@bengkelpro.com', '$2a$10$ELfHZMPyNoH95ANhszlM8OCtKezqfznEl3GczMgX1AjNppQzXTwyS', 'Super Admin', 'super_admin', '081234567890'),
    ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'sa01', 'sa01@bengkelpro.com', '$2a$10$ELfHZMPyNoH95ANhszlM8OCtKezqfznEl3GczMgX1AjNppQzXTwyS', 'Service Advisor 1', 'service_advisor', '081234567891'),
    ('00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'mekanik01', 'mekanik01@bengkelpro.com', '$2a$10$ELfHZMPyNoH95ANhszlM8OCtKezqfznEl3GczMgX1AjNppQzXTwyS', 'Mekanik 1', 'mekanik', '081234567892'),
    ('00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'kasir01', 'kasir01@bengkelpro.com', '$2a$10$ELfHZMPyNoH95ANhszlM8OCtKezqfznEl3GczMgX1AjNppQzXTwyS', 'Kasir 1', 'kasir', '081234567893'),
    ('00000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'gudang01', 'gudang01@bengkelpro.com', '$2a$10$ELfHZMPyNoH95ANhszlM8OCtKezqfznEl3GczMgX1AjNppQzXTwyS', 'Gudang 1', 'gudang', '081234567894');

-- Sample customers
INSERT INTO customers (id, tenant_id, branch_id, code, full_name, phone, email, address, category) VALUES
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'CUST001', 'Budi Santoso', '081298765432', 'budi@email.com', 'Jl. Melati No. 5, Jakarta', 'regular'),
    ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'CUST002', 'Siti Rahayu', '081298765433', 'siti@email.com', 'Jl. Mawar No. 12, Jakarta', 'member');

-- Sample vehicles
INSERT INTO vehicles (id, tenant_id, customer_id, plate_number, brand, model, year, color, last_km, next_service_km) VALUES
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'B 1234 CD', 'Toyota', 'Avanza', 2020, 'Putih', 45000, 50000),
    ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', 'B 5678 EF', 'Honda', 'Brio', 2021, 'Merah', 30000, 35000);

-- Sample spare parts
INSERT INTO spareparts (id, tenant_id, branch_id, code, name, category, brand, unit, purchase_price, selling_price, current_stock, min_stock) VALUES
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'SP001', 'Oli Mesin 10W-40', 'oli', 'Castrol', 'liter', 45000, 65000, 50, 10),
    ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'SP002', 'Kampas Rem Depan', 'rem', 'Bendix', 'set', 85000, 120000, 20, 5),
    ('00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'SP003', 'Aki Kering 12V 45Ah', 'aki', 'GS', 'pcs', 550000, 750000, 8, 3),
    ('00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'SP004', 'Filter Oli', 'filter', 'Sakura', 'pcs', 25000, 40000, 100, 20),
    ('00000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'SP005', 'Busi Iridium', 'kelistrikan', 'NGK', 'pcs', 75000, 110000, 40, 10),
    ('00000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'SP006', 'Ban 185/70 R14', 'ban', 'Bridgestone', 'pcs', 450000, 650000, 12, 4);

-- Sample supplier
INSERT INTO suppliers (id, tenant_id, name, contact_person, phone, email, address, tax_id) VALUES
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'PT Suku Cadang Jaya', 'Hendra', '021-9876543', 'sales@sukucadangjaya.com', 'Jl. Industri No. 10, Jakarta', '01.234.567.8-901.000');
