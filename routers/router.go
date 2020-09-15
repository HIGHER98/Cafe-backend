package routers

import (
	"cafe/middleware/jwt"
	"cafe/routers/api"
	"cafe/routers/api/admin"
	"cafe/routers/api/cust"
	"cafe/routers/api/staff/websocket"

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
	adminApi := r.Group("/admin/api")
	adminApi.Use(jwt.JWT())
	{
		adminApi.GET("/items", admin.GetItems)
		adminApi.PATCH("/item/:id", admin.UpdateItem)
		adminApi.POST("/additem", admin.AddItem)
		adminApi.POST("/delitem/:id", admin.DelItem)
	}
	r.GET("/wsqueue", func(c *gin.Context) {
		websocket.Wshandler(c.Writer, c.Request)
	})
	staffApi := r.Group("staff/api")
	staffApi.Use(jwt.JWT())
	{

	}

	return r
}
