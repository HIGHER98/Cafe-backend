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
	Notes          string `json:"notes"`
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
func AddPurchase(data *Purchase) error {
	logging.Info("Adding purchase: ", data)

	_, err := GetItemById(data.ItemId)
	if err != nil {
		return err
	}

	purchase := Purchase{
		ItemId:         data.ItemId,
		Email:          data.Email,
		Status:         PENDING_TRANSACTION,
		CustName:       data.CustName,
		DateTime:       time.Now().Format("2006-01-02 15:04:05"),
		CollectionTime: data.CollectionTime,
		Notes:          data.Notes,
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
	return purchase, nil
}

func Anonymise(days int) error {
	t := time.Now().AddDate(0, 0, (days * -1)).Format("2006-01-02")
	logging.Debug("Anonymising since: ", t)
	if err := db.Model(Purchase{}).Where("date_time <= ?", t).Updates(Purchase{CustName: "Anonymised", Email: "Anonymised"}).Error; err != nil {
		return err
	}
	return nil
}
