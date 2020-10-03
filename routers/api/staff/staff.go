package staff

import (
	"cafe/models"
	"cafe/pkg/app"
	"cafe/pkg/e"
	"cafe/pkg/logging"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Order struct {
	Id     int `json:"id"`
	Status int `json:"status"`
}

var O chan Order

func InitOrders() {
	O = make(chan Order)
}

//Signals websocket connections of an update to an order
func UpdateOrder(o chan Order, p Order) {
	o <- p
}

func UpdatePurchaseStatus(c *gin.Context) {
	logging.Info("Updating purchase...")
	appG := app.Gin{C: c}
	var p Order
	err := c.Bind(&p)
	if err != nil {
		logging.Debug(err)
		appG.Response(http.StatusBadRequest, e.FAILED_TO_BIND, nil)
		return
	}
	if err = models.UpdatePurchaseStatus(p.Id, p.Status); err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			appG.Response(http.StatusBadRequest, e.ID_NOT_FOUND, nil)
			return
		default:
			logging.Error(err)
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			return
		}
	}
	//If status is pending || confirmed
	if p.Status == 1 || p.Status == 2 {
		go UpdateOrder(O, p)
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func GetPurchases(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
