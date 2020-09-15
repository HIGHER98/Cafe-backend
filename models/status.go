package models

import "github.com/jinzhu/gorm"

const (
	PENDING_TRANSACTION = 1
	PENDING             = 2
	CONFIRMED           = 3
	COLLECTION          = 4
)

type Status struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	IsDel       int    `json:"is_del"`
}

//INSERT INTO TABLE status (description) VALUES (?);
func AddStatus(desc string) error {
	newStatus := Status{Description: desc}
	err := db.Create(&newStatus).Error
	return err
}

//Get Status
func GetStatus(id int) (*Status, error) {
	var status Status
	err := db.Where("id").Where("id=?", id).First(&status).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &status, err
	}
	return &status, nil
}

//UPDATE TABLE status SET is_del=0 WHERE id=?;
func DeleteStatus(id int) error {
	status, err := GetStatus(id)
	if err != nil {
		return err
	}
	db.Model(&status).Update("is_del", 1)
	return err
}
