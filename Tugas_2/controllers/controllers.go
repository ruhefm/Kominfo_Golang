package controllers

import (
	"encoding/json"
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

// func CreateItems(c *gin.Context) {
// 	var requestItems []models.Item
// 	// Code        string `json:"code" gorm:"type:varchar(10)"`
// 	// Description string `json:"description" gorm:"type:varchar(50)"`
// 	// Quantity    int64  `json:"quantity" gorm:"type:bigint"`
// 	// OrdersID     uint   `json:"order_id"`
// 	if err := c.BindJSON(&requestItems); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var responseItems []struct {
// 		ItemCode    string `json:"itemCode"`
// 		Description string `json:"description"`
// 		Quantity    int64  `json:"quantity"`
// 	}

// 	for _, requestItem := range requestItems {
// 		items := models.Item{
// 			OrdersID:    requestItem.OrdersID,
// 			Code:        requestItem.Code,
// 			Description: requestItem.Description,
// 			Quantity:    requestItem.Quantity,
// 		}

// 		if err := database.CreateItems(&items); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create items"})
// 			return
// 		}

// 		responseItems = append(responseItems, struct {
// 			ItemCode    string `json:"itemCode"`
// 			Description string `json:"description"`
// 			Quantity    int64  `json:"quantity"`
// 		}{
// 			ItemCode:    items.Code,
// 			Description: items.Description,
// 			Quantity:    items.Quantity,
// 		})
// 	}

// 	c.JSON(http.StatusOK, responseItems)
// }

// // @Summary Untuk menambahkan user di orders / to insert user in orders
// // @Description CreateOrder
// // @ID create-order
// // @Accept  json
// // @Produce  json
// // @Param orders body models.Orders true "Orders object that needs to be added to the order"
// // @Success 200 {object} models.Orders
// // @Router /orders [post]

// func CreateOrder(c *gin.Context) {
// 	var request []models.Orders

// 	if err := c.BindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var responses []struct {
// 		ID           uint      `json:"id"`
// 		OrderedAt    time.Time `json:"orderedAt"`
// 		CustomerName string    `json:"customerName"`
// 		Items        []Item    `json:"items"`
// 	}

// 	for _, requestItem := range request {
// 		var OrderedAtTime time.Time
// 		if requestItem.OrderedAt.IsZero() {
// 			OrderedAtTime = time.Now()
// 		} else {
// 			OrderedAtTime = requestItem.OrderedAt
// 		}

// 		order := models.Orders{
// 			CustomerName: requestItem.CustomerName,
// 			OrderedAt:    OrderedAtTime,
// 		}

// 		if err := database.CreateOrder(&order); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
// 			return
// 		}

// 		response := struct {
// 			ID           uint      `json:"id"`
// 			OrderedAt    time.Time `json:"orderedAt"`
// 			CustomerName string    `json:"customerName"`
// 			Items        []Item    `json:"items"`
// 		}{
// 			ID:           order.ID,
// 			OrderedAt:    order.OrderedAt,
// 			CustomerName: order.CustomerName,
// 		}

// 		responses = append(responses, response)
// 	}

// 	c.JSON(http.StatusOK, responses)
// }

func CreateOrders(c *gin.Context) {
	var request []struct {
		OrderedAt    time.Time `json:"orderedAt"`
		CustomerName string    `json:"customerName"`
		Items        []struct {
			ItemCode    string `json:"itemCode"`
			Description string `json:"description"`
			Quantity    int    `json:"quantity"`
		} `json:"items"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var response []struct {
		ID           uint      `json:"id"`
		OrderedAt    time.Time `json:"orderedAt"`
		CustomerName string    `json:"customerName"`
		Items        []Item    `json:"items"`
	}

	for _, req := range request {
		// Simpan order
		newOrder := models.Orders{
			CustomerName: req.CustomerName,
			OrderedAt:    req.OrderedAt,
		}
		if err := database.CreateOrder(&newOrder); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
			return
		}

		// Simpan items
		var items []Item
		for _, item := range req.Items {
			newItem := models.Item{
				OrdersID:    newOrder.ID,
				Code:        item.ItemCode,
				Description: item.Description,
				Quantity:    int64(item.Quantity),
			}
			if err := database.CreateItems(&newItem); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
				return
			}
			items = append(items, Item{
				ID:          newItem.ID,
				Code:        newItem.Code,
				Description: newItem.Description,
				Quantity:    int(newItem.Quantity),
			})
		}

		response = append(response, struct {
			ID           uint      `json:"id"`
			OrderedAt    time.Time `json:"orderedAt"`
			CustomerName string    `json:"customerName"`
			Items        []Item    `json:"items"`
		}{
			ID:           newOrder.ID,
			OrderedAt:    newOrder.OrderedAt,
			CustomerName: newOrder.CustomerName,
			Items:        items,
		})
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

func GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	orderID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := database.GetOrderByID(uint(orderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order"})
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

	response := struct {
		ID           uint      `json:"id"`
		OrderedAt    time.Time `json:"orderedAt"`
		CustomerName string    `json:"customerName"`
		Items        []Item    `json:"items"`
	}{
		ID:           order.ID,
		OrderedAt:    order.OrderedAt,
		CustomerName: order.CustomerName,
		Items:        items,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal response"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonResponse)
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
	var request []struct {
		ID           uint      `json:"id"`
		OrderedAt    time.Time `json:"orderedAt"`
		CustomerName string    `json:"customerName"`
		Items        []struct {
			ID          uint   `json:"id"`
			ItemCode    string `json:"itemCode"`
			Description string `json:"description"`
			Quantity    int    `json:"quantity"`
		} `json:"items"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var responses []struct {
		ID           uint      `json:"id"`
		OrderedAt    time.Time `json:"orderedAt"`
		CustomerName string    `json:"customerName"`
		Items        []Item    `json:"items"`
	}

	for _, req := range request {
		updateID := req.ID

		order := models.Orders{
			ID:           updateID,
			CustomerName: req.CustomerName,
			OrderedAt:    req.OrderedAt,
		}

		if err := database.UpdateOrder(updateID, &order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
			return
		}

		var items []Item
		for _, item := range req.Items {
			newItem := models.Item{
				ID:          item.ID,
				OrdersID:    updateID,
				Code:        item.ItemCode,
				Description: item.Description,
				Quantity:    int64(item.Quantity),
			}
			if err := database.UpdateItem(item.ID, &newItem); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
				return
			}
			items = append(items, Item{
				ID:          newItem.ID,
				Code:        newItem.Code,
				Description: newItem.Description,
				Quantity:    int(newItem.Quantity),
			})
		}

		responses = append(responses, struct {
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

	jsonResponse, err := json.Marshal(responses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal response"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonResponse)
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
