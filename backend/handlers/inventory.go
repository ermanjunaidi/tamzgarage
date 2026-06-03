package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"bengkelpro-backend/middleware"
	"bengkelpro-backend/models"
)

type InventoryHandler struct {
	DB *pgxpool.Pool
}

func NewInventoryHandler(db *pgxpool.Pool) *InventoryHandler {
	return &InventoryHandler{DB: db}
}

func (h *InventoryHandler) ListSpareparts(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	search := c.Query("search", "")
	category := c.Query("category", "")
	lowStock := c.Query("low_stock", "")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	offset := (page - 1) * limit

	query := `SELECT id, tenant_id, branch_id, code, name, category, brand, unit,
		purchase_price, selling_price, wholesale_price, member_price,
		current_stock, min_stock, max_stock, barcode, location, is_active, created_at, updated_at
		FROM spareparts WHERE tenant_id = $1 AND is_active = true`
	args := []interface{}{tenantID}
	argN := 2

	if search != "" {
		query += ` AND (name ILIKE $` + itoa(argN) + ` OR code ILIKE $` + itoa(argN) + ` OR brand ILIKE $` + itoa(argN) + `)`
		args = append(args, "%"+search+"%")
		argN++
	}
	if category != "" {
		query += ` AND category = $` + itoa(argN)
		args = append(args, category)
		argN++
	}
	if lowStock == "true" {
		query += ` AND current_stock <= min_stock`
	}

	query += ` ORDER BY name LIMIT $` + itoa(argN) + ` OFFSET $` + itoa(argN+1)
	args = append(args, limit, offset)

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	items := []models.Sparepart{}
	for rows.Next() {
		var s models.Sparepart
		if err := rows.Scan(&s.ID, &s.TenantID, &s.BranchID, &s.Code, &s.Name, &s.Category,
			&s.Brand, &s.Unit, &s.PurchasePrice, &s.SellingPrice, &s.WholesalePrice,
			&s.MemberPrice, &s.CurrentStock, &s.MinStock, &s.MaxStock, &s.Barcode,
			&s.Location, &s.IsActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
			continue
		}
		items = append(items, s)
	}

	return c.JSON(fiber.Map{"data": items, "page": page, "limit": limit})
}

func (h *InventoryHandler) GetSparepart(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	var s models.Sparepart
	err := h.DB.QueryRow(context.Background(), `
		SELECT id, tenant_id, branch_id, code, name, category, brand, unit,
		purchase_price, selling_price, wholesale_price, member_price,
		current_stock, min_stock, max_stock, barcode, location, is_active, created_at, updated_at
		FROM spareparts WHERE id=$1 AND tenant_id=$2
	`, id, tenantID).Scan(&s.ID, &s.TenantID, &s.BranchID, &s.Code, &s.Name, &s.Category,
		&s.Brand, &s.Unit, &s.PurchasePrice, &s.SellingPrice, &s.WholesalePrice,
		&s.MemberPrice, &s.CurrentStock, &s.MinStock, &s.MaxStock, &s.Barcode,
		&s.Location, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Sparepart not found"})
	}

	return c.JSON(s)
}

func (h *InventoryHandler) CreateSparepart(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var s models.Sparepart
	if err := c.BodyParser(&s); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := h.DB.QueryRow(context.Background(), `
		INSERT INTO spareparts (tenant_id, branch_id, code, name, category, brand, unit,
		purchase_price, selling_price, wholesale_price, member_price,
		current_stock, min_stock, max_stock, barcode, location)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		RETURNING id, code, name, category, current_stock, selling_price, created_at
	`, tenantID, s.BranchID, s.Code, s.Name, s.Category, s.Brand, s.Unit,
		s.PurchasePrice, s.SellingPrice, s.WholesalePrice, s.MemberPrice,
		s.CurrentStock, s.MinStock, s.MaxStock, s.Barcode, s.Location).Scan(
		&s.ID, &s.Code, &s.Name, &s.Category, &s.CurrentStock, &s.SellingPrice, &s.CreatedAt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(s)
}

func (h *InventoryHandler) UpdateSparepart(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	var s models.Sparepart
	if err := c.BodyParser(&s); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	_, err := h.DB.Exec(context.Background(), `
		UPDATE spareparts SET code=$1, name=$2, category=$3, brand=$4, unit=$5,
		purchase_price=$6, selling_price=$7, wholesale_price=$8, member_price=$9,
		min_stock=$10, max_stock=$11, barcode=$12, location=$13, updated_at=NOW()
		WHERE id=$14 AND tenant_id=$15
	`, s.Code, s.Name, s.Category, s.Brand, s.Unit, s.PurchasePrice, s.SellingPrice,
		s.WholesalePrice, s.MemberPrice, s.MinStock, s.MaxStock, s.Barcode, s.Location, id, tenantID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Updated"})
}

func (h *InventoryHandler) DeleteSparepart(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	_, err := h.DB.Exec(context.Background(), `UPDATE spareparts SET is_active=false WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Deleted"})
}

func (h *InventoryHandler) AddStock(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)

	var input struct {
		SparepartID string  `json:"sparepart_id"`
		Quantity    float64 `json:"quantity"`
		Notes       string  `json:"notes"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	tx, _ := h.DB.Begin(context.Background())
	defer tx.Rollback(context.Background())

	tx.Exec(context.Background(), `UPDATE spareparts SET current_stock = current_stock + $1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3`,
		input.Quantity, input.SparepartID, tenantID)

	tx.Exec(context.Background(), `INSERT INTO stock_mutations (tenant_id, sparepart_id, mutation_type, quantity, notes, created_by)
		VALUES ($1,$2,'masuk',$3,$4,$5)`, tenantID, input.SparepartID, input.Quantity, input.Notes, userID)

	tx.Commit(context.Background())

	return c.JSON(fiber.Map{"message": "Stock added"})
}

func (h *InventoryHandler) ListStockMutations(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	sparepartID := c.Query("sparepart_id", "")

	query := `SELECT sm.id, sm.tenant_id, sm.branch_id, sm.sparepart_id, sm.mutation_type,
		sm.quantity, sm.reference_type, sm.reference_id, sm.notes, sm.created_at,
		sp.name as sparepart_name
		FROM stock_mutations sm JOIN spareparts sp ON sm.sparepart_id = sp.id
		WHERE sm.tenant_id = $1`
	args := []interface{}{tenantID}

	if sparepartID != "" {
		query += ` AND sm.sparepart_id = $2`
		args = append(args, sparepartID)
	}

	query += ` ORDER BY sm.created_at DESC LIMIT 100`

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	mutations := []models.StockMutation{}
	for rows.Next() {
		var m models.StockMutation
		if err := rows.Scan(&m.ID, &m.TenantID, &m.BranchID, &m.SparepartID, &m.MutationType,
			&m.Quantity, &m.ReferenceType, &m.ReferenceID, &m.Notes, &m.CreatedAt, &m.SparepartName); err != nil {
			continue
		}
		mutations = append(mutations, m)
	}

	return c.JSON(mutations)
}
