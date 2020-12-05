package routers

import (
	"cafe/middleware/jwt"
	"cafe/routers/api"
	"cafe/routers/api/admin"
	"cafe/routers/api/cust"
	"cafe/routers/api/staff"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"content-type", "authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.LoadHTMLFiles("public/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	//public routes
	r.GET("/items", cust.GetItemsForSale)
	r.GET("/items/:id", cust.GetItem)
	r.POST("/create-checkout-session", cust.ProcessPayment)
	r.POST("/confirmpayment", cust.PaymentSuccess)

	r.POST("/auth", api.GetAuth)
	//Remove this route for live
	r.POST("/signup", api.CreateUser)

	//Admin routes
	adminApi := r.Group("/api/admin")
	adminApi.Use(jwt.JWT())
	{
		adminApi.GET("/items", admin.GetItems)
		adminApi.POST("/item", admin.AddItem)
		adminApi.PATCH("/item/:id", admin.UpdateItem)
		adminApi.DELETE("/item/:id", admin.DelItem)

		adminApi.GET("/users", admin.GetUsers)
		adminApi.DELETE("/user/:id", admin.DelUser)
		adminApi.PATCH("/user/:id", admin.UpdateUserRole)

		adminApi.GET("/categories", admin.GetCategories)
		adminApi.POST("/category", admin.AddCategory)
		adminApi.PATCH("/category/:id", admin.PatchCategory)
		adminApi.DELETE("/category/:id", admin.DeleteCategory)

		adminApi.GET("/tags", admin.GetTags)
		adminApi.POST("/tag", admin.AddTag)
		adminApi.PATCH("/tag/:id", admin.PatchTag)
		adminApi.DELETE("/tag/:id", admin.DeleteTag)

		adminApi.POST("/option", admin.AddItemOptions)
		adminApi.PATCH("/option/:id", admin.PatchItemOptions)
		adminApi.DELETE("/option/:id", admin.DeleteItemOptions)

		adminApi.POST("/size", admin.AddItemSize)
		adminApi.PATCH("/size/:id", admin.PatchItemSize)
		adminApi.DELETE("/size/:id", admin.DeleteItemSize)

		adminApi.GET("/stats/itempurchasestats", admin.GetItemPurchaseStats)
	}
	r.GET("/ws", func(c *gin.Context) {
		staff.Wshandler(c.Writer, c.Request)
	})
	staffApi := r.Group("/api/staff")
	staffApi.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH"},
		AllowHeaders:     []string{"content-type", "authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	staffApi.Use(jwt.JWT())
	{
		staffApi.GET("/purchases", staff.GetPurchases)
		staffApi.GET("/purchase/:id", staff.GetPurchaseById)
		staffApi.GET("/staff", api.GetStaffMembers)
		staffApi.GET("/status", staff.GetStatuses)
	}

	return r
}
