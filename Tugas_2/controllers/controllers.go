package controllers

import (
	"net/http"
	"strconv"
	"time"
	"tugas2/database"
	"tugas2/models"

	"github.com/gin-gonic/gin"
)

type Item struct {
	ID          uint   `json:"id"`
	Code        string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

type OrderedMap []struct {
	ID           uint      `json:"id"`
	OrderedAt    time.Time `json:"orderedAt"`
	CustomerName string    `json:"customerName"`
	Items        []Item    `json:"items"`
}

// @Summary Untuk menambahkan items ke user di orders / to insert items to user in orders
// @Description CreateItems
// @ID create-items
// @Accept  json
// @Produce  json
// @Param item body models.Item true "Item object that needs to be added to the order"
// @Success 200 {object} models.Item
// @Router /items [post]

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
	// Bisa tambahkan marshal
	response := struct {
		ItemCode    string `json:"itemCode"`
		Description string `json:"description"`
		Quantity    int64  `json:"quantity"`
	}{
		ItemCode:    items.Code,
		Description: items.Description,
		Quantity:    items.Quantity,
	}

	c.JSON(http.StatusOK, response)

}

// @Summary Untuk menambahkan user di orders / to insert user in orders
// @Description CreateOrder
// @ID create-order
// @Accept  json
// @Produce  json
// @Param orders body models.Orders true "Orders object that needs to be added to the order"
// @Success 200 {object} models.Orders
// @Router /orders [post]

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

	response := struct {
		ID           uint      `json:"id"`
		OrderedAt    time.Time `json:"orderedAt"`
		CustomerName string    `json:"customerName"`
		Items        []Item    `json:"items"`
	}{
		ID:           order.ID,
		OrderedAt:    order.OrderedAt,
		CustomerName: order.CustomerName,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Retrieve order data
// @Description Retrieve order data
// @ID retrieve-order-data
// @Accept json
// @Produce json
// @Success 200 {object} models.Orders Preload {object} models.Items
// @Router /orders [get]

func GetOrders(c *gin.Context) {
	orders, err := database.GetOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	var response OrderedMap
	for _, order := range orders {
		var items []Item
		for _, item := range order.Items {
			items = append(items, Item{
				ID:          item.ID,
				Code:        item.Code,
				Description: item.Description,
				Quantity:    int(item.Quantity),
			})
		}
		response = append(response, struct {
			ID           uint      `json:"id"`
			OrderedAt    time.Time `json:"orderedAt"`
			CustomerName string    `json:"customerName"`
			Items        []Item    `json:"items"`
		}{
			ID:           order.ID,
			OrderedAt:    order.OrderedAt,
			CustomerName: order.CustomerName,
			Items:        items,
		})
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Update order data
// @Description Update order data
// @ID update-order-data
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param order body models.Orders true "Order Object"
// @Success 200 {object} models.Orders
// @Router /orders/{id} [patch]

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

	var items []Item
	for _, item := range order.Items {
		items = append(items, Item{
			ID:          item.ID,
			Code:        item.Code,
			Description: item.Description,
			Quantity:    int(item.Quantity),
		})
	}
	response := OrderedMap{
		{
			ID:           order.ID,
			OrderedAt:    order.OrderedAt,
			CustomerName: order.CustomerName,
			Items:        items,
		},
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Delete order data
// @Description Delete order data
// @ID Delete-order-data
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 Success delete
// @Router /orders/{id} [delete]

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
