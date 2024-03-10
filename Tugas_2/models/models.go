package models

import "time"

type Orders struct {
	ID           uint      `json:"id" gorm:"primary_key;type:bigint"`
	CustomerName string    `json:"customer_name" gorm:"type:varchar(50)"`
	OrderedAt    time.Time `json:"ordered_at" gorm:"type:timestamp"`
	Items        []Items   `json:"items" gorm:"foreignKey:OrdersID"`
}

type Items struct {
	ID          uint   `json:"id" gorm:"primary_key;type:bigint"`
	Code        string `json:"code" gorm:"type:varchar(10)"`
	Description string `json:"description" gorm:"type:varchar(50)"`
	Quantity    int64  `json:"quantity" gorm:"type:bigint"`
	OrdersID    uint   `json:"orders_id" gorm:"type:bigint"`
}
