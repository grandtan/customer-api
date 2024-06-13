package handlers

import (
	"github.com/gin-gonic/gin"
)

type CustomerHandler interface {
	CreateCustomer(c *gin.Context)
	UpdateCustomer(c *gin.Context)
	DeleteCustomer(c *gin.Context)
	GetCustomer(c *gin.Context)
}

type handler struct{}

func NewHandler() CustomerHandler {
	return &handler{}
}
