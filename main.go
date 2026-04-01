package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	
	"github.com/user/grafana-analytics-app/internal/database"
	"github.com/user/grafana-analytics-app/internal/handler"
	"github.com/user/grafana-analytics-app/internal/repository"
	"github.com/user/grafana-analytics-app/internal/service"
	"github.com/user/grafana-analytics-app/internal/seeder"
)

var (
	ordersProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "retail_orders_processed_total",
		Help: "The total number of processed orders",
	})
	totalRevenueCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "retail_revenue_total",
		Help: "The total revenue grouped by country",
	}, []string{"country"})
)

func main() {
	log.Println("Starting Retail Analytics Engine...")

	// 1. Database Initialization
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

	// 2. Seeding Strategy
	log.Println("Force Refreshing Database for Massive Testing...")
	db.Exec("TRUNCATE TABLE orders, order_items, reviews, customers, products, promotions, staffs, stores, attributes, inventory_logs, refunds, shippings, payments CASCADE")

	ccount, _ := strconv.Atoi(os.Getenv("SEEDER_CUSTOMER_COUNT"))
	pcount, _ := strconv.Atoi(os.Getenv("SEEDER_PRODUCT_COUNT"))
	ocount, _ := strconv.Atoi(os.Getenv("SEEDER_ORDER_COUNT"))

	if ccount == 0 { ccount = 100 }
	if ocount == 0 { ocount = 25000 }

	err = seeder.SeedDatabase(db, seeder.Config{
		CustomerCount: ccount,
		ProductCount:  pcount,
		OrderCount:    ocount,
	})
	if err != nil {
		log.Printf("Seeding error: %v", err)
	}

	// Initialize our Domain Layers (Order)
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService, ordersProcessed, totalRevenueCount)

	// Initialize Customer Domain
	customerRepo := repository.NewCustomerRepository(db)
	customerService := service.NewCustomerService(customerRepo)
	customerHandler := handler.NewCustomerHandler(customerService)

	// Initialize Store Domain
	storeRepo := repository.NewStoreRepository(db)
	storeService := service.NewStoreService(storeRepo)
	storeHandler := handler.NewStoreHandler(storeService)

	// Initialize Product Domain
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	// Initialize Extended Domains
	promoRepo := repository.NewPromotionRepository(db)
	promoSvc := service.NewPromotionService(promoRepo)
	promoHandler := handler.NewPromotionHandler(promoSvc)

	payRepo := repository.NewPaymentRepository(db)
	paySvc := service.NewPaymentService(payRepo)
	payHandler := handler.NewPaymentHandler(paySvc)

	// 3. Setup API & Metrics
	r := gin.Default()

	// Metrics endpoint for Prometheus
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// Add proper API endpoints
	api := r.Group("/api")
	{
		api.POST("/orders", orderHandler.HandleProcessOrder)

		api.POST("/customers", customerHandler.HandleRegister)
		api.POST("/customers/:id/reward", customerHandler.HandleReward)

		api.POST("/stores", storeHandler.HandleCreateStore)

		api.POST("/products", productHandler.HandleAddProduct)

		api.POST("/promotions", promoHandler.HandleCreatePromo)
		api.POST("/payments", payHandler.HandleProcessPayment)
	}

	// Start background routine to update "Business Dashboard" metrics
	go updateBusinessMetrics(orderService)

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	log.Printf("Application serving on port %s", port)
	r.Run(":" + port)
}

func updateBusinessMetrics(svc service.OrderService) {
	for {
		total, err := svc.CalculateTotalRevenue()
		if err == nil {
			log.Printf("[METRICS] Total Revenue: $%.2f", total)
			prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "retail_current_total_revenue",
				Help: "Current total revenue across all regions",
			}).Set(total)
		} else {
			log.Printf("[METRICS] Error calculating revenue: %v", err)
		}
		// Calculate every 30 seconds for simplicity in demo
		time.Sleep(30 * time.Second)
	}
}
