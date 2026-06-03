package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"bengkelpro-backend/middleware"
	"bengkelpro-backend/models"
)

type EmployeeHandler struct {
	DB *pgxpool.Pool
}

func NewEmployeeHandler(db *pgxpool.Pool) *EmployeeHandler {
	return &EmployeeHandler{DB: db}
}

func (h *EmployeeHandler) ListMechanics(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	rows, err := h.DB.Query(context.Background(), `
		SELECT u.id, u.full_name, u.phone, u.skills, u.is_active,
			COALESCE((SELECT COUNT(*) FROM wo_mechanics wm JOIN work_orders wo ON wm.work_order_id=wo.id WHERE wm.user_id=u.id AND wo.status IN ('menunggu','diagnosis','menunggu_persetujuan','dikerjakan','qc')),0) as active_jobs
		FROM users u WHERE u.tenant_id=$1 AND u.role='mekanik' ORDER BY u.full_name
	`, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	type MechanicInfo struct {
		ID        string   `json:"id"`
		FullName  string   `json:"full_name"`
		Phone     string   `json:"phone"`
		Skills    []string `json:"skills"`
		IsActive  bool     `json:"is_active"`
		ActiveJobs int     `json:"active_jobs"`
	}

	mechanics := []MechanicInfo{}
	for rows.Next() {
		var m MechanicInfo
		rows.Scan(&m.ID, &m.FullName, &m.Phone, &m.Skills, &m.IsActive, &m.ActiveJobs)
		mechanics = append(mechanics, m)
	}

	return c.JSON(mechanics)
}

func (h *EmployeeHandler) GetEmployee(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	var u models.User
	err := h.DB.QueryRow(context.Background(), `
		SELECT id, tenant_id, branch_id, username, email, full_name, role, phone, skills, is_active, created_at
		FROM users WHERE id=$1 AND tenant_id=$2
	`, id, tenantID).Scan(&u.ID, &u.TenantID, &u.BranchID, &u.Username, &u.Email,
		&u.FullName, &u.Role, &u.Phone, &u.Skills, &u.IsActive, &u.CreatedAt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Employee not found"})
	}

	return c.JSON(u)
}
