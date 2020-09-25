package routers

import (
	"cafe/middleware/jwt"
	"cafe/routers/api"
	"cafe/routers/api/admin"
	"cafe/routers/api/cust"
	"cafe/routers/api/staff"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.LoadHTMLFiles("public/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	//public routes
	r.GET("/items", cust.GetItemsForSale)
	r.GET("/items/:id", cust.GetItem)
	r.POST("/purchase", cust.SubmitDetails)

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
	}
	r.GET("/ws", func(c *gin.Context) {
		staff.Wshandler(c.Writer, c.Request)
	})
	staffApi := r.Group("/api/staff")
	staffApi.Use(jwt.JWT())
	{
		staffApi.GET("/purchases", staff.GetPurchases)
		staffApi.PATCH("/purchase", staff.UpdatePurchaseStatus)
		/*	staffApi.GET("/ws", func(c *gin.Context){
			staff.Wshandler(c.Writer, c.Request)
		})*/

	}

	return r
}
