package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"bengkelpro-backend/middleware"
	"bengkelpro-backend/models"
)

type DashboardHandler struct {
	DB *pgxpool.Pool
}

func NewDashboardHandler(db *pgxpool.Pool) *DashboardHandler {
	return &DashboardHandler{DB: db}
}

func (h *DashboardHandler) GetStats(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	branchID := middleware.GetBranchID(c)

	stats := models.DashboardStats{}

	h.DB.QueryRow(context.Background(), `
		SELECT COALESCE(COUNT(*),0) FROM work_orders
		WHERE tenant_id=$1 AND ($2='' OR branch_id=$2::uuid) AND DATE(created_at)=CURRENT_DATE
	`, tenantID, branchID).Scan(&stats.TodayVehiclesIn)

	h.DB.QueryRow(context.Background(), `
		SELECT COALESCE(SUM(grand_total),0) FROM invoices
		WHERE tenant_id=$1 AND ($2='' OR branch_id=$2::uuid) AND DATE(paid_at)=CURRENT_DATE AND status='lunas'
	`, tenantID, branchID).Scan(&stats.TodayRevenue)

	h.DB.QueryRow(context.Background(), `
		SELECT COALESCE(SUM(grand_total),0) FROM invoices
		WHERE tenant_id=$1 AND ($2='' OR branch_id=$2::uuid) AND paid_at >= CURRENT_DATE - INTERVAL '7 days' AND status='lunas'
	`, tenantID, branchID).Scan(&stats.WeekRevenue)

	h.DB.QueryRow(context.Background(), `
		SELECT COALESCE(SUM(grand_total),0) FROM invoices
		WHERE tenant_id=$1 AND ($2='' OR branch_id=$2::uuid) AND paid_at >= DATE_TRUNC('month', CURRENT_DATE) AND status='lunas'
	`, tenantID, branchID).Scan(&stats.MonthRevenue)

	h.DB.QueryRow(context.Background(), `
		SELECT COALESCE(COUNT(*),0) FROM work_orders
		WHERE tenant_id=$1 AND ($2='' OR branch_id=$2::uuid) AND status IN ('menunggu','diagnosis','menunggu_persetujuan','dikerjakan','qc')
	`, tenantID, branchID).Scan(&stats.ActiveWorkOrders)

	h.DB.QueryRow(context.Background(), `
		SELECT COALESCE(COUNT(*),0) FROM work_orders
		WHERE tenant_id=$1 AND ($2='' OR branch_id=$2::uuid) AND status='menunggu'
	`, tenantID, branchID).Scan(&stats.WaitingQueue)

	h.DB.QueryRow(context.Background(), `
		SELECT COALESCE(COUNT(*),0) FROM work_orders
		WHERE tenant_id=$1 AND ($2='' OR branch_id=$2::uuid) AND status='dikerjakan'
	`, tenantID, branchID).Scan(&stats.InProgress)

	h.DB.QueryRow(context.Background(), `
		SELECT COALESCE(COUNT(*),0) FROM work_orders
		WHERE tenant_id=$1 AND ($2='' OR branch_id=$2::uuid) AND DATE(completed_at)=CURRENT_DATE
	`, tenantID, branchID).Scan(&stats.CompletedToday)

	h.DB.QueryRow(context.Background(), `
		SELECT COALESCE(COUNT(*),0) FROM spareparts
		WHERE tenant_id=$1 AND ($2='' OR branch_id=$2::uuid) AND current_stock <= min_stock AND is_active=true
	`, tenantID, branchID).Scan(&stats.LowStockItems)

	h.DB.QueryRow(context.Background(), `
		SELECT COALESCE(COUNT(*),0) FROM customers WHERE tenant_id=$1
	`, tenantID).Scan(&stats.TotalCustomers)

	h.DB.QueryRow(context.Background(), `
		SELECT COALESCE(COUNT(*),0) FROM spareparts WHERE tenant_id=$1 AND is_active=true
	`, tenantID).Scan(&stats.TotalSpareparts)

	return c.JSON(stats)
}

func (h *DashboardHandler) GetRevenueChart(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	branchID := middleware.GetBranchID(c)

	rows, err := h.DB.Query(context.Background(), `
		SELECT TO_CHAR(d.date, 'YYYY-MM-DD'), COALESCE(SUM(i.grand_total), 0)
		FROM generate_series(CURRENT_DATE - INTERVAL '30 days', CURRENT_DATE, '1 day') AS d(date)
		LEFT JOIN invoices i ON DATE(i.paid_at)=d.date
			AND i.tenant_id=$1 AND ($2='' OR i.branch_id=$2::uuid) AND i.status='lunas'
		GROUP BY d.date ORDER BY d.date
	`, tenantID, branchID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	chart := []models.RevenueChart{}
	for rows.Next() {
		var r models.RevenueChart
		if err := rows.Scan(&r.Date, &r.Amount); err != nil {
			continue
		}
		chart = append(chart, r)
	}

	return c.JSON(chart)
}
