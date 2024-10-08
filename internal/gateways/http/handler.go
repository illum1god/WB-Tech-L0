package http

import (
	"WB-L0/internal/connectors"
	"WB-L0/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Order interface {
	GetOrders(c *gin.Context)
	GetOrderByUID(c *gin.Context)
	SaveOrder(c *gin.Context)
}

type Handler struct {
	Order
}

func NewHandler(service usecase.Service) *Handler {
	return &Handler{Order: newOrder(service.Order)}
}

type order struct {
	service  usecase.Order
	validate *validator.Validate
}

func newOrder(service usecase.Order) Order {
	return &order{
		service:  service,
		validate: validator.New(),
	}
}

func (o *order) GetOrders(c *gin.Context) {
	orders, err := o.service.GetOrders(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve orders",
		})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No orders found",
		})
		return
	}

	responseOrders := connectors.DomainArrToResponse(orders)
	c.JSON(http.StatusOK, responseOrders)
}

func (o *order) GetOrderByUID(c *gin.Context) {
	orderUID := c.Param("uid")
	if orderUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Order UID is required",
		})
		return
	}

	orderRes, err := o.service.GetOrderByUID(c, orderUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Failed to get order " + orderUID,
		})
		return
	}

	responseOrder := connectors.DomainToResponse(orderRes)
	c.JSON(http.StatusOK, responseOrder)
}

func (o *order) SaveOrder(c *gin.Context) {
	var input usecase.OrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := o.service.SaveOrder(c, input.Order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order saved successfully"})
}
