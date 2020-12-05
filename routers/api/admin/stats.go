package admin

import (
	"cafe/models"
	"cafe/pkg/app"
	"cafe/pkg/e"
	"cafe/pkg/logging"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const layout = "2006-01-02"

func GetItemPurchaseStats(c *gin.Context) {
	appG := app.Gin{C: c}
	from, ok := c.GetQuery("from")
	if !ok {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}
	until, ok := c.GetQuery("until")
	if !ok {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}

	t1, err := time.Parse(layout, from)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}
	t2, err := time.Parse(layout, until)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}
	fromDate := t1.Format(layout)
	untilDate := t2.Format(layout)

	stats, err := models.ItemPurchaseStats(fromDate, untilDate)
	if err != nil {
		logging.Error("Failed to get item purchase stats: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, stats)
}
