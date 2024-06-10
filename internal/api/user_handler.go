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

// Create a user
func (u *HTTPHandler) RegisterUser(c *gin.Context) {
	var user *models.User
	if err := c.ShouldBind(&user); err != nil {
		util.Response(c, "invalid request", 400, "bad request body", nil)
		return
	}

	//validate user email
	if !util.IsValidEmail(user.Email) {
		util.Response(c, "invalid email", 400, "Bad request body", nil)
		return
	}

	//check if user already exists
	_, err := u.Repository.FindUserByEmail(user.Email)
	if err == nil {
		util.Response(c, "User already exists", 400, "Bad request body", nil)
		return
	}

	//hash password
	hashPass, err := util.HashPassword(user.Password)
	if err != nil {
		util.Response(c, "could not hash password", 500, "internal server error", nil)
		return
	}

	user.Password = hashPass

	//generate account number
	acctNo, err := util.GenerateAccountNumber()
	if err != nil {
		util.Response(c, "could not generate account number", 500, "internal server error", nil)
		return
	}

	user.AccountNo = acctNo

	//set available balance to zero
	user.AvailableBalance = 0.0

	//persist information in the data base
	err = u.Repository.CreateUser(user)
	if err != nil {
		util.Response(c, "user not created", 400, err.Error(), nil)
		return
	}
	util.Response(c, "user created", 200, "success", nil)
}

func (u *HTTPHandler) LoginUser(c *gin.Context) {
	var loginRequest *models.LoginRequest
	if err := c.ShouldBind(&loginRequest); err != nil {
		util.Response(c, "invalid request", 400, "bad request body", nil)
		return
	}

	if loginRequest.Email == "" || loginRequest.Password == "" {
		util.Response(c, "Please enter your email and/or password", 400, "bad request body", nil)
		return
	}

	user, err := u.Repository.FindUserByEmail(loginRequest.Email)
	if err != nil {
		util.Response(c, "user does not exist", 404, "user not found", nil)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		util.Response(c, "invalid email or password", 400, "invalid email or password", nil)
		return
	}

	//Generate token
	accessClaims, refreshClaims := middleware.GenerateClaims(user.Email)

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
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

func (u *HTTPHandler) GetUserByEmail(c *gin.Context) {
	_, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 500, "user not found", nil)
		return
	}

	email := c.Query("email")

	if email == "" {
		util.Response(c, "email is required", 400, "email is required", nil)
		return
	}

	user, err := u.Repository.FindUserByEmail(email)
	if err != nil {
		util.Response(c, "user not fount", 500, "user not found", nil)
		return
	}

	util.Response(c, "user found", 200, user, nil)
}

func (u *HTTPHandler) TransferFunds(c *gin.Context) {
	//declare request body struct
	var transferRequest *models.TransferRequest

	//bind JSON data to struct
	if err := c.ShouldBind(&transferRequest); err != nil {
		util.Response(c, "invalid request", 400, "bad request body", nil)
		return
	}

	//Get user from context
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 500, "user not found", nil)
		return
	}

	//validate the amount
	if transferRequest.Amount <= 0 {
		util.Response(c, "invalid amount", 400, "invalid amount", nil)
		return
	}

	//check if the account number exist
	recipient, err := u.Repository.FindUserByAccountNumber(transferRequest.AccountNumber)
	if err != nil {
		util.Response(c, "account number does not exist", 400, "account number does not exist", nil)
		return
	}

	//check if amount being transferred is less than the user's current balance
	if user.AvailableBalance < transferRequest.Amount {
		util.Response(c, "insufficient funds", 400, "insufficient funds", nil)
		return
	}

	//persist the data into the db
	err = u.Repository.TransferFunds(user, recipient, transferRequest.Amount)
	if err != nil {
		util.Response(c, "transfer failed", 500, "transfer failed", nil)
		return
	}

	util.Response(c, "transfer successful", 200, "transfer successful", nil)
}

// Add money to user account
func (u *HTTPHandler) AddMoney(c *gin.Context) {
	//declare request body struct
	var transferRequest *models.TransferRequest

	//bind JSON data to struct
	if err := c.ShouldBind(&transferRequest); err != nil {
		util.Response(c, "invalid request", 400, "bad request body", nil)
		return
	}

	//Get user from context
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 500, "user not found", nil)
		return
	}

	//validate the amount
	if transferRequest.Amount <= 0 {
		util.Response(c, "invalid amount", 400, "invalid amount", nil)
		return
	}

	//add the amount to the user's account
	user.AvailableBalance += transferRequest.Amount

	//persist the data into the db
	err = u.Repository.UpdateUser(user)
	if err != nil {
		util.Response(c, "add money failed", 500, "add money failed", nil)
		return
	}

	util.Response(c, "add money successful", 200, "add money successful", nil)
}

func (u *HTTPHandler) BalanceCheck(c *gin.Context) {

	// get user from context
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "user not fount", 500, "user not found", nil)
		return
	}
	// checking balance
	util.Response(c, "Balance retrieved successfully", 200, "sucess", nil)
	c.IndentedJSON(200, gin.H{"balance": user.AvailableBalance})

}

// Transaction history

func (u *HTTPHandler) TransactionHistory(c *gin.Context) {

	// get user from context
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "user not fount", 500, "user not found", nil)
		return
	}
	// Add/ create a  function to the repository  to retrive the transacton history
	transaction, err := u.Repository.Transaction(user.AccountNo)
	if err != nil {
		util.Response(c, "count not retrieve trasaction", 500, "not retrieved", nil)
		return
	}
	util.Response(c, "transaction successfully retrieved", 200, "successful", nil)
	c.IndentedJSON(200, gin.H{"Transaction History": transaction})
}
