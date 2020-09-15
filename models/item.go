package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Item struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	UploadDate  string  `json:"upload_date"`
	IsDel       int     `json:"is_del"`
}

//INSERT INTO item (name, description, price, del_price, upload_date) VALUES (?, ?, ?, ?, ?)
func AddItem(data map[string]interface{}) error {
	item := Item{
		Name:        data["Name"].(string),
		Description: data["Description"].(string),
		Price:       data["Price"].(float64),
		UploadDate:  time.Now().Format("2006-01-02"),
		IsDel:       0,
	}
	if err := db.Create(&item).Error; err != nil {
		return err
	}
	return nil
}

//SELECT * FROM item;
func GetItems() []*Item {
	var items []*Item
	db.Find(&items)
	return items
}

//SELECT * FROM item WHERE ID=? ORDER BY id LIMIT 1;
func GetItemById(id int) (*Item, error) {
	var item Item
	if err := db.Where("id=? AND is_del=0", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

//SELECT * FROM item WHERE is_del=0;
func GetItemsForSale() ([]*Item, error) {
	var items []*Item
	if err := db.Where("is_del=?", 0).Find(&items).Error; err != nil {
		return nil, err
	}
	fmt.Printf("%v", items)
	return items, nil
}

//UPDATE item SET is_del=1 WHERE id=?;
func SetIsDel(id int) error {
	item, err := GetItemById(id)
	if err != nil || err == gorm.ErrRecordNotFound {
		return err
	}
	if err = db.Model(&item).Update("is_del", 1).Error; err != nil {
		return err
	}
	return nil
}

//UPDATE item SET name=?, description=?, price=? WHERE id=?;
func UpdateItem(id int, data map[string]interface{}) error {
	item := Item{
		Name:        data["Name"].(string),
		Description: data["Description"].(string),
		Price:       data["Price"].(float64),
		UploadDate:  time.Now().Format("2006-01-02"),
	}
	if err := db.Model(&item).Where("id = ?", id).Update(item).Error; err != nil {
		return err
	}
	return nil
}
