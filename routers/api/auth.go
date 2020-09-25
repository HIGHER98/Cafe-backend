package api

import (
	"encoding/json"
	"net/http"

	"cafe/models"
	"cafe/pkg/app"
	"cafe/pkg/e"
	"cafe/pkg/logging"
	"cafe/pkg/util"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

//Sign in user and return JWT token
func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	var user User
	err := c.Bind(&user)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.FAILED_TO_BIND, nil)
		return
	}

	var userInterface map[string]interface{}
	inrec, err := json.Marshal(user)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.MARSHAL_ERROR, nil)
		return
	}
	json.Unmarshal(inrec, &userInterface)
	//Below works for test but not Insomnia
	/*data, err := c.GetRawData()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	var userMap map[string]string
	json.Unmarshal(data, &userMap)*/
	s, err := models.Check(userInterface["username"].(string), userInterface["password"].(string))
	if err != nil || !s {
		appG.Response(http.StatusUnauthorized, e.UNAUTHORIZED, nil)
		return
	}

	token, err := util.GenerateToken(userInterface["username"].(string), userInterface["password"].(string))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		logging.Error(err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}

//Create user, if they have supplied a unique username
func CreateUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var user User
	err := c.Bind(&user)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.FAILED_TO_BIND, nil)
		return
	}

	var userInterface map[string]interface{}
	inrec, err := json.Marshal(user)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.MARSHAL_ERROR, nil)
		return
	}
	json.Unmarshal(inrec, &userInterface)

	exists, err := models.CheckUserExists(userInterface["username"].(string))
	if exists {
		appG.Response(http.StatusBadRequest, e.USERNAME_TAKEN, nil)
		return
	} else if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	s, err := models.SignupUser(userInterface)
	if !s || err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		logging.Warn("%v", err)
		return
	}
	appG.Response(http.StatusCreated, e.CREATED, nil)
	//Frontend should redirect to GetAuth function on 201 code
}
