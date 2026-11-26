package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"devtrack/internal/auth"
	"devtrack/internal/config"
	"devtrack/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func fmtUUID(b [16]byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

type AuthHandler struct {
	queries   *database.Queries
	cfg       *config.Config
	dummyHash string
}

func NewAuthHandler(queries *database.Queries, cfg *config.Config) (*AuthHandler, error) {
	dummy, err := auth.HashPassword("password12ed12d")
	if err != nil {
		return nil, fmt.Errorf("Hash dummy password: %w", err)
	}
	return &AuthHandler{
		queries:   queries,
		cfg:       cfg,
		dummyHash: dummy,
	}, nil
}

const pgUniqueViolation = "23505"

type RegistrationRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	FullName string `json:"full_name" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegistrationRequest
	//1. validate request body
	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Email = strings.ToLower(req.Email)

	//2. Hash password with bcrypt
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}
	//3. Save user to db via sqlc
	user, err := h.queries.CreateUser(c.Request.Context(), database.CreateUserParams{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Role:         "dev",
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}
		log.Printf("create user failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	//4. Convert UUID bytes to string and generate JWT
	userIDStr := fmtUUID(user.ID.Bytes)
	token, err := auth.GenerateToken(userIDStr, user.Email, user.Role, h.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate auth token"})
		return
	}
	//5. Return token and user profile
	c.JSON(http.StatusCreated, gin.H{
		"token": token,
	})
}

type LoginRequest struct { //login payload defination
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,max=72"`
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	//1. validate req body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Email = strings.ToLower(req.Email)

	//2. look up the user by mail via sqlc
	user, err := h.queries.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { //errors.Is-->takes both aurguments as errors; Does this error contain the error i'm looking for
			auth.CheckPassword(req.Password, h.dummyHash)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return

		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	//3. compare password with stored bcrypt hash
	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	//4. Generate JWT
	userIDStr := fmtUUID(user.ID.Bytes)
	token, err := auth.GenerateToken(userIDStr, user.Email, user.Role, h.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate auth token"})
		return
	}

	//5. return token and user profile
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":        userIDStr,
			"email":     user.Email,
			"full_name": user.FullName,
			"role":      user.Role,
		},
	})

}
