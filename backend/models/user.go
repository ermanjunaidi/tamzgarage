package models

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Code      string     `json:"code"`
	Address   *string    `json:"address"`
	Phone     *string    `json:"phone"`
	Email     *string    `json:"email"`
	LogoURL   *string    `json:"logo_url"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Branch struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Address   *string   `json:"address"`
	Phone     *string   `json:"phone"`
	Email     *string   `json:"email"`
	TaxRate   float64   `json:"tax_rate"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     *uuid.UUID `json:"tenant_id"`
	BranchID     *uuid.UUID `json:"branch_id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	FullName     string     `json:"full_name"`
	Role         string     `json:"role"`
	Phone        *string    `json:"phone"`
	Address      *string    `json:"address"`
	Skills       []string   `json:"skills"`
	IsActive     bool       `json:"is_active"`
	LastLogin    *time.Time `json:"last_login"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	User     User   `json:"user"`
	Branch   *Branch `json:"branch"`
}

type AuditLog struct {
	ID         uuid.UUID       `json:"id"`
	TenantID   *uuid.UUID      `json:"tenant_id"`
	UserID     *uuid.UUID      `json:"user_id"`
	Action     string          `json:"action"`
	EntityType string          `json:"entity_type"`
	EntityID   *uuid.UUID      `json:"entity_id"`
	OldData    interface{}     `json:"old_data"`
	NewData    interface{}     `json:"new_data"`
	IPAddress  *string         `json:"ip_address"`
	CreatedAt  time.Time       `json:"created_at"`
}
