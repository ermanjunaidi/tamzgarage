package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"bengkelpro-backend/config"
	"bengkelpro-backend/handlers"
	"bengkelpro-backend/middleware"
)

func Setup(app *fiber.App, db *pgxpool.Pool) {
	// Set JWT secret
	cfg := config.Load()
	middleware.SetJWTSecret(cfg.JWTSecret)

	// Handlers
	authH := handlers.NewAuthHandler(db)
	dashH := handlers.NewDashboardHandler(db)
	custH := handlers.NewCustomerHandler(db)
	vehH := handlers.NewVehicleHandler(db)
	woH := handlers.NewWorkOrderHandler(db)
	invH := handlers.NewInventoryHandler(db)
	empH := handlers.NewEmployeeHandler(db)
	posH := handlers.NewPOSHandler(db)
	supH := handlers.NewSupplierHandler(db)
	rptH := handlers.NewReportHandler(db)

	api := app.Group("/api")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/login", authH.Login)

	// Protected routes
	protected := api.Group("", middleware.AuthRequired())

	// Profile
	protected.Get("/profile", authH.GetProfile)

	// Users (super_admin only)
	users := protected.Group("/users", middleware.RequireRole("super_admin", "admin_cabang"))
	users.Get("/", authH.ListUsers)
	users.Post("/", authH.CreateUser)

	// Dashboard
	protected.Get("/dashboard/stats", dashH.GetStats)
	protected.Get("/dashboard/revenue-chart", dashH.GetRevenueChart)

	// Customers
	customers := protected.Group("/customers")
	customers.Get("/", custH.List)
	customers.Get("/:id", custH.Get)
	customers.Post("/", custH.Create)
	customers.Put("/:id", custH.Update)
	customers.Delete("/:id", middleware.RequireRole("super_admin", "admin_cabang"), custH.Delete)

	// Vehicles
	vehicles := protected.Group("/vehicles")
	vehicles.Get("/", vehH.List)
	vehicles.Get("/:id", vehH.Get)
	vehicles.Post("/", vehH.Create)
	vehicles.Put("/:id", vehH.Update)
	vehicles.Delete("/:id", middleware.RequireRole("super_admin", "admin_cabang"), vehH.Delete)

	// Work Orders
	workorders := protected.Group("/work-orders")
	workorders.Get("/", woH.List)
	workorders.Get("/:id", woH.Get)
	workorders.Post("/", woH.Create)
	workorders.Patch("/:id/status", woH.UpdateStatus)
	workorders.Delete("/:id", middleware.RequireRole("super_admin", "admin_cabang"), woH.Delete)

	// Inventory / Spare Parts
	spareparts := protected.Group("/spareparts")
	spareparts.Get("/", invH.ListSpareparts)
	spareparts.Get("/:id", invH.GetSparepart)
	spareparts.Post("/", middleware.RequireRole("super_admin", "admin_cabang", "gudang"), invH.CreateSparepart)
	spareparts.Put("/:id", middleware.RequireRole("super_admin", "admin_cabang", "gudang"), invH.UpdateSparepart)
	spareparts.Delete("/:id", middleware.RequireRole("super_admin", "admin_cabang"), invH.DeleteSparepart)
	spareparts.Post("/stock-in", middleware.RequireRole("super_admin", "admin_cabang", "gudang"), invH.AddStock)

	// Stock Mutations
	protected.Get("/stock-mutations", invH.ListStockMutations)

	// Employees / Mechanics
	employees := protected.Group("/employees")
	employees.Get("/mechanics", empH.ListMechanics)
	employees.Get("/:id", empH.GetEmployee)

	// POS / Invoices
	invoices := protected.Group("/invoices")
	invoices.Get("/", posH.ListInvoices)
	invoices.Get("/:id", posH.GetInvoice)
	invoices.Post("/", middleware.RequireRole("super_admin", "admin_cabang", "kasir", "service_advisor"), posH.CreateInvoice)

	// Payments
	protected.Post("/payments", middleware.RequireRole("super_admin", "admin_cabang", "kasir"), posH.ProcessPayment)

	// Suppliers
	suppliers := protected.Group("/suppliers")
	suppliers.Get("/", supH.List)
	suppliers.Get("/:id", supH.Get)
	suppliers.Post("/", middleware.RequireRole("super_admin", "admin_cabang", "gudang"), supH.Create)
	suppliers.Put("/:id", middleware.RequireRole("super_admin", "admin_cabang", "gudang"), supH.Update)
	suppliers.Delete("/:id", middleware.RequireRole("super_admin", "admin_cabang"), supH.Delete)

	// Purchase Orders
	posGroup := protected.Group("/purchase-orders")
	posGroup.Get("/", supH.ListPOs)
	posGroup.Post("/", middleware.RequireRole("super_admin", "admin_cabang", "gudang"), supH.CreatePO)
	posGroup.Post("/:id/receive", middleware.RequireRole("super_admin", "admin_cabang", "gudang"), supH.ReceivePO)

	// Reports
	reports := protected.Group("/reports")
	reports.Get("/revenue", rptH.RevenueReport)
	reports.Get("/stock", rptH.StockReport)
	reports.Get("/work-orders", rptH.WorkOrderReport)
	reports.Get("/top-customers", rptH.TopCustomers)
}
