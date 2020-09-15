package routers

import (
	"bytes"
	"cafe/models"
	"cafe/pkg/logging"
	"cafe/pkg/setting"
	"cafe/pkg/util"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func init() {
	setting.Setup()
	models.Setup()
	util.Setup()
	logging.Setup()
	router = InitRouter()
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var resp Response

func Test_GetItemsForSale(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/items", nil)
	if err != nil {
		t.Errorf("Failed to get items: %v", err)
	}
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	err = json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Errorf("Failed to marshal response to json: %v", err)
	}
	assert.Equal(t, "ok", resp.Msg)
}

func Test_GetItem(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/items/1", nil)
	if err != nil {
		t.Errorf("Failed to get item: %v", err)
	}
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	err = json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Errorf("Failed to marshal response to json: %v", err)
	}
	assert.Equal(t, "ok", resp.Msg)
}

type Purchase struct {
	ItemId         int    `json:"item_id"`
	Email          string `json:"email"`
	CustName       string `json:"cust_name"`
	CollectionTime string `json:"collection_time"`
}

func Test_SubmitDetails(t *testing.T) {
	purchase := &Purchase{ItemId: 1, Email: "test@test.com", CustName: "Test Testington", CollectionTime: "2020-09-20"}
	jsonPurchase, err := json.Marshal(purchase)
	if err != nil {
		t.Errorf("Error martialling: %v", err)
	}
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/purchase", bytes.NewBuffer(jsonPurchase))
	if err != nil {
		t.Errorf("Failed to make request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	err = json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Errorf("Failed to marshal response to json: %v", err)
	}
	assert.Equal(t, "Created", resp.Msg)
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

type T struct {
	Token string `json:"token"`
}

var j T

/*
func Test_CreateUser(t *testing.T) {
	user := &User{Username: "testy", Password: "test", Role: 1}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		t.Errorf("Failed to marshal json: %v", err)
	}
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonUser))
	if err != nil {
		t.Errorf("Failed to make request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	err = json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Errorf("Failed to marshal response to json: %v", err)
	}
	assert.Equal(t, "Created", resp.Msg)
}
*/

func Test_SignIn(t *testing.T) {
	user := &User{Username: "testy", Password: "test"}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		t.Errorf("Failed to marshal json: %v", err)
	}
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/auth", bytes.NewBuffer(jsonUser))
	if err != nil {
		t.Errorf("Failed to make request: %v", err)
	}
	req.Header.Set("Context-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	err = json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Errorf("Failed to marshal response to json: %v", err)
	}
	assert.Equal(t, "ok", resp.Msg)
	j.Token = resp.Data.(map[string]interface{})["token"].(string)
}

func Test_AdminGetItems(t *testing.T) {
	w := httptest.NewRecorder()
	jsonToken, err := json.Marshal(j)
	if err != nil {
		t.Errorf("Failed to marshal json: %v", err)
	}
	req, err := http.NewRequest("GET", "/admin/api/items", bytes.NewBuffer(jsonToken))
	if err != nil {
		t.Errorf("Failed to get item: %v", err)
	}
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	err = json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Errorf("Failed to marshal response to json: %v", err)
	}
	assert.Equal(t, "ok", resp.Msg)
}

/*
type Item struct {
	Name        string  `form:"name"`
	Description string  `form:"description"`
	Price       float32 `form:"price"`
}

func Test_AddItem(t *testing.T) {
	item := &Item{Name: "", Description: "", Price: 0.00}
	w := httptest.NewRecorder()
	jsonToken, err := json.Marshal(j)
	if err != nil {
		t.Errorf("Failed to marshal json: %v", err)
	}
	req, err := http.NewRequest("POST", "/admin/api/additem", bytes.NewBuffer(jsonToken))
	if err != nil {
		t.Errorf("Failed to get item: %v", err)
	}
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
*/
