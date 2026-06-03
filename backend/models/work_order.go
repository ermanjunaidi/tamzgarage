package models

import (
	"time"

	"github.com/google/uuid"
)

type WorkOrder struct {
	ID                   uuid.UUID  `json:"id"`
	TenantID             uuid.UUID  `json:"tenant_id"`
	BranchID             *uuid.UUID `json:"branch_id"`
	WONumber             string     `json:"wo_number"`
	CustomerID           *uuid.UUID `json:"customer_id"`
	VehicleID            *uuid.UUID `json:"vehicle_id"`
	Status               string     `json:"status"`
	Complaint            string     `json:"complaint"`
	Diagnosis            *string    `json:"diagnosis"`
	MechanicNotes        *string    `json:"mechanic_notes"`
	EstimatedCost        float64    `json:"estimated_cost"`
	EstimatedHours       *float64   `json:"estimated_hours"`
	ActualCost           float64    `json:"actual_cost"`
	LaborCost            float64    `json:"labor_cost"`
	TotalSparepartCost   float64    `json:"total_sparepart_cost"`
	Discount             float64    `json:"discount"`
	TaxAmount            float64    `json:"tax_amount"`
	GrandTotal           float64    `json:"grand_total"`
	CustomerApproved     bool       `json:"customer_approved"`
	QueueNumber          *int       `json:"queue_number"`
	CreatedBy            *uuid.UUID `json:"created_by"`
	StartedAt            *time.Time `json:"started_at"`
	CompletedAt          *time.Time `json:"completed_at"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	CustomerName         *string    `json:"customer_name,omitempty"`
	VehiclePlate         *string    `json:"vehicle_plate,omitempty"`
	Services             []WOService    `json:"services,omitempty"`
	Spareparts           []WOSparepart  `json:"spareparts,omitempty"`
	Mechanics            []WOMechanic   `json:"mechanics,omitempty"`
}

type WOService struct {
	ID          uuid.UUID `json:"id"`
	WorkOrderID uuid.UUID `json:"work_order_id"`
	ServiceName string    `json:"service_name"`
	Description *string   `json:"description"`
	Quantity    float64   `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
	CreatedAt   time.Time `json:"created_at"`
}

type WOSparepart struct {
	ID            uuid.UUID `json:"id"`
	WorkOrderID   uuid.UUID `json:"work_order_id"`
	SparepartID   uuid.UUID `json:"sparepart_id"`
	SparepartName string    `json:"sparepart_name"`
	Quantity      float64   `json:"quantity"`
	UnitPrice     float64   `json:"unit_price"`
	TotalPrice    float64   `json:"total_price"`
	CreatedAt     time.Time `json:"created_at"`
}

type WOMechanic struct {
	ID          uuid.UUID `json:"id"`
	WorkOrderID uuid.UUID `json:"work_order_id"`
	UserID      uuid.UUID `json:"user_id"`
	Commission  float64   `json:"commission"`
	AssignedAt  time.Time `json:"assigned_at"`
	FullName    *string   `json:"full_name,omitempty"`
}
