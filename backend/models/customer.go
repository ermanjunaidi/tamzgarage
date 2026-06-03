package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	BranchID      *uuid.UUID `json:"branch_id"`
	Code          *string    `json:"code"`
	FullName      string     `json:"full_name"`
	Phone         string     `json:"phone"`
	Email         *string    `json:"email"`
	Address       *string    `json:"address"`
	IDNumber      *string    `json:"id_number"`
	Category      string     `json:"category"`
	LoyaltyPoints int        `json:"loyalty_points"`
	Tags          []string   `json:"tags"`
	Notes         *string    `json:"notes"`
	IsBlacklisted bool       `json:"is_blacklisted"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	VehicleCount  int        `json:"vehicle_count,omitempty"`
}

type Vehicle struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	CustomerID    *uuid.UUID `json:"customer_id"`
	PlateNumber   string     `json:"plate_number"`
	Brand         string     `json:"brand"`
	Model         *string    `json:"model"`
	Variant       *string    `json:"variant"`
	Year          *int       `json:"year"`
	Color         *string    `json:"color"`
	VIN           *string    `json:"vin"`
	EngineNumber  *string    `json:"engine_number"`
	LastKM        *int       `json:"last_km"`
	NextServiceKM *int       `json:"next_service_km"`
	Notes         *string    `json:"notes"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CustomerName  *string    `json:"customer_name,omitempty"`
}
