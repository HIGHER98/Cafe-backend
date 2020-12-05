package models

type ItemSizes struct {
	Id          int      `json:"size_id,omitempty" gorm:"column:id"`
	ItemId      int      `json:"item_id"`
	ItemSize    string   `json:"item_size"`
	AddPrice    *float64 `json:"size_price" gorm:"column:"add_price" default:0"`
	Description *string  `json:"size_description" gorm:"default:null column:description"`
}

//INSERT INTO item_sizes VALUES (item_id, item_size, add_price, description) VALUES (?, ?, ?, ?);
func (size *ItemSizes) AddItemSize() error {
	if err := db.Create(size).Error; err != nil {
		return err
	}
	return nil
}

//SELECT * FROM item_sizes WHERE id = ?;
func GetItemSize(id int) (*ItemSizes, error) {
	var size ItemSizes
	if err := db.Where("id=?", id).First(&size).Error; err != nil {
		return nil, err
	}
	return &size, nil
}

//UPDATE item_sizes SET item_size=?, add_price-?, description=? WHERE id=sizeId;
func (size *ItemSizes) UpdateItemSize(id int) error {
	if err := db.Model(size).Where("id=?", id).Updates(&size).Error; err != nil {
		return err
	}
	return nil
}

//UPDATE item_sizes SET is_del=1 WHERE id=sizeId;
func DeleteItemSize(id int) error {
	var size *ItemSizes
	if err := db.Model(&size).Where("id=?", id).Update("is_del", 1).Error; err != nil {
		return err
	}
	return nil
}

//UPDATE item_sizes SET is_del=1 WHERE item_id=?;
func DeleteAllItemSizes(itemId int) error {
	if err := db.Model(&ItemSizes{}).Where("item_id = ?", itemId).Updates(map[string]interface{}{"is_del": 1}).Error; err != nil {
		return err
	}
	return nil
}
