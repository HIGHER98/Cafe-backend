package models

type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetTagById(id int) (*Tag, error) {
	var tag *Tag
	if err := db.Select(&tag).Where("id=?", id).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

func GetAllTags() ([]*Tag, error) {
	var tags []*Tag
	if err := db.Where("is_del = 0").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func AddTag(name string) error {
	tag := Tag{Name: name}
	if err := db.Create(&tag).Error; err != nil {
		return err
	}
	return nil
}

func EditTag(id int, name string) error {
	if err := db.Model(&Tag{}).Where("id=?", id).Update("name", name).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTag(id int) error {
	if err := db.Model(&Tag{}).Where("id=?", id).Update("is_del", 1).Error; err != nil {
		return err
	}
	return nil
}
