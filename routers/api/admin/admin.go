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

//VIEW/EDIT/DELETE/ADD staff
func GetUsers(c *gin.Context) {
	appG := app.Gin{C: c}
	users, err := models.GetActiveUsers()
	if err != nil {
		logging.Error(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
	}
	appG.Response(http.StatusOK, e.SUCCESS, users)
}

//Sets is_del=0
func DelUser(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logging.Debug(id)
		appG.Response(http.StatusBadRequest, e.FAILED_ATOI, nil)
		return
	}
	if err = models.MarkUserDeleted(id); err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			appG.Response(http.StatusBadRequest, e.ID_NOT_FOUND, nil)
			return
		default:
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			logging.Error(err)
			return
		}
	}
	appG.Response(http.StatusOK, e.DELETED, nil)
}

type role struct {
	Role int `json:"role"`
}

//Updates user
func UpdateUserRole(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logging.Debug(id)
		appG.Response(http.StatusBadRequest, e.FAILED_ATOI, nil)
		return
	}
	var r role
	err = c.Bind(&r)
	if err != nil {
		logging.Debug(err)
		appG.Response(http.StatusBadRequest, e.FAILED_TO_BIND, nil)
		return
	}

	if err = models.UpdateUser(id, r.Role); err != nil {
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
