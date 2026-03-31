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
	log.Println("--- Seeding ENTERPRISE RETAIL DATA (MASSIVE) ---")

	// 1. Seed Stores
	locations := []string{"Jakarta", "Surabaya", "Bandung", "Medan", "Bali"}
	var stores []models.Store
	for _, loc := range locations {
		s := models.Store{Name: "Store " + loc, Location: loc}
		db.Create(&s)
		stores = append(stores, s)
	}

	// 2. Seed Staff per Store
	roles := []string{"Sales", "Cashier", "Manager"}
	var allStaff []models.Staff
	for _, store := range stores {
		for i := 0; i < 10; i++ {
			s := models.Staff{
				StoreID: store.ID,
				Name:    fmt.Sprintf("Staff %s-%d", store.Location, i),
				Role:    roles[rand.Intn(len(roles))],
			}
			db.Create(&s)
			allStaff = append(allStaff, s)
		}
	}

	// 3. Seed Promotions
	promos := []models.Promotion{
		{Code: "SALE20", Discount: 20.0, ValidUntil: time.Now().AddDate(0, 1, 0)},
		{Code: "VIP50", Discount: 50.0, ValidUntil: time.Now().AddDate(0, 1, 0)},
		{Code: "FLASH10", Discount: 10.0, ValidUntil: time.Now().AddDate(0, 1, 0)},
	}
	for i := range promos { db.Create(&promos[i]) }

	// 4. Products & Customers (Reuse)
	var products []models.Product
	for i := 0; i < cfg.ProductCount; i++ {
		p := models.Product{Name: fmt.Sprintf("Prop %d", i), Price: 50 + rand.Float64()*950, Cost: 10 + rand.Float64()*40, Quantity: 100 + rand.Intn(1000)}
		db.Create(&p)
		products = append(products, p)
	}

	var customers []models.Customer
	segments := []string{"VIP", "Regular", "New"}
	countries := []string{"Indonesia", "USA", "Germany", "Japan", "Brazil"}
	for i := 0; i < cfg.CustomerCount; i++ {
		c := models.Customer{Name: fmt.Sprintf("Cust %d", i), Email: fmt.Sprintf("c%d%d@mail.com", i, time.Now().UnixNano()), Segment: segments[rand.Intn(len(segments))], Country: countries[rand.Intn(len(countries))]}
		db.Create(&c)
		customers = append(customers, c)
	}

	// 5. SEED ORDERS (10k+)
	log.Printf("Seeding %d Enterprise Orders...", cfg.OrderCount)
	for i := 0; i < cfg.OrderCount; i++ {
		cust := customers[rand.Intn(len(customers))]
		store := stores[rand.Intn(len(stores))]
		staff := allStaff[rand.Intn(len(allStaff))]
		
		var promoID *uint
		discount := 0.0
		if rand.Float64() < 0.3 {
			p := promos[rand.Intn(len(promos))]
			promoID = &p.ID
			discount = p.Discount
		}

		order := models.Order{
			CustomerID:    cust.ID,
			StoreID:       store.ID,
			StaffReferral: staff.ID,
			PromotionID:   promoID,
			Status:        "PAID",
			OrderDate:     time.Now().AddDate(0, 0, -rand.Intn(90)),
		}
		db.Create(&order)

		// Items & Calculation
		totalPr := 0.0
		for j := 0; j < 1+rand.Intn(4); j++ {
			p := products[rand.Intn(len(products))]
			qty := 1 + rand.Intn(3)
			item := models.OrderItem{OrderID: order.ID, ProductID: p.ID, Quantity: qty, UnitPrice: p.Price}
			db.Create(&item)
			totalPr += (p.Price * float64(qty))
		}
		finalPr := totalPr * (1 - (discount / 100))
		db.Model(&order).Updates(map[string]interface{}{"total_price": totalPr, "final_price": finalPr})

		// 6. Seed Reviews (30% customers leave review)
		if rand.Float64() < 0.3 {
			rev := models.Review{
				ProductID:  products[rand.Intn(len(products))].ID,
				CustomerID: cust.ID,
				Rating:     1 + rand.Intn(5),
				Comment:    "Good product!",
			}
			db.Create(&rev)
		}
	}

	log.Println("--- Enterprise Seeding COMPLETED! ---")
	return nil
}
