package seeder

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/user/grafana-analytics-app/internal/models"
	"gorm.io/gorm"
)

type Config struct {
	CustomerCount int
	ProductCount  int
	OrderCount    int
}

func SeedDatabase(db *gorm.DB, cfg Config) error {
	log.Println("--- Seeding COMPREHENSIVE ENTERPRISE DATA ---")

	locations := []string{"Jakarta", "Surabaya", "Bandung", "Medan", "Bali"}
	var stores []models.Store
	for _, loc := range locations {
		tax := 0.11
		if loc == "Bali" { tax = 0.15 } // Higher tax in Bali
		s := models.Store{Name: "Store " + loc, Location: loc, TaxRate: tax}
		db.Create(&s)
		stores = append(stores, s)
	}

	roles := []string{"Sales", "Cashier", "Manager"}
	var staffs []models.Staff
	for _, store := range stores {
		for i := 0; i < 5; i++ {
			st := models.Staff{
				StoreID: store.ID,
				Name:    fmt.Sprintf("Staff %s-%d", store.Location, i),
				Role:    roles[rand.Intn(len(roles))],
			}
			db.Create(&st)
			staffs = append(staffs, st)
		}
	}

	promos := []models.Promotion{
		{Code: "MEGA20", Discount: 20.0, ValidUntil: time.Now().AddDate(0, 1, 0)},
		{Code: "NEWYEAR50", Discount: 50.0, ValidUntil: time.Now().AddDate(0, 1, 0)},
	}
	for i := range promos { db.Create(&promos[i]) }

	var products []models.Product
	cats := []string{"Electronics", "Fashion", "Home", "Sports"}
	for i := 0; i < cfg.ProductCount; i++ {
		p := models.Product{
			Name:     fmt.Sprintf("SKU-%d", i),
			Category: cats[rand.Intn(len(cats))],
			Price:    50 + rand.Float64()*950,
			Cost:     10 + rand.Float64()*40,
			Stock:    100 + rand.Intn(1000),
		}
		db.Create(&p)
		
		// Attributes & Logs
		db.Create(&models.Attribute{ProductID: p.ID, Key: "Condition", Value: "New"})
		db.Create(&models.InventoryLog{ProductID: p.ID, Change: p.Stock, Reason: "Initial Stock", CreatedAt: time.Now()})
		
		products = append(products, p)
	}

	var customers []models.Customer
	segments := []string{"VIP", "Regular", "New"}
	countries := []string{"Indonesia", "USA", "Germany", "Japan"}
	for i := 0; i < cfg.CustomerCount; i++ {
		c := models.Customer{
			Name:          fmt.Sprintf("Cust %d", i),
			Email:         fmt.Sprintf("c%d%d@mail.com", i, time.Now().UnixNano()),
			Segment:       segments[rand.Intn(len(segments))],
			Country:       countries[rand.Intn(len(countries))],
			LoyaltyPoints: rand.Intn(1000),
			CreatedAt:     time.Now().AddDate(0, -rand.Intn(12), -rand.Intn(30)),
		}
		db.Create(&c)
		customers = append(customers, c)
	}

	log.Printf("Commencing Massive Batch Order Insert (%d)...", cfg.OrderCount)
	batchSize := 500
	for i := 0; i < cfg.OrderCount; i += batchSize {
		for j := 0; j < batchSize; j++ {
			if i+j >= cfg.OrderCount { break }

			cust := customers[rand.Intn(len(customers))]
			store := stores[rand.Intn(len(stores))]
			staff := staffs[rand.Intn(len(staffs))]
			
			var promoID *uint
			discountRate := 0.0
			if rand.Float64() < 0.2 { // 20% use promos
				p := promos[rand.Intn(len(promos))]
				promoID = &p.ID
				discountRate = p.Discount
			}

			subtotal := 0.0
			orderedProducts := make([]models.Product, 1+rand.Intn(3))
			for k := range orderedProducts {
				p := products[rand.Intn(len(products))]
				orderedProducts[k] = p
				subtotal += p.Price
			}

			taxAmount := subtotal * store.TaxRate
			finalPrice := (subtotal - (subtotal * discountRate / 100)) + taxAmount

			status := "PAID"
			if rand.Float64() < 0.05 { status = "REFUNDED" }

			// Make 20% of orders happen within the last 24 hours (for "Last 6 hours" default Grafana views)
			var orderDate time.Time
			if rand.Float64() < 0.2 {
				orderDate = time.Now().Add(-time.Duration(rand.Intn(24)) * time.Hour)
			} else {
				orderDate = time.Now().AddDate(0, 0, -1 - rand.Intn(180))
			}

			order := models.Order{
				CustomerID:    cust.ID,
				StoreID:       store.ID,
				StaffReferral: staff.ID,
				PromotionID:   promoID,
				Status:        status,
				SubTotal:      subtotal,
				TaxAmount:     taxAmount,
				FinalPrice:    finalPrice,
				OrderDate:     orderDate,
			}
			db.Create(&order)

			// Items
			for _, p := range orderedProducts {
				db.Create(&models.OrderItem{OrderID: order.ID, ProductID: p.ID, Quantity: 1, UnitPrice: p.Price})
			}

			// Payment
			db.Create(&models.Payment{OrderID: order.ID, Method: "Gateway", Status: "SUCCESS", PaidAt: orderDate})

			// Shipping
			carrier := []string{"FastEx", "SlowPost", "EagleDelivery"}[rand.Intn(3)]
			db.Create(&models.Shipping{OrderID: order.ID, Carrier: carrier, ShippingCost: 15.0, EstimatedDays: 1+rand.Intn(5)})

			// Refund
			if status == "REFUNDED" {
				reason := []string{"Broken Item", "Late Delivery", "Changed Mind"}[rand.Intn(3)]
				db.Create(&models.Refund{OrderID: order.ID, Amount: finalPrice * 0.5, Reason: reason, CreatedAt: orderDate})
			}
			
			// Customer Review
			if rand.Float64() < 0.1 { // 10% leave review
				db.Create(&models.Review{
					ProductID:  orderedProducts[0].ID,
					CustomerID: cust.ID,
					Rating:     1 + rand.Intn(5),
					Comment:    "Verified Purchase",
					CreatedAt:  orderDate.AddDate(0, 0, 5),
				})
			}
		}
		log.Printf("Progress: %d executed...", i)
	}

	log.Println("--- COMPREHENSIVE SEEDING COMPLETED! ---")
	return nil
}
