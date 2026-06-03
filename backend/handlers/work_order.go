package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"bengkelpro-backend/middleware"
	"bengkelpro-backend/models"
)

type WorkOrderHandler struct {
	DB *pgxpool.Pool
}

func NewWorkOrderHandler(db *pgxpool.Pool) *WorkOrderHandler {
	return &WorkOrderHandler{DB: db}
}

func (h *WorkOrderHandler) List(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	branchID := middleware.GetBranchID(c)
	status := c.Query("status", "")
	search := c.Query("search", "")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	offset := (page - 1) * limit

	query := `
		SELECT w.id, w.tenant_id, w.branch_id, w.wo_number, w.customer_id, w.vehicle_id,
		       w.status, w.complaint, w.diagnosis, w.mechanic_notes, w.estimated_cost,
		       w.estimated_hours, w.actual_cost, w.labor_cost, w.total_sparepart_cost,
		       w.discount, w.tax_amount, w.grand_total, w.customer_approved, w.queue_number,
		       w.created_at, w.updated_at, w.started_at, w.completed_at,
		       cu.full_name as customer_name, v.plate_number as vehicle_plate
		FROM work_orders w
		LEFT JOIN customers cu ON w.customer_id = cu.id
		LEFT JOIN vehicles v ON w.vehicle_id = v.id
		WHERE w.tenant_id = $1
	`
	args := []interface{}{tenantID}
	argN := 2

	if branchID != "" {
		query += ` AND w.branch_id = $` + itoa(argN)
		args = append(args, branchID)
		argN++
	}
	if status != "" {
		query += ` AND w.status = $` + itoa(argN)
		args = append(args, status)
		argN++
	}
	if search != "" {
		query += ` AND (w.wo_number ILIKE $` + itoa(argN) + ` OR cu.full_name ILIKE $` + itoa(argN) + ` OR v.plate_number ILIKE $` + itoa(argN) + `)`
		args = append(args, "%"+search+"%")
		argN++
	}

	query += ` ORDER BY w.created_at DESC LIMIT $` + itoa(argN) + ` OFFSET $` + itoa(argN+1)
	args = append(args, limit, offset)

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	orders := []models.WorkOrder{}
	for rows.Next() {
		var w models.WorkOrder
		if err := rows.Scan(&w.ID, &w.TenantID, &w.BranchID, &w.WONumber, &w.CustomerID,
			&w.VehicleID, &w.Status, &w.Complaint, &w.Diagnosis, &w.MechanicNotes,
			&w.EstimatedCost, &w.EstimatedHours, &w.ActualCost, &w.LaborCost,
			&w.TotalSparepartCost, &w.Discount, &w.TaxAmount, &w.GrandTotal,
			&w.CustomerApproved, &w.QueueNumber, &w.CreatedAt, &w.UpdatedAt,
			&w.StartedAt, &w.CompletedAt, &w.CustomerName, &w.VehiclePlate); err != nil {
			continue
		}
		orders = append(orders, w)
	}

	return c.JSON(fiber.Map{"data": orders, "page": page, "limit": limit})
}

func (h *WorkOrderHandler) Get(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	var w models.WorkOrder
	err := h.DB.QueryRow(context.Background(), `
		SELECT w.id, w.tenant_id, w.branch_id, w.wo_number, w.customer_id, w.vehicle_id,
		       w.status, w.complaint, w.diagnosis, w.mechanic_notes, w.estimated_cost,
		       w.estimated_hours, w.actual_cost, w.labor_cost, w.total_sparepart_cost,
		       w.discount, w.tax_amount, w.grand_total, w.customer_approved, w.queue_number,
		       w.created_at, w.updated_at, w.started_at, w.completed_at,
		       cu.full_name as customer_name, v.plate_number as vehicle_plate
		FROM work_orders w
		LEFT JOIN customers cu ON w.customer_id = cu.id
		LEFT JOIN vehicles v ON w.vehicle_id = v.id
		WHERE w.id = $1 AND w.tenant_id = $2
	`, id, tenantID).Scan(&w.ID, &w.TenantID, &w.BranchID, &w.WONumber, &w.CustomerID,
		&w.VehicleID, &w.Status, &w.Complaint, &w.Diagnosis, &w.MechanicNotes,
		&w.EstimatedCost, &w.EstimatedHours, &w.ActualCost, &w.LaborCost,
		&w.TotalSparepartCost, &w.Discount, &w.TaxAmount, &w.GrandTotal,
		&w.CustomerApproved, &w.QueueNumber, &w.CreatedAt, &w.UpdatedAt,
		&w.StartedAt, &w.CompletedAt, &w.CustomerName, &w.VehiclePlate)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Work order not found"})
	}

	// Load services
	sRows, _ := h.DB.Query(context.Background(), `SELECT id, work_order_id, service_name, description, quantity, unit_price, total_price FROM wo_services WHERE work_order_id=$1`, id)
	if sRows != nil {
		defer sRows.Close()
		for sRows.Next() {
			var s models.WOService
			sRows.Scan(&s.ID, &s.WorkOrderID, &s.ServiceName, &s.Description, &s.Quantity, &s.UnitPrice, &s.TotalPrice)
			w.Services = append(w.Services, s)
		}
	}

	// Load spareparts
	spRows, _ := h.DB.Query(context.Background(), `SELECT id, work_order_id, sparepart_id, sparepart_name, quantity, unit_price, total_price FROM wo_spareparts WHERE work_order_id=$1`, id)
	if spRows != nil {
		defer spRows.Close()
		for spRows.Next() {
			var sp models.WOSparepart
			spRows.Scan(&sp.ID, &sp.WorkOrderID, &sp.SparepartID, &sp.SparepartName, &sp.Quantity, &sp.UnitPrice, &sp.TotalPrice)
			w.Spareparts = append(w.Spareparts, sp)
		}
	}

	// Load mechanics
	mRows, _ := h.DB.Query(context.Background(), `SELECT wm.id, wm.work_order_id, wm.user_id, wm.commission, wm.assigned_at, u.full_name FROM wo_mechanics wm JOIN users u ON wm.user_id=u.id WHERE wm.work_order_id=$1`, id)
	if mRows != nil {
		defer mRows.Close()
		for mRows.Next() {
			var m models.WOMechanic
			mRows.Scan(&m.ID, &m.WorkOrderID, &m.UserID, &m.Commission, &m.AssignedAt, &m.FullName)
			w.Mechanics = append(w.Mechanics, m)
		}
	}

	return c.JSON(w)
}

