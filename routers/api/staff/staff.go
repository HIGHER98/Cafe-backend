package staff

import (
	"cafe/models"
	"cafe/pkg/app"
	"cafe/pkg/e"
	"cafe/pkg/logging"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	logging.Debug("An order is being updated: Channel: ", o, " Order: ", p)
	O <- p
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
	//If status is not pending payment
	if p.Status != 1 {
		go UpdateOrder(O, p)
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func GetPurchases(c *gin.Context) {
	appG := app.Gin{C: c}
	purchases, err := models.GetAllTodayOrActivePurchasesView()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, purchases)
}

func GetPurchaseById(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}
	purchase, err := models.GetItemsFromPurchaseView(id)
	if err == gorm.ErrRecordNotFound {
		appG.Response(http.StatusOK, e.ID_NOT_FOUND, nil)
	} else if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
	}
	appG.Response(http.StatusOK, e.SUCCESS, purchase)
}
