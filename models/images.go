package models

import "github.com/jinzhu/gorm"

type Image struct {
	Id     int    `json:"id"`
	ItemId int    `json:"item_id"`
	S3Link string `json:"s3_link"`
	IsDel  string `json:"is_del"`
}

//INSERT INTO TBALE image (item_id, s3_link) VALUES (?, ?);
func AddImage(itemId int, s3link string) error {
	img := Image{ItemId: itemId, S3Link: s3link}
	if err := db.Create(&img).Error; err != nil {
		return err
	}
	return nil
}

//SELECT * FROM Image WHERE id=? LIMIT 1;
func GetImage(id int) (*Image, error) {
	var img Image
	if err := db.Where("id = ?", id).First(&img).Error; err != nil {
		return &img, err
	}
	return &img, nil
}

//SELECT * FROM image WHERE item_id=?;
func GetImagesForItem(itemId int) ([]Image, error) {
	var images []Image
	err := db.Where("item_id = ?", itemId).Find(&images).Error
	return images, err
}

//UPDATE TABLE IMAGE SET is_del=1 WHERE id=?;
func DelImage(id int) error {
	img, err := GetImage(id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err = db.Model(&img).Update("is_del", 1).Error; err != nil {
		return err
	}
	return nil
}
