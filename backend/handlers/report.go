package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"bengkelpro-backend/middleware"
)

type ReportHandler struct {
	DB *pgxpool.Pool
}

func NewReportHandler(db *pgxpool.Pool) *ReportHandler {
	return &ReportHandler{DB: db}
}

func (h *ReportHandler) RevenueReport(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	branchID := middleware.GetBranchID(c)
	startDate := c.Query("start", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.Query("end", time.Now().Format("2006-01-02"))

	rows, err := h.DB.Query(context.Background(), `
		SELECT DATE(i.paid_at), COUNT(*), COALESCE(SUM(i.grand_total),0), COALESCE(SUM(i.amount_paid),0)
		FROM invoices i
		WHERE i.tenant_id=$1 AND ($2='' OR i.branch_id=$2::uuid)
		AND i.paid_at BETWEEN $3 AND $4 AND i.status='lunas'
		GROUP BY DATE(i.paid_at) ORDER BY DATE(i.paid_at)
	`, tenantID, branchID, startDate, endDate)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	type RevenueRow struct {
		Date      string  `json:"date"`
		Count     int     `json:"count"`
		Total     float64 `json:"total"`
		Paid      float64 `json:"paid"`
	}

	report := []RevenueRow{}
	var totalRev, totalPaid float64
	var totalCount int

	for rows.Next() {
		var r RevenueRow
		rows.Scan(&r.Date, &r.Count, &r.Total, &r.Paid)
		report = append(report, r)
		totalRev += r.Total
		totalPaid += r.Paid
		totalCount += r.Count
	}

	return c.JSON(fiber.Map{
		"data":        report,
		"total_count": totalCount,
		"total_rev":   totalRev,
		"total_paid":  totalPaid,
	})
}

func (h *ReportHandler) StockReport(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	rows, err := h.DB.Query(context.Background(), `
		SELECT code, name, category, brand, current_stock, min_stock, selling_price,
			(current_stock * selling_price) as stock_value,
			CASE WHEN current_stock <= min_stock THEN true ELSE false END as low_stock
		FROM spareparts WHERE tenant_id=$1 AND is_active=true
		ORDER BY category, name
	`, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	type StockRow struct {
		Code        string  `json:"code"`
		Name        string  `json:"name"`
		Category    string  `json:"category"`
		Brand       *string `json:"brand"`
		CurrentStock float64 `json:"current_stock"`
		MinStock    float64 `json:"min_stock"`
		SellingPrice float64 `json:"selling_price"`
		StockValue  float64 `json:"stock_value"`
		LowStock    bool    `json:"low_stock"`
	}

	items := []StockRow{}
	var totalValue float64
	for rows.Next() {
		var s StockRow
		rows.Scan(&s.Code, &s.Name, &s.Category, &s.Brand, &s.CurrentStock, &s.MinStock, &s.SellingPrice, &s.StockValue, &s.LowStock)
		items = append(items, s)
		totalValue += s.StockValue
	}

	return c.JSON(fiber.Map{"data": items, "total_value": totalValue})
}

func (h *ReportHandler) WorkOrderReport(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	branchID := middleware.GetBranchID(c)

	rows, err := h.DB.Query(context.Background(), `
		SELECT status, COUNT(*), COALESCE(SUM(grand_total),0)
		FROM work_orders WHERE tenant_id=$1 AND ($2='' OR branch_id=$2::uuid)
		GROUP BY status ORDER BY status
	`, tenantID, branchID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	type WOStatusRow struct {
		Status string  `json:"status"`
		Count  int     `json:"count"`
		Value  float64 `json:"value"`
	}

	report := []WOStatusRow{}
	for rows.Next() {
		var r WOStatusRow
		rows.Scan(&r.Status, &r.Count, &r.Value)
		report = append(report, r)
	}

	return c.JSON(report)
}

func (h *ReportHandler) TopCustomers(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	rows, err := h.DB.Query(context.Background(), `
		SELECT c.full_name, c.phone, COUNT(wo.id) as wo_count, COALESCE(SUM(i.grand_total),0) as total_spent
		FROM customers c
		LEFT JOIN work_orders wo ON c.id=wo.customer_id
		LEFT JOIN invoices i ON wo.id=i.work_order_id AND i.status='lunas'
		WHERE c.tenant_id=$1
		GROUP BY c.id, c.full_name, c.phone
		ORDER BY total_spent DESC LIMIT 10
	`, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	type CustomerRow struct {
		FullName   string  `json:"full_name"`
		Phone      string  `json:"phone"`
		WOCount    int     `json:"wo_count"`
		TotalSpent float64 `json:"total_spent"`
	}

	customers := []CustomerRow{}
	for rows.Next() {
		var cu CustomerRow
		rows.Scan(&cu.FullName, &cu.Phone, &cu.WOCount, &cu.TotalSpent)
		customers = append(customers, cu)
	}

	return c.JSON(customers)
}
