package models

import (
	"time"
	"gorm.io/gorm"
)

type Customer struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `json:"name"`
	Email         string         `gorm:"unique" json:"email"`
	Segment       string         `json:"segment"` 
	Country       string         `json:"country"`
	LoyaltyPoints int            `json:"loyalty_points"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Orders        []Order        `json:"orders"`
	Reviews       []Review       `json:"reviews"`
	CreatedAt     time.Time      `json:"created_at"`
}

type Store struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	TaxRate   float64   `json:"tax_rate"`
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
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `json:"name"`
	Category      string         `json:"category"`
	Price         float64        `json:"price"`
	Cost          float64        `json:"cost"`
	Stock         int            `json:"stock"`
	Attributes    []Attribute    `gorm:"foreignKey:ProductID" json:"attributes"`
	InventoryLogs []InventoryLog `gorm:"foreignKey:ProductID" json:"logs"`
	Reviews       []Review       `gorm:"foreignKey:ProductID" json:"reviews"`
}

type Attribute struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ProductID uint   `json:"product_id"`
	Key       string `json:"key"`   // e.g. "Color", "Material"
	Value     string `json:"value"` // e.g. "Red", "Cotton"
}

type InventoryLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `json:"product_id"`
	Change    int       `json:"change"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

type Review struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProductID  uint      `json:"product_id"`
	CustomerID uint      `json:"customer_id"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
}

type Promotion struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Code       string    `gorm:"unique" json:"code"`
	Discount   float64   `json:"discount"` // Percentage
	ValidUntil time.Time `json:"valid_until"`
	Orders     []Order   `gorm:"foreignKey:PromotionID" json:"orders"`
}

type Order struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CustomerID    uint           `json:"customer_id"`
	StoreID       uint           `json:"store_id"`
	StaffReferral uint           `json:"staff_referral"`
	PromotionID   *uint          `json:"promotion_id"`
	Status        string         `json:"status"` // PAID, REFUNDED, SHIPPED
	SubTotal      float64        `json:"sub_total"`
	TaxAmount     float64        `json:"tax_amount"`
	FinalPrice    float64        `json:"final_price"`
	OrderDate     time.Time      `gorm:"index" json:"order_date"`
	Payment       Payment        `gorm:"foreignKey:OrderID" json:"payment"`
	Shipping      Shipping       `gorm:"foreignKey:OrderID" json:"shipping"`
	Refund        *Refund        `gorm:"foreignKey:OrderID" json:"refund"`
	OrderItems    []OrderItem    `gorm:"foreignKey:OrderID" json:"items"`
}

type Refund struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `json:"order_id"`
	Amount    float64   `json:"amount"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

type Payment struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OrderID       uint      `json:"order_id"`
	Method        string    `json:"method"` // Gateway, Cash, Crypto
	Status        string    `json:"status"`
	PaidAt        time.Time `json:"paid_at"`
}

type Shipping struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	OrderID        uint       `json:"order_id"`
	Carrier        string     `json:"carrier"`
	ShippingCost   float64    `json:"shipping_cost"`
	EstimatedDays  int        `json:"estimated_days"`
	ShippedAt      *time.Time `json:"shipped_at"`
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}
