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
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644) //path file, write read only, create if tidak ada, trunc if ada, granted power privileges.
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

// type Credit struct {
// 	operatorName string `json:"operatorName"`
// 	credits      int    `json:"credits"`
// }

func saveCredits(credits map[string]int, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.Marshal(credits)
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	fmt.Println("Credits saved successfully")
	return nil
}

func loadCredits(filePath string) (map[string]int, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	credits := make(map[string]int)

	if len(file) > 0 {
		err = json.Unmarshal(file, &credits)
		if err != nil {
			return nil, err
		}
	}

	return credits, nil
}
func StartServer() *gin.Engine {
	router := gin.Default()
	router.Static("/public", "./views")

	router.GET("/", func(c *gin.Context) {
		c.File("./views/halaman.html")
	})
	router.POST("/pump_water_wind", func(c *gin.Context) {
		filePath := "./views/status.json"
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
	})

	router.POST("/addCredit", func(c *gin.Context) {
		filePath := "./views/credit.json"
		operatorName := c.PostForm("name")
		if operatorName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for operator name"})
			return
		}
		credits, err := loadCredits(filePath)
		if err != nil {
			fmt.Println("Error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load credits"})
			return
		}
		credit, ok := credits[operatorName]
		if !ok {
			credits[operatorName] = 1
		} else {
			credits[operatorName] = credit + 1
		}

		err = saveCredits(credits, filePath)
		if err != nil {
			fmt.Println("Error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save credits"})
			return
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
