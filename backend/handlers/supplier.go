package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"bengkelpro-backend/middleware"
	"bengkelpro-backend/models"
)

type SupplierHandler struct {
	DB *pgxpool.Pool
}

func NewSupplierHandler(db *pgxpool.Pool) *SupplierHandler {
	return &SupplierHandler{DB: db}
}

func (h *SupplierHandler) List(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	search := c.Query("search", "")

	query := `SELECT id, tenant_id, code, name, contact_person, phone, email, address, tax_id, bank_account, notes, is_active, created_at FROM suppliers WHERE tenant_id=$1`
	args := []interface{}{tenantID}

	if search != "" {
		query += ` AND (name ILIKE $2 OR contact_person ILIKE $2 OR phone ILIKE $2)`
		args = append(args, "%"+search+"%")
	}

	query += ` ORDER BY name`

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	suppliers := []models.Supplier{}
	for rows.Next() {
		var s models.Supplier
		rows.Scan(&s.ID, &s.TenantID, &s.Code, &s.Name, &s.ContactPerson, &s.Phone,
			&s.Email, &s.Address, &s.TaxID, &s.BankAccount, &s.Notes, &s.IsActive, &s.CreatedAt)
		suppliers = append(suppliers, s)
	}

	return c.JSON(suppliers)
}

func (h *SupplierHandler) Get(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	var s models.Supplier
	err := h.DB.QueryRow(context.Background(), `
		SELECT id, tenant_id, code, name, contact_person, phone, email, address, tax_id, bank_account, notes, is_active, created_at
		FROM suppliers WHERE id=$1 AND tenant_id=$2
	`, id, tenantID).Scan(&s.ID, &s.TenantID, &s.Code, &s.Name, &s.ContactPerson, &s.Phone,
		&s.Email, &s.Address, &s.TaxID, &s.BankAccount, &s.Notes, &s.IsActive, &s.CreatedAt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Supplier not found"})
	}

	return c.JSON(s)
}

func (h *SupplierHandler) Create(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var s models.Supplier
	if err := c.BodyParser(&s); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := h.DB.QueryRow(context.Background(), `
		INSERT INTO suppliers (tenant_id, code, name, contact_person, phone, email, address, tax_id, bank_account, notes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id, code, name, created_at
	`, tenantID, s.Code, s.Name, s.ContactPerson, s.Phone, s.Email, s.Address, s.TaxID, s.BankAccount, s.Notes).Scan(&s.ID, &s.Code, &s.Name, &s.CreatedAt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(s)
}

func (h *SupplierHandler) Update(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	var s models.Supplier
	if err := c.BodyParser(&s); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	_, err := h.DB.Exec(context.Background(), `
		UPDATE suppliers SET code=$1, name=$2, contact_person=$3, phone=$4, email=$5,
		address=$6, tax_id=$7, bank_account=$8, notes=$9, updated_at=NOW()
		WHERE id=$10 AND tenant_id=$11
	`, s.Code, s.Name, s.ContactPerson, s.Phone, s.Email, s.Address, s.TaxID, s.BankAccount, s.Notes, id, tenantID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Updated"})
}

func (h *SupplierHandler) Delete(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	_, err := h.DB.Exec(context.Background(), `DELETE FROM suppliers WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Deleted"})
}

func (h *SupplierHandler) CreatePO(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)

	var input struct {
		SupplierID string            `json:"supplier_id"`
		Items      []models.POItem   `json:"items"`
		Notes      string            `json:"notes"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	var totalAmount float64
	for _, item := range input.Items {
		totalAmount += item.TotalPrice
	}

	var poNum string
	h.DB.QueryRow(context.Background(), `SELECT 'PO-'||TO_CHAR(CURRENT_DATE,'YYYYMMDD')||'-'||LPAD(COALESCE(COUNT(*)+1,1)::text,4,'0') FROM purchase_orders WHERE tenant_id=$1 AND DATE(created_at)=CURRENT_DATE`, tenantID).Scan(&poNum)

	tx, _ := h.DB.Begin(context.Background())
	defer tx.Rollback(context.Background())

	var po models.PurchaseOrder
	err := tx.QueryRow(context.Background(), `
		INSERT INTO purchase_orders (tenant_id, po_number, supplier_id, total_amount, notes, created_by)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, po_number, total_amount, created_at
	`, tenantID, poNum, input.SupplierID, totalAmount, input.Notes, userID).Scan(&po.ID, &po.PONumber, &po.TotalAmount, &po.CreatedAt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	for _, item := range input.Items {
		tx.Exec(context.Background(), `INSERT INTO po_items (purchase_order_id, sparepart_id, item_name, quantity, unit_price, total_price) VALUES ($1,$2,$3,$4,$5,$6)`,
			po.ID, item.SparepartID, item.ItemName, item.Quantity, item.UnitPrice, item.TotalPrice)
	}

	tx.Commit(context.Background())
	return c.Status(201).JSON(po)
}

func (h *SupplierHandler) ListPOs(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	status := c.Query("status", "")

	query := `SELECT po.id, po.po_number, po.supplier_id, po.status, po.total_amount, po.notes, po.created_at, s.name as supplier_name
		FROM purchase_orders po JOIN suppliers s ON po.supplier_id = s.id
		WHERE po.tenant_id=$1`
	args := []interface{}{tenantID}

	if status != "" {
		query += ` AND po.status=$2`
		args = append(args, status)
	}

	query += ` ORDER BY po.created_at DESC LIMIT 50`

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	pos := []models.PurchaseOrder{}
	for rows.Next() {
		var po models.PurchaseOrder
		rows.Scan(&po.ID, &po.PONumber, &po.SupplierID, &po.Status, &po.TotalAmount, &po.Notes, &po.CreatedAt, &po.SupplierName)
		pos = append(pos, po)
	}

	return c.JSON(pos)
}

func (h *SupplierHandler) ReceivePO(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	tx, _ := h.DB.Begin(context.Background())
	defer tx.Rollback(context.Background())

	// Update PO status
	tx.Exec(context.Background(), `UPDATE purchase_orders SET status='diterima', received_at=NOW(), updated_at=NOW() WHERE id=$1 AND tenant_id=$2`, id, tenantID)

	// Add stock for each item
	rows, _ := tx.Query(context.Background(), `SELECT sparepart_id, quantity FROM po_items WHERE purchase_order_id=$1`, id)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var spID string
			var qty float64
			rows.Scan(&spID, &qty)
			if spID != "" {
				tx.Exec(context.Background(), `UPDATE spareparts SET current_stock = current_stock + $1 WHERE id=$2 AND tenant_id=$3`, qty, spID, tenantID)
				tx.Exec(context.Background(), `INSERT INTO stock_mutations (tenant_id, sparepart_id, mutation_type, quantity, reference_type, reference_id, notes) VALUES ($1,$2,'masuk',$3,'PO',$4,'Receiving PO')`, tenantID, spID, qty, id)
			}
		}
	}

	tx.Commit(context.Background())
	return c.JSON(fiber.Map{"message": "PO received, stock updated"})
}
