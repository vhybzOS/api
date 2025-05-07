package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vhybZApp/api/config"
	"github.com/vhybZApp/api/database"
	"github.com/vhybZApp/api/models"
)

type Claims struct {
	Username string `json:"username"`
	Type     string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

func generateToken(username, tokenType string, expiresIn time.Duration) (string, error) {
	claims := Claims{
		Username: username,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Authorization header is required"))
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid token"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid token claims"))
			c.Abort()
			return
		}

		if claims.Type != TokenTypeAccess {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid token type"))
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}

// @Summary Register a new user
// @Description Create a new user account with username, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User registration data"
// @Success 201 {object} models.RegisterResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/register [post]
func register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	var user database.DBUser
	user.Username = req.Username
	user.Email = req.Email
	if err := user.HashPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Error hashing password"))
		return
	}

	if err := database.GetDB().Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Error creating user"))
		return
	}

	c.JSON(http.StatusCreated, models.NewRegisterResponse())
}

// @Summary Login user
// @Description Authenticate user with email and password, return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login [post]
func login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	var user database.DBUser
	if err := database.GetDB().Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid credentials"))
		return
	}

	if err := user.CheckPassword(req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid credentials"))
		return
	}

	// Generate access token (15 minutes)
	accessToken, err := generateToken(user.Username, TokenTypeAccess, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Error generating access token"))
		return
	}

	// Generate refresh token (7 days)
	refreshToken, err := generateToken(user.Username, TokenTypeRefresh, 7*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Error generating refresh token"))
		return
	}

	c.JSON(http.StatusOK, models.NewTokenResponse(accessToken, refreshToken, 15*60))
}

// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body models.RefreshRequest true "Refresh token"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/refresh [post]
func refresh(c *gin.Context) {
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	token, err := jwt.ParseWithClaims(req.RefreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid refresh token"))
		return
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || claims.Type != TokenTypeRefresh {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid refresh token"))
		return
	}

	// Generate new access token
	accessToken, err := generateToken(claims.Username, TokenTypeAccess, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Error generating access token"))
		return
	}

	// Generate new refresh token
	refreshToken, err := generateToken(claims.Username, TokenTypeRefresh, 7*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Error generating refresh token"))
		return
	}

	c.JSON(http.StatusOK, models.NewTokenResponse(accessToken, refreshToken, 15*60))
}

// @Summary Get user profile
// @Description Get the authenticated user's profile information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/profile [get]
func getProfile(c *gin.Context) {
	username := c.GetString("username")
	var user database.DBUser
	if err := database.GetDB().Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse("User not found"))
		return
	}

	c.JSON(http.StatusOK, models.NewProfileResponse(user.Username, user.Email))
}
