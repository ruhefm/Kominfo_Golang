package routers

import (
	"tugas2/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.POST("/orders", controllers.CreateOrder)
	router.GET("/orders", controllers.GetOrders)
	router.PATCH("/orders/:id", controllers.UpdateOrder)
	router.DELETE("/orders/:id", controllers.DeleteOrder)
	router.POST("/items", controllers.CreateItems)
	router.GET("/items", controllers.GetOrders)

	return router
}
