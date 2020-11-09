package staff

import (
	"cafe/models"
	"cafe/pkg/app"
	"cafe/pkg/e"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
