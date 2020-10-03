package cust

import (
	"cafe/models"
	"cafe/pkg/app"
	"cafe/pkg/e"
	"cafe/pkg/logging"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//Get a single item by ID
func GetItem(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}
	item, err := models.GetItemViewById(id)
	if err == gorm.ErrRecordNotFound {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	} else if err != nil {
		logging.Error(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, item)
}

//Get all items for sale
func GetItemsForSale(c *gin.Context) {
	appG := app.Gin{C: c}
	items, err := models.GetAllActiveItems()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.NOT_FOUND, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, items)
}

//Submit details for purchasing an item
func SubmitDetails(c *gin.Context) {
	appG := app.Gin{C: c}
	var purchase models.Purchase
	err := c.Bind(&purchase)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.FAILED_TO_BIND, nil)
		return
	}

	id, err := models.AddPurchase(&purchase)
	if err == gorm.ErrRecordNotFound {
		appG.Response(http.StatusBadRequest, e.ID_NOT_FOUND, nil)
		return
	} else if err != nil {
		logging.Error("Failed to add purchase: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	} else if id == 0 {
		logging.Error("Failed to add purchase. Id returned 0")
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	logging.Info("Adding purchase items for id: ", id)
	err = models.AddPurchaseItems(id, purchase.Item)
	if err != nil {
		logging.Error("Failed to add purchase data: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		//Need to delete purchase by id here.
		return
	}
	appG.Response(http.StatusCreated, e.CREATED, nil)
}

//Process payment - Speak to paypal - Marks item as sold
func ProcessPayment(c *gin.Context) {

}
