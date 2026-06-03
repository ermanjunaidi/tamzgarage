package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"bengkelpro-backend/middleware"
	"bengkelpro-backend/models"
)

type AuthHandler struct {
	DB *pgxpool.Pool
}

func NewAuthHandler(db *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	var user models.User
	var tenantID, branchID uuid.UUID
	var branchName *string

	err := h.DB.QueryRow(context.Background(), `
		SELECT u.id, u.tenant_id, u.branch_id, u.username, u.email, u.password_hash,
			   u.full_name, u.role, u.phone, u.is_active,
			   b.name as branch_name
		FROM users u
		LEFT JOIN branches b ON u.branch_id = b.id
		WHERE u.username = $1 AND u.is_active = true
	`, req.Username).Scan(
		&user.ID, &tenantID, &branchID, &user.Username, &user.Email,
		&user.PasswordHash, &user.FullName, &user.Role, &user.Phone, &user.IsActive,
		&branchName,
	)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := middleware.GenerateToken(user.ID, tenantID, branchID, user.Role, user.Username)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	user.TenantID = &tenantID
	user.BranchID = &branchID

	var branch *models.Branch
	if branchName != nil {
		branch = &models.Branch{
			ID:   branchID,
			Name: *branchName,
		}
	}

	return c.JSON(models.LoginResponse{
		Token:  token,
		User:   user,
		Branch: branch,
	})
}

func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var user models.User
	err := h.DB.QueryRow(context.Background(), `
		SELECT id, tenant_id, branch_id, username, email, full_name, role, phone, is_active
		FROM users WHERE id = $1
	`, userID).Scan(&user.ID, &user.TenantID, &user.BranchID, &user.Username, &user.Email,
		&user.FullName, &user.Role, &user.Phone, &user.IsActive)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

func (h *AuthHandler) ListUsers(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	rows, err := h.DB.Query(context.Background(), `
		SELECT id, tenant_id, branch_id, username, email, full_name, role, phone, is_active, created_at
		FROM users WHERE tenant_id = $1 ORDER BY created_at DESC
	`, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.TenantID, &u.BranchID, &u.Username, &u.Email,
			&u.FullName, &u.Role, &u.Phone, &u.IsActive, &u.CreatedAt); err != nil {
			continue
		}
		users = append(users, u)
	}

	return c.JSON(users)
}

func (h *AuthHandler) CreateUser(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
		Role     string `json:"role"`
		Phone    string `json:"phone"`
		BranchID string `json:"branch_id"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	var user models.User
	err = h.DB.QueryRow(context.Background(), `
		INSERT INTO users (tenant_id, branch_id, username, email, password_hash, full_name, role, phone)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, tenant_id, username, email, full_name, role, phone, is_active, created_at
	`, tenantID, input.BranchID, input.Username, input.Email, string(hash),
		input.FullName, input.Role, input.Phone).Scan(
		&user.ID, &user.TenantID, &user.Username, &user.Email,
		&user.FullName, &user.Role, &user.Phone, &user.IsActive, &user.CreatedAt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(user)
}
