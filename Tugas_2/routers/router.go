package routers

import (
	"tugas2/controllers"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.POST("/orders", controllers.CreateOrders)
	router.GET("/orders", controllers.GetOrders)
	router.GET("/orders/:id", controllers.GetOrderByID)
	router.PATCH("/orders", controllers.UpdateOrder)
	router.DELETE("/orders/:id", controllers.DeleteOrder)
	//Hanya digunakan untuk generate secara terpisah, sementara yang atas dapat dilakukan batch keduanya.
	router.POST("/items", controllers.CreateItems)
	router.POST("/users", controllers.CreateUsers)
	router.GET("/items", controllers.GetOrders)
	//Hanya digunakan untuk generate secara terpisah, sementara yang atas dapat dilakukan batch keduanya.
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/docs/swagger.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})
	router.Use(Cors())
	return router
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
