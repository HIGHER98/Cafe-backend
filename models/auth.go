package models

import (
	"cafe/pkg/logging"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"`
	IsDel    int    `json:"is_del"`
}

//TODO Return roll as well
func Check(username, password string) (bool, error) {
	var user Users
	err := db.Select("id, password, role").Where(Users{Username: username, IsDel: 0}).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, err
	}

	if user.ID > 0 {
		return true, nil
	}
	return false, nil
}

//Returns false, nil if no user exists
func CheckUserExists(username string) (bool, error) {
	var user Users
	err := db.Select("id").Where(Users{Username: username, IsDel: 0}).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if user.ID > 0 {
		return true, nil
	}
	return false, nil
}

func SignupUser(data map[string]interface{}) (bool, error) {
	var role int
	if roleid, ok := data["role"].(int); ok {
		role = roleid
	} else if roleid, ok := data["role"].(float64); ok {
		role = int(roleid)
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(data["password"].(string)), 8)
	if err != nil {
		return false, err
	}

	user := Users{
		Username: data["username"].(string),
		Password: string(hashedPass),
		Role:     role,
	}
	logging.Info("Signing up new user: ", user)
	if err = db.Create(&user).Error; err != nil {
		return false, err
	}
	return true, nil
}
