package api

import (
	"github.com/gin-gonic/gin"

	"payment-system-one/internal/util"
)

// Readiness is to check if server is up
func (u *HTTPHandler) Readiness(c *gin.Context) {
	data := "server is up and running"

	// healthcheck
	util.Response(c, "Ready to go", 200, data, nil)
}

// login

// call a protected route

//query parameter
//path parameter

//100 ---- informtional
//200 ---- success 200, 201, 202
//300 ---- redirect
//400 ---- client error
//500 ----- server error

//syntax error
