package controllers

import (
	"net/http"
	"strconv"
	"time"
	"tugas2/database"
	"tugas2/models"

	"github.com/gin-gonic/gin"
)

func CreateItems(c *gin.Context) {
	var requestItems models.Item
	// Code        string `json:"code" gorm:"type:varchar(10)"`
	// Description string `json:"description" gorm:"type:varchar(50)"`
	// Quantity    int64  `json:"quantity" gorm:"type:bigint"`
	// OrdersID     uint   `json:"order_id"`
	if err := c.BindJSON(&requestItems); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items := models.Item{
		OrdersID:    requestItems.OrdersID,
		Code:        requestItems.Code,
		Description: requestItems.Description,
		Quantity:    requestItems.Quantity,
	}

	if err := database.CreateItems(&items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create items"})
		return
	}
	c.JSON(http.StatusCreated, items)
}

func CreateOrder(c *gin.Context) {
	var request models.Orders

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var OrderedAtTime time.Time
	if request.OrderedAt.IsZero() {
		OrderedAtTime = time.Now()
	} else {
		OrderedAtTime = request.OrderedAt
	}

	order := models.Orders{
		CustomerName: request.CustomerName,
		OrderedAt:    OrderedAtTime,
	}

	if err := database.CreateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func GetOrders(c *gin.Context) {
	orders, err := database.GetOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var request models.Orders
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateID, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}
	order := models.Orders{
		ID: uint(updateID),
	}

	existingOrder, err := database.GetOrderByID(uint(updateID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order"})
		return
	}
	if request.CustomerName == "" {
		order.CustomerName = existingOrder.CustomerName
	}
	if request.CustomerName != "" {
		order.CustomerName = request.CustomerName
	}
	if !request.OrderedAt.IsZero() {
		order.OrderedAt = request.OrderedAt
	}
	if request.OrderedAt.IsZero() {
		order.OrderedAt = existingOrder.OrderedAt
	}

	if err := database.UpdateOrder(uint(updateID), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	deleteID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	if err := database.DeleteOrder(uint(deleteID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.String(http.StatusOK, "Success delete")
}
