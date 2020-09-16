package models

import (
	"cafe/pkg/logging"
	"time"
)

type Purchase struct {
	Id     int    `json:"id"`
	ItemId int    `json:"item_id"`
	Email  string `json:"email"`
	//Status options: 'Pending transaction', 'Pending', 'Confirmed', 'Delivered'
	Status         int    `json:"status"`
	CustName       string `json:"cust_name"`
	DateTime       string `json:"date_time"`
	CollectionTime string `json:"collection_time"`
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

//INSERT INTO purchase (item_id, email, status, cust_name, date_time, collection_time) VALUES (?, ?, 1, ?, ?);
func AddPurchase(data map[string]interface{}) error {
	var id int
	if itemid, ok := data["item_id"].(int); ok {
		id = itemid
	} else if itemid, ok := data["item_id"].(float64); ok {
		id = int(itemid)
	}
	item, err := GetItemById(id)
	if err != nil || item == nil {
		return err
	}

	purchase := Purchase{
		ItemId:         id,
		Email:          data["email"].(string),
		Status:         PENDING_TRANSACTION,
		CustName:       data["cust_name"].(string),
		DateTime:       time.Now().Format("2006-01-02 15:04:05"),
		CollectionTime: data["collection_time"].(string),
	}
	if err := db.Create(&purchase).Error; err != nil {
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
	if s == nil || err != nil {
		return err
	}

	if err = db.Model(&purchase).Where("id=?", id).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func GetTodaysOrders() ([]*Purchase, error) {
	var purchase []*Purchase
	t := time.Now().Format("2006-01-02")
	t1 := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	if err := db.Where("collection_time > ? AND collection_time < ? AND status IN (?)", t, t1, []int{2, 3}).Find(&purchase).Error; err != nil {
		return nil, err
	}
	logging.Info(purchase)
	return purchase, nil
}
