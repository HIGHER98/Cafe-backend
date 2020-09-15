package cust

import (
	"cafe/models"
	"cafe/pkg/app"
	"cafe/pkg/e"
	"cafe/pkg/logging"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//Get a single item by ID
func GetItem(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}
	item, err := models.GetItemById(id)
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
	items, err := models.GetItemsForSale()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.NOT_FOUND, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, items)
}

type Purchase struct {
	ItemID         int    `json:"item_id"`
	Email          string `json:"email"`
	CustName       string `json:"cust_name"`
	CollectionTime string `json:"collection_time"`
}

//Submit details for purchasing an item
func SubmitDetails(c *gin.Context) {
	appG := app.Gin{C: c}
	var purchase Purchase
	err := c.Bind(&purchase)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.FAILED_TO_BIND, nil)
		return
	}
	var purchaseInterface map[string]interface{}
	inrec, err := json.Marshal(purchase)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.MARSHAL_ERROR, nil)
		return
	}
	json.Unmarshal(inrec, &purchaseInterface)

	err = models.AddPurchase(purchaseInterface)
	if err != nil {
		log.Fatalf("%v", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusCreated, e.CREATED, nil)
}

//Process payment - Speak to paypal - Marks item as sold
func ProcessPayment(c *gin.Context) {

}
