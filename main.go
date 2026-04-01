package main

import (
    "log"
    "net/http"
    "os"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/user/grafana-analytics-app/internal/database"
    "github.com/user/grafana-analytics-app/internal/seeder"
    "github.com/user/grafana-analytics-app/internal/models"
    "gorm.io/gorm"
)

var (
    ordersProcessed = promauto.NewCounter(prometheus.CounterOpts{
        Name: "retail_orders_processed_total",
        Help: "The total number of processed orders",
    })
    totalRevenue = promauto.NewCounterVec(prometheus.CounterOpts{
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

    // 2. Seeding Strategy (Always Refresh for Testing)
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

    // 3. Setup API & Metrics
    r := gin.Default()

    // Metrics endpoint for Prometheus
    r.GET("/metrics", gin.WrapH(promhttp.Handler()))

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "UP"})
    })

    // Simulated Order Endpoint (Adds to metrics)
    r.POST("/api/orders", func(c *gin.Context) {
        ordersProcessed.Inc()
        // Dummy data for country
        countries := []string{"Indonesia", "USA", "Germany", "Japan"}
        country := countries[models.Customer{}.ID % uint(len(countries))]
        totalRevenue.WithLabelValues(country).Add(150.0)
        c.JSON(http.StatusCreated, gin.H{"message": "Order processed"})
    })

    // Start background routine to update "Business Dashboard" metrics
    go updateBusinessMetrics(db)

    port := os.Getenv("PORT")
    if port == "" { port = "8080" }
    log.Printf("Application serving on port %s", port)
    r.Run(":" + port)
}

func updateBusinessMetrics(db *gorm.DB) {
	// Continuously update metrics for real-time dashboard feel
	for {
		var total float64
		db.Model(&models.Order{}).Where("status = ?", "PAID").Select("SUM(total_price)").Scan(&total)
		// Custom logic could be added here for periodic metrics syncing
		log.Printf("[METRICS] Total Revenue: $%.2f", total)
		prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "retail_current_total_revenue",
			Help: "Current total revenue across all regions",
		}).Set(total)
		
		// In a real app we would use more sophisticated metrics update patterns, 
		// but this works for demo data sync.
		break 
	}
}
