package admin

import (
	"cafe/models"
	"cafe/pkg/app"
	"cafe/pkg/e"
	"cafe/pkg/logging"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Item struct {
	Name        string  `form:"name"`
	Description string  `form:"description"`
	Price       float32 `form:"price"`
}

//Get all live items
func GetItems(c *gin.Context) {
	appG := app.Gin{C: c}
	items, err := models.GetItemsForSale()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, items)
}

//Update specific item
func UpdateItem(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}

	var item Item
	err = c.Bind(&item)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}
	var itemInterface map[string]interface{}
	inrec, err := json.Marshal(item)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}
	json.Unmarshal(inrec, &itemInterface)
	logging.Debug("Updating item with: ", itemInterface)

	err = models.UpdateItem(id, itemInterface)
	if err == gorm.ErrRecordNotFound {
		appG.Response(http.StatusOK, e.ID_NOT_FOUND, nil)
		return
	} else if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		logging.Error(err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//Add item to the menu
func AddItem(c *gin.Context) {
	appG := app.Gin{C: c}
	var item Item
	err := c.Bind(&item)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}
	var itemInterface map[string]interface{}
	inrec, err := json.Marshal(item)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.MARSHAL_ERROR, nil)
		return
	}
	json.Unmarshal(inrec, &itemInterface)
	logging.Debug("Adding item: ", itemInterface)
	err = models.AddItem(itemInterface)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		logging.Warn("Failed to create item: ", itemInterface, "\nError: ", err)
		return
	}
	appG.Response(http.StatusCreated, e.CREATED, nil)
}

func DelItem(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}
	err = models.SetIsDel(id)
	if err == gorm.ErrRecordNotFound {
		appG.Response(http.StatusBadRequest, e.ID_NOT_FOUND, nil)
		return
	} else if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		logging.Error(err)
		return
	}
	appG.Response(http.StatusOK, e.DELETED, nil)
}

//VIEW/EDIT/DELETE/ADD staff
