package database

import (
	"fmt"
	"tugas2/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	username = "postgres"
	password = "test123456"
	dbName   = "db-go-sql"
	port     = 5432
)

var (
	db  *gorm.DB
	err error
)

func StartDB() {
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, username, password, dbName, port)
	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Orders{}, &models.Item{})
	fmt.Println("DB LOG: DB Connected")
}
func GetDB() *gorm.DB {
	return db
}

func CreateOrder(order *models.Orders) error {
	db := GetDB()

	if err := db.Create(order).Error; err != nil {
		return err
	}

	return nil
}

func CreateItems(order *models.Item) error {
	db := GetDB()

	if err := db.Create(order).Error; err != nil {
		return err
	}

	return nil
}

func GetOrders() ([]models.Orders, error) {
	var orders []models.Orders
	db := GetDB()
	if err := db.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func UpdateOrder(id uint, updatedOrder *models.Orders) error {
	db := GetDB()

	var order models.Orders
	if err := db.First(&order, id).Error; err != nil {
		return err
	}

	order.CustomerName = updatedOrder.CustomerName

	if err := db.Save(&order).Error; err != nil {
		return err
	}

	return nil
}

func DeleteOrder(id uint) error {
	db := GetDB()

	var order models.Orders
	if err := db.First(&order, id).Error; err != nil {
		return err
	}

	if err := db.Delete(&order).Error; err != nil {
		return err
	}

	return nil
}
