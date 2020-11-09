package models

import (
	"cafe/pkg/logging"
	"errors"
	"time"
)

type PurchaseItems struct {
	Id         int `json:"id"`
	PurchaseId int `json:"purchase_id"`
	ItemId     int `json:"item_id"`
	SizeId     int `gorm:"default:null" json:"size_id"`
	OptId      int `gorm:"default:null" json:"opt_id"`
}

type Purchase struct {
	Id             int             `json:"id"`
	Item           []PurchaseItems `json:"items"`
	Email          string          `json:"email"`
	CustName       string          `json:"cust_name"`
	DateTime       string          `json:"date_time"`
	CollectionTime string          `json:"collection_time"`
	Status         int             `json:"status"` //Status options: 'Pending transaction', 'Pending', 'Confirmed', 'Collected'
	Notes          string          `json:"notes"`
	Uuid           string          `json:"uuid"`
}

type PurchaseViews struct {
	PurchasesId     int     `json:"purchases_id"`
	CustName        string  `json:"cust_name"`
	Email           string  `json:"email"`
	Cost            float64 `json:"cost"`
	DateTime        string  `json:"date_time"`
	CollectionTime  string  `json:"collection_time"`
	Notes           string  `json:"notes"`
	ItemName        string  `json:"item_name"`
	Opt             string  `json:"opt"`
	ItemSize        string  `json:"item_size"`
	PurchaseItemsId int     `json:"purchase_items_id"`
	Status          string  `json:"status"`
}

//SELECT * FROM purchase;
func GetAllPurchases() ([]*Purchase, error) {
	var purchases []*Purchase
	if err := db.Find(&purchases).Error; err != nil {
		return nil, err
	}
	return purchases, nil
}

//SELECT * FROM purchase WHERE id=?;
func GetPurchaseById(id int) (*Purchase, error) {
	var purchase Purchase
	if err := db.Where("id=?", id).First(&purchase).Error; err != nil {
		return nil, err
	}
	return &purchase, nil
}

//INSERT INTO purchase (email, status, cust_name, date_time, collection_time) VALUES (?, ?, 1, ?, ?);
func AddPurchase(data *Purchase) (int, error) {
	logging.Info("Adding purchase: ", *data)
	for _, j := range data.Item {
		err := ValidItemOrder(j.ItemId, j.SizeId, j.OptId)
		if err != nil {
			return 0, err
		}
	}
	purchase := Purchase{
		Email:          data.Email,
		Status:         PENDING_TRANSACTION,
		CustName:       data.CustName,
		DateTime:       time.Now().Format("2006-01-02 15:04:05"),
		CollectionTime: data.CollectionTime,
		Notes:          data.Notes,
		Uuid:           data.Uuid,
	}
	if err := db.Create(&purchase).Error; err != nil {
		return 0, err
	}
	return purchase.Id, nil
}

//Returns ID, #rowsAffected, error
func ConfirmPurchase(uuid string) (int, int64, error) {
	logging.Info("Confirming purchase with UUID: ", uuid)
	var purchase Purchase
	result := db.Model(&purchase).Where("uuid = ? AND status=1", uuid).Update("status", 2)
	if result.Error != nil {
		return 0, 0, result.Error
	}
	if err := db.Where("uuid = ?", uuid).Find(&purchase).Error; err != nil {
		return 0, 0, err
	}
	return purchase.Id, result.RowsAffected, nil
}

func AddPurchaseItems(purchaseId int, items []PurchaseItems) error {
	logging.Info("Adding items for order: ", purchaseId)
	for i := range items {
		items[i].PurchaseId = purchaseId
	}
	logging.Info("Adding items: ", items, " for purchase ", purchaseId)
	if err := db.Create(items).Error; err != nil {
		return err
	}
	return nil
}

//UPDATE purchase SET (status=?) WHERE id=?
func UpdatePurchaseStatus(id, status int) error {
	purchase, err := GetPurchaseById(id)
	if err != nil {
		return err
	}
	s, err := GetStatus(status)
	if err != nil {
		return err
	} else if s == nil {
		return errors.New("Invalid status")
	}

	if err = db.Model(&purchase).Where("id=?", id).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

//Returns a list of orders for todays date that are pending/collected
//SELECT * FROM purchase WHERE collection_time > today AND collection_time < tomorrow AND status IN (2,3)
func GetTodaysOrders() ([]*Purchase, error) {
	var purchase []*Purchase
	t := time.Now().Format("2006-01-02")
	t1 := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	if err := db.Where("collection_time > ? AND collection_time < ? AND status IN (?)", t, t1, []int{2, 3}).Find(&purchase).Error; err != nil {
		return nil, err
	}
	return purchase, nil
}

//SELECT * FROM purchase_views WHERE id=?;
func GetItemsFromPurchaseView(id int) ([]*PurchaseViews, error) {
	var purchaseItems []*PurchaseViews
	if err := db.Where("purchases_id=?", id).Find(&purchaseItems).Error; err != nil {
		return nil, err
	}
	return purchaseItems, nil
}

//SELECT * FROM purchase_views;
func GetAllTodayOrActivePurchasesView() ([]*PurchaseViews, error) {
	var purchaseItems []*PurchaseViews
	today := "%" + time.Now().Format("2006-01-02") + "%"
	if err := db.Where("collection_time LIKE ? AND status NOT LIKE 'Pending transaction'", today).Find(&purchaseItems).Error; err != nil {
		return nil, err
	}
	return purchaseItems, nil
}

//Anonymises data in the purchase table after `days` old
//UPDATE purchase SET cust_name=Anonymised AND email=Anonymised WHERE date_time <= today-days
func Anonymise(days int) error {
	t := time.Now().AddDate(0, 0, (days * -1)).Format("2006-01-02")
	logging.Debug("Anonymising since: ", t)
	if err := db.Model(Purchase{}).Where("date_time <= ?", t).Updates(Purchase{CustName: "Anonymised", Email: "Anonymised"}).Error; err != nil {
		return err
	}
	return nil
}

type PurchaseActivities struct {
	Id         int    `json:"id"`
	PurchaseId int    `json:"purchase_id"`
	StatusSet  int    `json:"status_set"`
	SetBy      int    `json:"set_by"`
	UpdatedAt  string `json:"updated_at"`
}

func AddPurchaseActivity(purchaseId, status, setBy int) error {
	purchaseActivity := PurchaseActivities{PurchaseId: purchaseId, StatusSet: status, SetBy: setBy, UpdatedAt: time.Now().Format("2006-01-02 15:04:05")}
	if err := db.Create(&purchaseActivity).Error; err != nil {
		return err
	}
	return nil
}
