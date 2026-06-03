// Auth
export interface User {
  id: string
  tenant_id: string
  branch_id: string
  username: string
  email: string
  full_name: string
  role: UserRole
  phone: string
  is_active: boolean
  created_at: string
}

export type UserRole =
  | 'super_admin'
  | 'admin_cabang'
  | 'service_advisor'
  | 'mekanik'
  | 'kasir'
  | 'gudang'
  | 'pelanggan'

export interface Branch {
  id: string
  name: string
  code: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
  branch: Branch | null
}

// Dashboard
export interface DashboardStats {
  today_vehicles_in: number
  today_revenue: number
  week_revenue: number
  month_revenue: number
  active_work_orders: number
  waiting_queue: number
  in_progress: number
  completed_today: number
  low_stock_items: number
  total_customers: number
  total_spareparts: number
}

export interface RevenueChart {
  date: string
  amount: number
}

// Customer
export interface Customer {
  id: string
  tenant_id: string
  branch_id?: string
  code?: string
  full_name: string
  phone: string
  email?: string
  address?: string
  id_number?: string
  category: 'regular' | 'member' | 'fleet'
  loyalty_points: number
  tags: string[]
  notes?: string
  is_blacklisted: boolean
  created_at: string
  updated_at: string
  vehicle_count: number
}

// Vehicle
export interface Vehicle {
  id: string
  tenant_id: string
  customer_id?: string
  plate_number: string
  brand: string
  model?: string
  variant?: string
  year?: number
  color?: string
  vin?: string
  engine_number?: string
  last_km?: number
  next_service_km?: number
  notes?: string
  created_at: string
  updated_at: string
  customer_name?: string
}

// Work Order
export interface WorkOrder {
  id: string
  tenant_id: string
  branch_id?: string
  wo_number: string
  customer_id?: string
  vehicle_id?: string
  status: WOStatus
  complaint: string
  diagnosis?: string
  mechanic_notes?: string
  estimated_cost: number
  estimated_hours?: number
  actual_cost: number
  labor_cost: number
  total_sparepart_cost: number
  discount: number
  tax_amount: number
  grand_total: number
  customer_approved: boolean
  queue_number?: number
  created_by?: string
  started_at?: string
  completed_at?: string
  created_at: string
  updated_at: string
  customer_name?: string
  vehicle_plate?: string
  services?: WOService[]
  spareparts?: WOSparepart[]
  mechanics?: WOMechanic[]
}

export type WOStatus =
  | 'menunggu'
  | 'diagnosis'
  | 'menunggu_persetujuan'
  | 'dikerjakan'
  | 'qc'
  | 'selesai'
  | 'diambil'
  | 'batal'

export interface WOService {
  id: string
  work_order_id: string
  service_name: string
  description?: string
  quantity: number
  unit_price: number
  total_price: number
}

export interface WOSparepart {
  id: string
  work_order_id: string
  sparepart_id: string
  sparepart_name: string
  quantity: number
  unit_price: number
  total_price: number
}

export interface WOMechanic {
  id: string
  work_order_id: string
  user_id: string
  commission: number
  assigned_at: string
  full_name?: string
}

// Sparepart
export interface Sparepart {
  id: string
  tenant_id: string
  branch_id?: string
  code: string
  name: string
  category: string
  brand?: string
  unit: string
  purchase_price: number
  selling_price: number
  wholesale_price?: number
  member_price?: number
  current_stock: number
  min_stock: number
  max_stock?: number
  barcode?: string
  location?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface StockMutation {
  id: string
  tenant_id: string
  sparepart_id: string
  mutation_type: 'masuk' | 'keluar' | 'retur' | 'opname' | 'rusak'
  quantity: number
  reference_type?: string
  reference_id?: string
  notes?: string
  created_at: string
  sparepart_name?: string
}

// Supplier
export interface Supplier {
  id: string
  tenant_id: string
  code?: string
  name: string
  contact_person?: string
  phone?: string
  email?: string
  address?: string
  tax_id?: string
  bank_account?: string
  notes?: string
  is_active: boolean
  created_at: string
}

export interface PurchaseOrder {
  id: string
  po_number: string
  supplier_id: string
  status: string
  total_amount: number
  notes?: string
  created_at: string
  supplier_name?: string
  items?: POItem[]
}

export interface POItem {
  sparepart_id?: string
  item_name: string
  quantity: number
  unit_price: number
  total_price: number
}

// Invoice & Payment
export interface Invoice {
  id: string
  invoice_number: string
  work_order_id?: string
  customer_id?: string
  subtotal: number
  discount: number
  tax_amount: number
  grand_total: number
  amount_paid: number
  balance_due: number
  status: 'belum_bayar' | 'dp' | 'lunas' | 'batal'
  notes?: string
  created_at: string
  paid_at?: string
  customer_name?: string
  payments?: Payment[]
}

export interface Payment {
  id: string
  amount: number
  method: string
  reference?: string
  notes?: string
  created_at: string
}

// Employee / Mechanic
export interface Mechanic {
  id: string
  full_name: string
  phone: string
  skills: string[]
  is_active: boolean
  active_jobs: number
}

// Paginated Response
export interface PaginatedResponse<T> {
  data: T[]
  total?: number
  page: number
  limit: number
}
