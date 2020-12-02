package models

type ItemOptions struct {
	Id          int      `json:"opt_id,omitempty" gorm:"column:id"`
	ItemId      int      `json:"item_id"`
	Opt         string   `json:"opt"`
	AddPrice    *float64 `json:"option_price" gorm:"column:"add_price" default:0"`
	Description *string  `json:"option_description" gorm:"default:null column:description"`
}

//INSERT INTO item_options VALUES (item_id, opt, add_price, description) VALUES (?, ?, ?, ?);
func (option *ItemOptions) AddItemOption() error {
	if err := db.Create(option).Error; err != nil {
		return err
	}
	return nil
}

//SELECT * FROM item_option WHERE id = ?;
func GetItemOption(id int) (*ItemOptions, error) {
	var option ItemOptions
	if err := db.Where("id=?", id).First(&option).Error; err != nil {
		return nil, err
	}
	return &option, nil
}

//UPDATE item_options SET opt=?, add_price-?, description=? WHERE id=optId;
func (option *ItemOptions) UpdateItemOption(id int) error {
	if err := db.Model(option).Where("id=?", id).Updates(&option).Error; err != nil {
		return err
	}
	return nil
}

//UPDATE item_options SET is_del=1 WHERE id=optId;
func DeleteItemOption(id int) error {
	var option *ItemOptions
	if err := db.Model(&option).Where("id=?", id).Update("is_del", 1).Error; err != nil {
		return err
	}
	return nil
}

//UPDATE item_options SET is_del=1 WHERE item_id=?;
func DeleteAllItemOptions(itemId int) error {
	if err := db.Model(&ItemOptions{}).Where("item_id = ?", itemId).Updates(map[string]interface{}{"is_del": 1}).Error; err != nil {
		return err
	}
	return nil
}
