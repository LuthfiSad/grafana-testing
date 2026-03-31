package models

import (
	"time"
	"gorm.io/gorm"
)

type Customer struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	Email     string         `gorm:"unique" json:"email"`
	Segment   string         `json:"segment"` 
	Country   string         `json:"country"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Orders    []Order        `json:"orders"`
	Reviews   []Review       `json:"reviews"`
	CreatedAt time.Time      `json:"created_at"`
}

type Store struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Staff     []Staff   `gorm:"foreignKey:StoreID" json:"staff"`
	Orders    []Order   `gorm:"foreignKey:StoreID" json:"orders"`
}

type Staff struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StoreID   uint      `json:"store_id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Orders    []Order   `gorm:"foreignKey:StaffReferral" json:"orders"`
}

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CategoryID  uint      `json:"category_id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Cost        float64   `json:"cost"`
	Quantity    int       `json:"quantity"`
	Reviews     []Review  `gorm:"foreignKey:ProductID" json:"reviews"`
}

type Review struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProductID  uint      `json:"product_id"`
	CustomerID uint      `json:"customer_id"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment"`
}

type Promotion struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Code       string    `gorm:"unique" json:"code"`
	Discount   float64   `json:"discount"`
	ValidUntil time.Time `json:"valid_until"`
	Orders     []Order   `gorm:"foreignKey:PromotionID" json:"orders"`
}

type Order struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CustomerID    uint           `json:"customer_id"`
	StoreID       uint           `json:"store_id"`
	StaffReferral uint           `json:"staff_referral"`
	PromotionID   *uint          `json:"promotion_id"`
	Status        string         `json:"status"`
	TotalPrice    float64        `json:"total_price"`
	FinalPrice    float64        `json:"final_price"`
	OrderDate     time.Time      `json:"order_date"`
	OrderItems    []OrderItem    `gorm:"foreignKey:OrderID" json:"items"`
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}
