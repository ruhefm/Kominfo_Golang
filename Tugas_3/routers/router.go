package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func jsonSimpan(status Status, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777) //path file, write read only, create if tidak ada, trunc if ada, granted power privileges.
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(map[string]Status{"status": status})
	if err != nil {
		return err
	}

	fmt.Println("Updated status:", status)
	return nil
}

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

func StartServer() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.File("./views/halaman.html")
	})
	router.GET("/status.json", func(c *gin.Context) {
		c.File("status.json")
	})
	router.GET("/danger.m4a", func(c *gin.Context) {
		c.File("./views/danger.m4a")
	})
	router.GET("/succed.mp3", func(c *gin.Context) {
		c.File("./views/succed.mp3")
	})
	router.POST("/pump_water_wind", func(c *gin.Context) {
		filePath := "status.json"
		waterIn, err := strconv.Atoi(c.PostForm("water"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for water"})
			return
		}

		windIn, err := strconv.Atoi(c.PostForm("wind"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for wind"})
			return
		}
		status := Status{
			Water: waterIn,
			Wind:  windIn,
		}
		jsonSimpan(status, filePath)
		if err != nil {
			fmt.Println("Error: ", err)
		}
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
