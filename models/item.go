package models

import (
	"fmt"
	"time"

	"cafe/pkg/logging"
	"gorm.io/gorm"
)

type Item struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    int     `json:"category"`
	Tag         int     `json:"tag"`
	UploadDate  string  `json:"upload_date"`
	IsDel       int     `json:"is_del"`
}

//INSERT INTO item (name, description, price, del_price, upload_date) VALUES (?, ?, ?, ?, ?)
func AddItem(data Item) error {
	/*item := Item{
		Name:        data["Name"].(string),
		Description: data["Description"].(string),
		Price:       data["Price"].(float64),
		Category:    data["Category"].(int),
		Tag:         data["Tag"].(int),
		UploadDate:  time.Now().Format("2006-01-02"),
		IsDel:       0,
	}*/
	data.UploadDate = time.Now().Format("2006-01-02")
	if err := db.Create(&data).Error; err != nil {
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
func SetItemIsDel(id int) error {
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
func UpdateItem(id int, data Item) error {
	/*
		item := Item{
			Name:        data["Name"].(string),
			Description: data["Description"].(string),
			Price:       data["Price"].(float64),
			Category:    data["Category"].(int),
			Tag:         data["Tag"].(int),
			UploadDate:  time.Now().Format("2006-01-02"),
		}*/
	data.UploadDate = time.Now().Format("2006-01-02")
	if err := db.Model(&data).Where("id=?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

//Using new Database structure

type ItemView struct {
	Id          int     `json:"id"`
	ItemName    string  `json:"item_name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	OptId       int     `json:"opt_id"`
	Opt         string  `json:"opt"`
	OptionPrice float64 `json:"option_price"`
	SizeId      int     `json:"size_id"`
	ItemSize    string  `json:"item_size"`
	SizePrice   float64 `json:"size_price"`
	Category    string  `json:"category"`
	Tag         string  `json:"tag"`
}

func GetAllActiveItems() ([]*ItemView, error) {
	var items []*ItemView
	if err := db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

//Returns the item and the available options & sizing options if they exist
func GetItemViewById(id int) ([]*ItemView, error) {
	var item []*ItemView
	if err := db.Where("id = ?", id).Find(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func ValidItemOrder(itemId, sizeId, optId interface{}) error {
	logging.Info("Checking validity for ID: ", itemId, "\nSizeid: ", sizeId, "\noptId: ", optId)

	if optId.(int) == 0 && sizeId.(int) == 0 {
		if err := db.Where("id=?", itemId.(int)).First(&ItemView{}).Error; err != nil {
			logging.Error(err)
			return err
		}
	} else if optId.(int) == 0 && sizeId.(int) != 0 {
		if err := db.Where("id = ? AND size_id = ?", itemId.(int), sizeId.(int)).First(&ItemView{}).Error; err != nil {
			logging.Error(err)
			return err
		}
	} else if sizeId.(int) == 0 && optId.(int) != 0 {
		if err := db.Where("id = ? AND opt_id = ?", itemId.(int), optId.(int)).First(&ItemView{}).Error; err != nil {
			logging.Error(err)
			return err
		}
	} else {
		if err := db.Where("id = ? AND opt_id = ? AND size_id = ?", itemId.(int), optId.(int), sizeId.(int)).First(&ItemView{}).Error; err != nil {
			logging.Error("Unable to find for ", itemId, " ", sizeId, " ", optId, err)
			return err
		}
	}
	return nil
}
