package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"bengkelpro-backend/middleware"
	"bengkelpro-backend/models"
)

type POSHandler struct {
	DB *pgxpool.Pool
}

func NewPOSHandler(db *pgxpool.Pool) *POSHandler {
	return &POSHandler{DB: db}
}

func (h *POSHandler) CreateInvoice(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)

	var input struct {
		WorkOrderID string  `json:"work_order_id"`
		CustomerID  string  `json:"customer_id"`
		Subtotal    float64 `json:"subtotal"`
		Discount    float64 `json:"discount"`
		TaxAmount   float64 `json:"tax_amount"`
		GrandTotal  float64 `json:"grand_total"`
		Notes       string  `json:"notes"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	var inv models.Invoice
	var invNum string
	h.DB.QueryRow(context.Background(), `SELECT 'INV-'||TO_CHAR(CURRENT_DATE,'YYYYMMDD')||'-'||LPAD(COALESCE(COUNT(*)+1,1)::text,4,'0') FROM invoices WHERE tenant_id=$1 AND DATE(created_at)=CURRENT_DATE`, tenantID).Scan(&invNum)

	err := h.DB.QueryRow(context.Background(), `
		INSERT INTO invoices (tenant_id, branch_id, invoice_number, work_order_id, customer_id, subtotal, discount, tax_amount, grand_total, balance_due, notes, created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		RETURNING id, invoice_number, subtotal, discount, tax_amount, grand_total, amount_paid, balance_due, status, created_at
	`, tenantID, middleware.GetBranchID(c), invNum, input.WorkOrderID, input.CustomerID,
		input.Subtotal, input.Discount, input.TaxAmount, input.GrandTotal, input.GrandTotal, input.Notes, userID).Scan(
		&inv.ID, &inv.InvoiceNumber, &inv.Subtotal, &inv.Discount,
		&inv.TaxAmount, &inv.GrandTotal, &inv.AmountPaid, &inv.BalanceDue, &inv.Status, &inv.CreatedAt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(inv)
}

func (h *POSHandler) ProcessPayment(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)

	var input struct {
		InvoiceID string  `json:"invoice_id"`
		Amount    float64 `json:"amount"`
		Method    string  `json:"method"`
		Reference string  `json:"reference"`
		Notes     string  `json:"notes"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	tx, _ := h.DB.Begin(context.Background())
	defer tx.Rollback(context.Background())

	// Insert payment
	_, err := tx.Exec(context.Background(), `
		INSERT INTO payments (tenant_id, branch_id, invoice_id, amount, method, reference, notes, created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, tenantID, middleware.GetBranchID(c), input.InvoiceID, input.Amount, input.Method, input.Reference, input.Notes, userID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Update invoice
	var amountPaid, grandTotal float64
	tx.QueryRow(context.Background(), `SELECT COALESCE(SUM(amount),0) FROM payments WHERE invoice_id=$1`, input.InvoiceID).Scan(&amountPaid)
	tx.QueryRow(context.Background(), `SELECT grand_total FROM invoices WHERE id=$1`, input.InvoiceID).Scan(&grandTotal)

	status := "belum_bayar"
	if amountPaid >= grandTotal {
		status = "lunas"
	} else if amountPaid > 0 {
		status = "dp"
	}

	tx.Exec(context.Background(), `UPDATE invoices SET amount_paid=$1, balance_due=grand_total-$1, status=$2, paid_at=CASE WHEN $2='lunas' THEN NOW() ELSE paid_at END, updated_at=NOW() WHERE id=$3`,
		amountPaid, status, input.InvoiceID)

	tx.Commit(context.Background())

	return c.JSON(fiber.Map{"message": "Payment processed", "status": status, "amount_paid": amountPaid})
}

func (h *POSHandler) ListInvoices(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	status := c.Query("status", "")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	offset := (page - 1) * limit

	query := `SELECT i.id, i.invoice_number, i.work_order_id, i.customer_id, i.subtotal,
		i.discount, i.tax_amount, i.grand_total, i.amount_paid, i.balance_due, i.status,
		i.created_at, i.paid_at, c.full_name as customer_name
		FROM invoices i LEFT JOIN customers c ON i.customer_id = c.id
		WHERE i.tenant_id=$1`
	args := []interface{}{tenantID}
	argN := 2

	if status != "" {
		query += ` AND i.status=$` + itoa(argN)
		args = append(args, status)
		argN++
	}

	query += ` ORDER BY i.created_at DESC LIMIT $` + itoa(argN) + ` OFFSET $` + itoa(argN+1)
	args = append(args, limit, offset)

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	invoices := []models.Invoice{}
	for rows.Next() {
		var i models.Invoice
		if err := rows.Scan(&i.ID, &i.InvoiceNumber, &i.WorkOrderID, &i.CustomerID,
			&i.Subtotal, &i.Discount, &i.TaxAmount, &i.GrandTotal, &i.AmountPaid,
			&i.BalanceDue, &i.Status, &i.CreatedAt, &i.PaidAt, &i.CustomerName); err != nil {
			continue
		}
		invoices = append(invoices, i)
	}

	return c.JSON(fiber.Map{"data": invoices, "page": page, "limit": limit})
}

func (h *POSHandler) GetInvoice(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	var i models.Invoice
	err := h.DB.QueryRow(context.Background(), `
		SELECT i.id, i.invoice_number, i.work_order_id, i.customer_id, i.subtotal,
		i.discount, i.tax_amount, i.grand_total, i.amount_paid, i.balance_due, i.status,
		i.notes, i.created_at, i.paid_at, c.full_name as customer_name
		FROM invoices i LEFT JOIN customers c ON i.customer_id = c.id
		WHERE i.id=$1 AND i.tenant_id=$2
	`, id, tenantID).Scan(&i.ID, &i.InvoiceNumber, &i.WorkOrderID, &i.CustomerID,
		&i.Subtotal, &i.Discount, &i.TaxAmount, &i.GrandTotal, &i.AmountPaid,
		&i.BalanceDue, &i.Status, &i.Notes, &i.CreatedAt, &i.PaidAt, &i.CustomerName)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Invoice not found"})
	}

	// Load payments
	pRows, _ := h.DB.Query(context.Background(), `SELECT id, amount, method, reference, notes, created_at FROM payments WHERE invoice_id=$1`, id)
	if pRows != nil {
		defer pRows.Close()
		for pRows.Next() {
			var p models.Payment
			pRows.Scan(&p.ID, &p.Amount, &p.Method, &p.Reference, &p.Notes, &p.CreatedAt)
			i.Payments = append(i.Payments, p)
		}
	}

	return c.JSON(i)
}
