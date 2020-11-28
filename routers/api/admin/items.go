package admin

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

type Item struct {
	Name        string  `form:"name"`
	Description string  `form:"description"`
	Price       float32 `form:"price"`
	Category    int     `form:"category"`
	Tag         int     `form:"tag"`
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

	var item models.Item
	err = c.Bind(&item)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}

	err = models.UpdateItem(id, item)
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
	var item models.Item
	err := c.Bind(&item)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}

	logging.Debug("Adding item: ", item)
	err = models.AddItem(item)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		logging.Warn("Failed to create item: ", item, "\nError: ", err)
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
	err = models.SetItemIsDel(id)
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

func GetCategories(c *gin.Context) {
	appG := app.Gin{C: c}
	categories, err := models.GetAllCategories()
	if err != nil {
		logging.Error("Failed to get categories: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, categories)
}

func AddCategories(c *gin.Context)    {}
func PatchCategories(c *gin.Context)  {}
func DeleteCategories(c *gin.Context) {}

func GetTags(c *gin.Context)    {}
func AddTags(c *gin.Context)    {}
func PatchTags(c *gin.Context)  {}
func DeleteTags(c *gin.Context) {}

func AddItemTags(c *gin.Context)    {}
func DeleteItemTags(c *gin.Context) {}

//Should these do options as well
func AddItemOptions(c *gin.Context)    {}
func PatchItemOptions(c *gin.Context)  {}
func DeleteItemOptions(c *gin.Context) {}
