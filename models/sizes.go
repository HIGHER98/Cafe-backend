package models

type ItemSizes struct {
	ItemId      int     `json:"item_id"`
	ItemSize    string  `json:"item_size"`
	AddPrice    float64 `json:"add_price,string,omitempty"`
	Description *string `json:"description" gorm:"default:null"`
}

//INSERT INTO item_sizes VALUES (item_id, item_size, add_price, description) VALUES (?, ?, ?, ?);
func (size *ItemSizes) AddItemSize() error {
	if err := db.Create(size).Error; err != nil {
		return err
	}
	return nil
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
