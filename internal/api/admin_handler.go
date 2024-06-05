package api

import (
	"net/http"
	"os"
	"payment-system-one/internal/middleware"
	"payment-system-one/internal/models"
	"payment-system-one/internal/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (u *HTTPHandler) RegisterAdmin(c *gin.Context) {
	var admin *models.Admin
	if err := c.ShouldBind(&admin); err != nil {
		util.Response(c, "invalid request", 400, "bad request body", nil)
		return
	}

	//validate admin email
	if !util.IsValidEmail(admin.Email) {
		util.Response(c, "invalid email", 400, "Bad request body", nil)
		return
	}

	//check if admin already exists
	_, err := u.Repository.FindAdminByEmail(admin.Email)
	if err == nil {
		util.Response(c, "admin already exists", 400, "Bad request body", nil)
		return
	}

	hashPass, err := util.HashPassword(admin.Password)
	if err != nil {
		util.Response(c, "could not hash password", 500, "internal server error", nil)
		return
	}

	admin.Password = hashPass

	//persist information in the data base
	err = u.Repository.CreateAdmin(admin)
	if err != nil {
		util.Response(c, "admin not created", 400, err.Error(), nil)
		return
	}
	util.Response(c, "admin created", 200, "success", nil)
}

func (u *HTTPHandler) LoginAdmin(c *gin.Context) {
	var loginRequest *models.LoginRequest
	if err := c.ShouldBind(&loginRequest); err != nil {
		util.Response(c, "invalid request", 400, "bad request body", nil)
		return
	}

	if loginRequest.Email == "" || loginRequest.Password == "" {
		util.Response(c, "Please enter your email and/or password", 400, "bad request body", nil)
		return
	}

	admin, err := u.Repository.FindAdminByEmail(loginRequest.Email)
	if err != nil {
		util.Response(c, "admin does not exist", 404, "admin not found", nil)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(loginRequest.Password)); err != nil {
		util.Response(c, "invalid email or password", 400, "invalid email or password", nil)
		return
	}

	//Generate token
	accessClaims, refreshClaims := middleware.GenerateClaims(admin.Email)

	secret := os.Getenv("JWT_SECRET")

	accessToken, err := middleware.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
	if err != nil {
		util.Response(c, "error generating access token", 500, "error generating access token", nil)
		return
	}
	refreshToken, err := middleware.GenerateToken(jwt.SigningMethodHS256, refreshClaims, &secret)
	if err != nil {
		util.Response(c, "error generating refresh token", 500, "error generating refresh token", nil)
		return
	}
	c.Header("access_token", *accessToken)
	c.Header("refresh_token", *refreshToken)

	util.Response(c, "login successful", http.StatusOK, gin.H{
		"admin":         admin,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}
