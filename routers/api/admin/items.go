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

func AddCategory(c *gin.Context) {
	appG := app.Gin{C: c}
	var category models.Category
	err := c.Bind(&category)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}
	err = models.AddCategory(category.Name)
	if err != nil {
		logging.Error("Failed to add category: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusCreated, e.CREATED, nil)
}

func PatchCategory(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}
	var category models.Category
	err = c.Bind(&category)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}
	err = models.EditCategory(id, category.Name)
	if err != nil {
		logging.Error("Failed to edit category: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//TODO
func DeleteCategories(c *gin.Context) {

}

func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}
	tags, err := models.GetAllTags()
	if err != nil {
		logging.Error("Failed to get tags: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, tags)
}

func AddTag(c *gin.Context) {
	appG := app.Gin{C: c}
	var tag models.Tag
	err := c.Bind(&tag)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}
	err = models.AddTag(tag.Name)
	if err != nil {
		logging.Error("Failed to add tag: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusCreated, e.CREATED, nil)
}

func PatchTag(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}
	var tag models.Tag
	err = c.Bind(&tag)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}
	err = models.EditTag(id, tag.Name)
	if err != nil {
		logging.Error("Failed to edit tag: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func DeleteTag(c *gin.Context) {}

//Add/Delete tags for a specific item
func AddItemTags(c *gin.Context)    {}
func DeleteItemTags(c *gin.Context) {}

//Should these do options as well??
func AddItemOptions(c *gin.Context)    {}
func PatchItemOptions(c *gin.Context)  {}
func DeleteItemOptions(c *gin.Context) {}
