package models

type roles struct {
	Id int `json:"id"`
}

func CheckRoleExists(id int) (bool, error) {
	var r roles
	if err := db.Select("id").Where("id=?", id).First(&r).Error; err != nil {
		return false, err
	}
	return true, nil
}