func (h *WorkOrderHandler) Create(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)

	var input struct {
		models.WorkOrder
		Services   []models.WOService   `json:"services"`
		Spareparts []models.WOSparepart `json:"spareparts"`
		Mechanics  []models.WOMechanic  `json:"mechanics"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Generate WO number
	var woNum string
	h.DB.QueryRow(context.Background(), `SELECT 'WO-'||TO_CHAR(CURRENT_DATE,'YYYYMMDD')||'-'||LPAD(COALESCE(COUNT(*)+1,1)::text,4,'0') FROM work_orders WHERE tenant_id=$1 AND DATE(created_at)=CURRENT_DATE`, tenantID).Scan(&woNum)

	tx, _ := h.DB.Begin(context.Background())
	defer tx.Rollback(context.Background())

	err := tx.QueryRow(context.Background(), `
		INSERT INTO work_orders (tenant_id, branch_id, wo_number, customer_id, vehicle_id, status, complaint, estimated_cost, estimated_hours, created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id, wo_number, status, created_at
	`, tenantID, input.BranchID, woNum, input.CustomerID, input.VehicleID,
		"menunggu", input.Complaint, input.EstimatedCost, input.EstimatedHours, userID).Scan(
		&input.ID, &input.WONumber, &input.Status, &input.CreatedAt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Insert services
	for _, s := range input.Services {
		tx.Exec(context.Background(), `INSERT INTO wo_services (work_order_id, service_name, description, quantity, unit_price, total_price) VALUES ($1,$2,$3,$4,$5,$6)`,
			input.ID, s.ServiceName, s.Description, s.Quantity, s.UnitPrice, s.TotalPrice)
	}

	// Insert spareparts
	for _, sp := range input.Spareparts {
		tx.Exec(context.Background(), `INSERT INTO wo_spareparts (work_order_id, sparepart_id, sparepart_name, quantity, unit_price, total_price) VALUES ($1,$2,$3,$4,$5,$6)`,
			input.ID, sp.SparepartID, sp.SparepartName, sp.Quantity, sp.UnitPrice, sp.TotalPrice)
	}

	// Insert mechanics
	for _, m := range input.Mechanics {
		tx.Exec(context.Background(), `INSERT INTO wo_mechanics (work_order_id, user_id, commission) VALUES ($1,$2,$3)`,
			input.ID, m.UserID, m.Commission)
	}

	tx.Commit(context.Background())
	return c.Status(201).JSON(input.WorkOrder)
}

func (h *WorkOrderHandler) UpdateStatus(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	var input struct {
		Status           string  `json:"status"`
		Diagnosis        string  `json:"diagnosis"`
		MechanicNotes    string  `json:"mechanic_notes"`
		EstimatedCost    float64 `json:"estimated_cost"`
		CustomerApproved bool    `json:"customer_approved"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	_, err := h.DB.Exec(context.Background(), `
		UPDATE work_orders SET status=$1, diagnosis=$2, mechanic_notes=$3, estimated_cost=$4,
		customer_approved=$5,
		started_at = CASE WHEN $1='dikerjakan' AND started_at IS NULL THEN NOW() ELSE started_at END,
		completed_at = CASE WHEN $1 IN ('selesai','diambil') THEN NOW() ELSE completed_at END,
		updated_at = NOW()
		WHERE id=$6 AND tenant_id=$7
	`, input.Status, input.Diagnosis, input.MechanicNotes, input.EstimatedCost,
		input.CustomerApproved, id, tenantID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Status updated"})
}

func (h *WorkOrderHandler) Delete(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	_, err := h.DB.Exec(context.Background(), `DELETE FROM work_orders WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Deleted"})
}
