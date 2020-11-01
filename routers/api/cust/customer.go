package cust

import (
	"cafe/models"
	"cafe/pkg/app"
	"cafe/pkg/e"
	"cafe/pkg/logging"
	"fmt"
	"net/http"
	"strconv"

	"cafe/routers/api/staff"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/checkout/session"
	"gorm.io/gorm"
)

//Get a single item by ID
func GetItem(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.BAD_REQUEST, nil)
		return
	}
	item, err := models.GetItemViewById(id)
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
	items, err := models.GetAllActiveItems()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.NOT_FOUND, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, items)
}

//Submit details for purchasing an item
func SubmitDetails(c *gin.Context) {
	appG := app.Gin{C: c}
	var purchase models.Purchase
	err := c.Bind(&purchase)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.FAILED_TO_BIND, nil)
		return
	}

	id, err := models.AddPurchase(&purchase)
	if err == gorm.ErrRecordNotFound {
		appG.Response(http.StatusBadRequest, e.ID_NOT_FOUND, nil)
		return
	} else if err != nil {
		logging.Error("Failed to add purchase: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	} else if id == 0 {
		logging.Error("Failed to add purchase. Id returned 0")
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	logging.Info("Adding purchase items for id: ", id)
	err = models.AddPurchaseItems(id, purchase.Item)
	if err != nil {
		logging.Error("Failed to add purchase data: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		//Need to delete purchase by id here.
		return
	}
	appG.Response(http.StatusCreated, e.CREATED, nil)
}

//Submit details for purchasing an item
func SubmitDetailsHelper(purchase *models.Purchase) (int, int, int) {
	id, err := models.AddPurchase(purchase)
	if err == gorm.ErrRecordNotFound {
		return 0, http.StatusBadRequest, e.ID_NOT_FOUND
	} else if err != nil {
		logging.Error("Failed to add purchase: ", err)
		return 0, http.StatusInternalServerError, e.ERROR
	} else if id == 0 {
		logging.Error("Failed to add purchase. Id returned 0")
		return 0, http.StatusInternalServerError, e.ERROR
	}
	logging.Info("Adding purchase items for id: ", id)
	err = models.AddPurchaseItems(id, purchase.Item)
	if err != nil {
		logging.Error("Failed to add purchase data: ", err)
		//Need to delete purchase by id here.
		return 0, http.StatusInternalServerError, e.ERROR
	}
	return id, http.StatusCreated, e.CREATED
}

type CreateCheckoutSessionResponse struct {
	SessionID string `json:"id"`
}

func makeStripeLineItemList(purchaseItems []*models.PurchaseViews) ([]*stripe.CheckoutSessionLineItemParams, error) {
	var s []*stripe.CheckoutSessionLineItemParams
	for _, j := range purchaseItems {
		logging.Info("Iterating through purchaseItems: ", j)
		s = append(s,
			&stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("eur"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(j.ItemName + " " + j.Opt + " " + j.ItemSize),
					},
					UnitAmount: stripe.Int64(int64(j.Cost * 100)),
				},
				Quantity: stripe.Int64(1),
			},
		)
	}
	return s, nil
}

//Process payment - Speak to stripe - Marks item as sold
func ProcessPayment(c *gin.Context) {
	appG := app.Gin{C: c}
	var purchase models.Purchase
	err := c.Bind(&purchase)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.FAILED_TO_BIND, nil)
		return
	}
	uuid := uuid.New().String()
	purchase.Uuid = uuid
	purchaseId, status, code := SubmitDetailsHelper(&purchase)
	if purchaseId == 0 || status != http.StatusCreated || code != e.CREATED {
		appG.Response(status, code, nil)
		return
	}

	p, err := models.GetItemsFromPurchaseView(purchaseId)
	if err != nil {
		logging.Error("Failed to get items from purchase_view: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	s, _ := makeStripeLineItemList(p)
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems:  s,
		SuccessURL: stripe.String("http://192.168.8.105:5000/success/" + uuid),
		CancelURL:  stripe.String("http://192.168.8.105:5000/cancel"),
	}

	session, err := session.New(params)
	if err != nil {
		logging.Error("Failed to create Stripe session: ", err)
		appG.Response(http.StatusInternalServerError, e.STRIPE_CREATE_SESSION_ERROR, nil)
		return
	}

	data := CreateCheckoutSessionResponse{
		SessionID: session.ID,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

type paymentSuccess struct {
	Uuid string `json:"uuid"`
}

func PaymentSuccess(c *gin.Context) {
	appG := app.Gin{C: c}
	var uuid paymentSuccess
	err := c.Bind(&uuid)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.FAILED_TO_BIND, nil)
		return
	}
	id, rowsAffected, err := models.ConfirmPurchase(uuid.Uuid)
	if err != nil {
		logging.Error("Failed to confirm purchase with UUID: ", uuid.Uuid, "\t Error: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	} else if id == 0 {
		logging.Error("Id returned 0 when confirming purchase")
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	//If it's a new purchase and not just someone refreshing the page
	//TODO: Only send if it's an order for today
	if rowsAffected != 0 {
		fmt.Println("Before updating order in customer")
		go func(id int) {
			order := staff.Order{Id: id, Status: 1}
			staff.UpdateOrder(staff.O, order)
		}(id)

	}
	fmt.Println("After updating order in customer")
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
