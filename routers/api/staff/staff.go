package staff

import (
	"cafe/models"
	"cafe/pkg/app"
	"cafe/pkg/e"
	"cafe/pkg/logging"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type purchase struct {
	Id     int `json:"id"`
	Status int `json:"status"`
}

func UpdatePurchaseStatus(c *gin.Context) {
	appG := app.Gin{C: c}
	var p purchase
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
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func GetPurchases(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
