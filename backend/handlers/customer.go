package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"bengkelpro-backend/middleware"
	"bengkelpro-backend/models"
)

type CustomerHandler struct {
	DB *pgxpool.Pool
}

func NewCustomerHandler(db *pgxpool.Pool) *CustomerHandler {
	return &CustomerHandler{DB: db}
}

func (h *CustomerHandler) List(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	search := c.Query("search", "")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	offset := (page - 1) * limit

	query := `
		SELECT c.id, c.tenant_id, c.branch_id, c.code, c.full_name, c.phone, c.email,
		       c.address, c.id_number, c.category, c.loyalty_points, c.tags, c.notes,
		       c.is_blacklisted, c.created_at, c.updated_at,
		       COALESCE((SELECT COUNT(*) FROM vehicles v WHERE v.customer_id = c.id), 0) as vehicle_count
		FROM customers c WHERE c.tenant_id = $1
	`
	args := []interface{}{tenantID}

	if search != "" {
		query += ` AND (c.full_name ILIKE $2 OR c.phone ILIKE $2 OR c.code ILIKE $2)`
		args = append(args, "%"+search+"%")
	}

	query += ` ORDER BY c.created_at DESC LIMIT $` + itoa(len(args)+1) + ` OFFSET $` + itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	customers := []models.Customer{}
	for rows.Next() {
		var cu models.Customer
		if err := rows.Scan(&cu.ID, &cu.TenantID, &cu.BranchID, &cu.Code, &cu.FullName, &cu.Phone,
			&cu.Email, &cu.Address, &cu.IDNumber, &cu.Category, &cu.LoyaltyPoints,
			&cu.Tags, &cu.Notes, &cu.IsBlacklisted, &cu.CreatedAt, &cu.UpdatedAt, &cu.VehicleCount); err != nil {
			continue
		}
		customers = append(customers, cu)
	}

	var total int
	countQ := `SELECT COUNT(*) FROM customers WHERE tenant_id = $1`
	h.DB.QueryRow(context.Background(), countQ, tenantID).Scan(&total)

	return c.JSON(fiber.Map{"data": customers, "total": total, "page": page, "limit": limit})
}

func (h *CustomerHandler) Get(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	var cu models.Customer
	err := h.DB.QueryRow(context.Background(), `
		SELECT c.id, c.tenant_id, c.branch_id, c.code, c.full_name, c.phone, c.email,
		       c.address, c.id_number, c.category, c.loyalty_points, c.tags, c.notes,
		       c.is_blacklisted, c.created_at, c.updated_at,
		       COALESCE((SELECT COUNT(*) FROM vehicles v WHERE v.customer_id = c.id), 0) as vehicle_count
		FROM customers c WHERE c.id = $1 AND c.tenant_id = $2
	`, id, tenantID).Scan(&cu.ID, &cu.TenantID, &cu.BranchID, &cu.Code, &cu.FullName, &cu.Phone,
		&cu.Email, &cu.Address, &cu.IDNumber, &cu.Category, &cu.LoyaltyPoints,
		&cu.Tags, &cu.Notes, &cu.IsBlacklisted, &cu.CreatedAt, &cu.UpdatedAt, &cu.VehicleCount)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Customer not found"})
	}

	return c.JSON(cu)
}

func (h *CustomerHandler) Create(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var cu models.Customer
	if err := c.BodyParser(&cu); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := h.DB.QueryRow(context.Background(), `
		INSERT INTO customers (tenant_id, branch_id, code, full_name, phone, email, address, id_number, category, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, tenant_id, code, full_name, phone, email, address, id_number, category, notes, is_blacklisted, created_at
	`, tenantID, cu.BranchID, cu.Code, cu.FullName, cu.Phone, cu.Email, cu.Address,
		cu.IDNumber, cu.Category, cu.Notes).Scan(
		&cu.ID, &cu.TenantID, &cu.Code, &cu.FullName, &cu.Phone,
		&cu.Email, &cu.Address, &cu.IDNumber, &cu.Category, &cu.Notes,
		&cu.IsBlacklisted, &cu.CreatedAt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(cu)
}

func (h *CustomerHandler) Update(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	var cu models.Customer
	if err := c.BodyParser(&cu); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	_, err := h.DB.Exec(context.Background(), `
		UPDATE customers SET full_name=$1, phone=$2, email=$3, address=$4, id_number=$5,
		category=$6, notes=$7, is_blacklisted=$8, updated_at=NOW()
		WHERE id=$9 AND tenant_id=$10
	`, cu.FullName, cu.Phone, cu.Email, cu.Address, cu.IDNumber, cu.Category,
		cu.Notes, cu.IsBlacklisted, id, tenantID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Updated"})
}

func (h *CustomerHandler) Delete(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	_, err := h.DB.Exec(context.Background(), `DELETE FROM customers WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Deleted"})
}
