package models

import (
	"time"

	"github.com/google/uuid"
)

type Sparepart struct {
	ID             uuid.UUID  `json:"id"`
	TenantID       uuid.UUID  `json:"tenant_id"`
	BranchID       *uuid.UUID `json:"branch_id"`
	Code           string     `json:"code"`
	Name           string     `json:"name"`
	Category       string     `json:"category"`
	Brand          *string    `json:"brand"`
	Unit           string     `json:"unit"`
	PurchasePrice  float64    `json:"purchase_price"`
	SellingPrice   float64    `json:"selling_price"`
	WholesalePrice *float64   `json:"wholesale_price"`
	MemberPrice    *float64   `json:"member_price"`
	CurrentStock   float64    `json:"current_stock"`
	MinStock       float64    `json:"min_stock"`
	MaxStock       *float64   `json:"max_stock"`
	Barcode        *string    `json:"barcode"`
	Location       *string    `json:"location"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type StockMutation struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	BranchID      *uuid.UUID `json:"branch_id"`
	SparepartID   uuid.UUID  `json:"sparepart_id"`
	MutationType  string     `json:"mutation_type"`
	Quantity      float64    `json:"quantity"`
	ReferenceType *string    `json:"reference_type"`
	ReferenceID   *uuid.UUID `json:"reference_id"`
	Notes         *string    `json:"notes"`
	CreatedBy     *uuid.UUID `json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
	SparepartName *string    `json:"sparepart_name,omitempty"`
}

type Supplier struct {
	ID            uuid.UUID `json:"id"`
	TenantID      uuid.UUID `json:"tenant_id"`
	Code          *string   `json:"code"`
	Name          string    `json:"name"`
	ContactPerson *string   `json:"contact_person"`
	Phone         *string   `json:"phone"`
	Email         *string   `json:"email"`
	Address       *string   `json:"address"`
	TaxID         *string   `json:"tax_id"`
	BankAccount   *string   `json:"bank_account"`
	Notes         *string   `json:"notes"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PurchaseOrder struct {
	ID         uuid.UUID  `json:"id"`
	TenantID   uuid.UUID  `json:"tenant_id"`
	BranchID   *uuid.UUID `json:"branch_id"`
	PONumber   string     `json:"po_number"`
	SupplierID uuid.UUID  `json:"supplier_id"`
	Status     string     `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	Notes      *string    `json:"notes"`
	CreatedBy  *uuid.UUID `json:"created_by"`
	OrderedAt  *time.Time `json:"ordered_at"`
	ReceivedAt *time.Time `json:"received_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	SupplierName *string  `json:"supplier_name,omitempty"`
	Items        []POItem `json:"items,omitempty"`
}

type POItem struct {
	ID              uuid.UUID  `json:"id"`
	PurchaseOrderID uuid.UUID  `json:"purchase_order_id"`
	SparepartID     *uuid.UUID `json:"sparepart_id"`
	ItemName        string     `json:"item_name"`
	Quantity        float64    `json:"quantity"`
	UnitPrice       float64    `json:"unit_price"`
	TotalPrice      float64    `json:"total_price"`
	ReceivedQty     float64    `json:"received_qty"`
}

type Invoice struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	BranchID      *uuid.UUID `json:"branch_id"`
	InvoiceNumber string     `json:"invoice_number"`
	WorkOrderID   *uuid.UUID `json:"work_order_id"`
	CustomerID    *uuid.UUID `json:"customer_id"`
	Subtotal      float64    `json:"subtotal"`
	Discount      float64    `json:"discount"`
	TaxAmount     float64    `json:"tax_amount"`
	GrandTotal    float64    `json:"grand_total"`
	AmountPaid    float64    `json:"amount_paid"`
	BalanceDue    float64    `json:"balance_due"`
	Status        string     `json:"status"`
	Notes         *string    `json:"notes"`
	CreatedBy     *uuid.UUID `json:"created_by"`
	PaidAt        *time.Time `json:"paid_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CustomerName  *string    `json:"customer_name,omitempty"`
	Payments      []Payment  `json:"payments,omitempty"`
}

type Payment struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	BranchID      *uuid.UUID `json:"branch_id"`
	InvoiceID     uuid.UUID  `json:"invoice_id"`
	PaymentNumber *string    `json:"payment_number"`
	Amount        float64    `json:"amount"`
	Method        string     `json:"method"`
	Reference     *string    `json:"reference"`
	Notes         *string    `json:"notes"`
	CreatedBy     *uuid.UUID `json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
}

type DashboardStats struct {
	TodayVehiclesIn    int     `json:"today_vehicles_in"`
	TodayRevenue       float64 `json:"today_revenue"`
	WeekRevenue        float64 `json:"week_revenue"`
	MonthRevenue       float64 `json:"month_revenue"`
	ActiveWorkOrders   int     `json:"active_work_orders"`
	WaitingQueue       int     `json:"waiting_queue"`
	InProgress         int     `json:"in_progress"`
	CompletedToday     int     `json:"completed_today"`
	LowStockItems      int     `json:"low_stock_items"`
	TotalCustomers     int     `json:"total_customers"`
	TotalSpareparts    int     `json:"total_spareparts"`
}

type RevenueChart struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}
