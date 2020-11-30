package models

type ItemOptions struct {
	ItemId      int     `json:"item_id"`
	Opt         string  `json:"opt"`
	AddPrice    float64 `json:"add_price"`
	Description string  `json:"description"`
}

//INSERT INTO item_options VALUES (item_id, opt, add_price, description) VALUES (?, ?, ?, ?);
func (option *ItemOptions) AddItemOption() error {
	if err := db.Create(option).Error; err != nil {
		return err
	}
	return nil
}

//UPDATE item_options SET opt=?, add_price-?, description=? WHERE id=optId;
func (option *ItemOptions) UpdateItemOption(id int) error {
	if err := db.Model(option).Where("id=?", id).Updates(&option).Error; err != nil {
		return err
	}
	return nil
}

//UPDATE item_option SET is_del=1 WHERE id=optId;
func DeleteItemOption(id int) error {
	var option *ItemOptions
	if err := db.Model(&option).Where("id=?", id).Update("is_del", 1).Error; err != nil {
		return err
	}
	return nil
}
