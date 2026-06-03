package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"bengkelpro-backend/middleware"
	"bengkelpro-backend/models"
)

type VehicleHandler struct {
	DB *pgxpool.Pool
}

func NewVehicleHandler(db *pgxpool.Pool) *VehicleHandler {
	return &VehicleHandler{DB: db}
}

func (h *VehicleHandler) List(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	customerID := c.Query("customer_id")
	search := c.Query("search", "")

	query := `
		SELECT v.id, v.tenant_id, v.customer_id, v.plate_number, v.brand, v.model,
		       v.variant, v.year, v.color, v.vin, v.engine_number, v.last_km,
		       v.next_service_km, v.notes, v.created_at, v.updated_at,
		       c.full_name as customer_name
		FROM vehicles v
		LEFT JOIN customers c ON v.customer_id = c.id
		WHERE v.tenant_id = $1
	`
	args := []interface{}{tenantID}
	argN := 2

	if customerID != "" {
		query += ` AND v.customer_id = $` + itoa(argN)
		args = append(args, customerID)
		argN++
	}
	if search != "" {
		query += ` AND (v.plate_number ILIKE $` + itoa(argN) + ` OR v.brand ILIKE $` + itoa(argN) + ` OR c.full_name ILIKE $` + itoa(argN) + `)`
		args = append(args, "%"+search+"%")
		argN++
	}

	query += ` ORDER BY v.created_at DESC`

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	vehicles := []models.Vehicle{}
	for rows.Next() {
		var v models.Vehicle
		if err := rows.Scan(&v.ID, &v.TenantID, &v.CustomerID, &v.PlateNumber, &v.Brand,
			&v.Model, &v.Variant, &v.Year, &v.Color, &v.VIN, &v.EngineNumber,
			&v.LastKM, &v.NextServiceKM, &v.Notes, &v.CreatedAt, &v.UpdatedAt,
			&v.CustomerName); err != nil {
			continue
		}
		vehicles = append(vehicles, v)
	}

	return c.JSON(vehicles)
}

func (h *VehicleHandler) Get(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	var v models.Vehicle
	err := h.DB.QueryRow(context.Background(), `
		SELECT v.id, v.tenant_id, v.customer_id, v.plate_number, v.brand, v.model,
		       v.variant, v.year, v.color, v.vin, v.engine_number, v.last_km,
		       v.next_service_km, v.notes, v.created_at, v.updated_at,
		       c.full_name as customer_name
		FROM vehicles v LEFT JOIN customers c ON v.customer_id = c.id
		WHERE v.id = $1 AND v.tenant_id = $2
	`, id, tenantID).Scan(&v.ID, &v.TenantID, &v.CustomerID, &v.PlateNumber, &v.Brand,
		&v.Model, &v.Variant, &v.Year, &v.Color, &v.VIN, &v.EngineNumber,
		&v.LastKM, &v.NextServiceKM, &v.Notes, &v.CreatedAt, &v.UpdatedAt, &v.CustomerName)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Vehicle not found"})
	}

	return c.JSON(v)
}

func (h *VehicleHandler) Create(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var v models.Vehicle
	if err := c.BodyParser(&v); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := h.DB.QueryRow(context.Background(), `
		INSERT INTO vehicles (tenant_id, customer_id, plate_number, brand, model, variant, year, color, vin, engine_number, last_km, next_service_km, notes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id, tenant_id, plate_number, brand, model, variant, year, color, vin, engine_number, last_km, next_service_km, notes, created_at
	`, tenantID, v.CustomerID, v.PlateNumber, v.Brand, v.Model, v.Variant,
		v.Year, v.Color, v.VIN, v.EngineNumber, v.LastKM, v.NextServiceKM, v.Notes).Scan(
		&v.ID, &v.TenantID, &v.PlateNumber, &v.Brand, &v.Model, &v.Variant,
		&v.Year, &v.Color, &v.VIN, &v.EngineNumber, &v.LastKM, &v.NextServiceKM, &v.Notes, &v.CreatedAt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(v)
}

func (h *VehicleHandler) Update(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	var v models.Vehicle
	if err := c.BodyParser(&v); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	_, err := h.DB.Exec(context.Background(), `
		UPDATE vehicles SET customer_id=$1, plate_number=$2, brand=$3, model=$4, variant=$5,
		year=$6, color=$7, vin=$8, engine_number=$9, last_km=$10, next_service_km=$11, notes=$12, updated_at=NOW()
		WHERE id=$13 AND tenant_id=$14
	`, v.CustomerID, v.PlateNumber, v.Brand, v.Model, v.Variant, v.Year,
		v.Color, v.VIN, v.EngineNumber, v.LastKM, v.NextServiceKM, v.Notes, id, tenantID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Updated"})
}

func (h *VehicleHandler) Delete(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id := c.Params("id")

	_, err := h.DB.Exec(context.Background(), `DELETE FROM vehicles WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Deleted"})
}
