package models

import (
	"cafe/pkg/logging"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Users struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"`
	IsDel    int    `json:"is_del"`
}

//SELECT id, password, role FROM users WHERE username=? AND is_del=0 LIMIT 1;
func Check(username, password string) (bool, error) {
	var user Users
	if err := db.Select("id, password, role").Where(Users{Username: username, IsDel: 0}).First(&user).Error; err != nil {
		return false, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, err
	}

	if user.ID > 0 {
		return true, nil
	}
	return false, nil
}

//Returns false, nil if no user exists
//If user exists is marked as is_del=1, they are deleted from the database and false, nil is returned
func CheckUserExists(username string) (bool, error) {
	var user Users
	err := db.Select("id").Where(Users{Username: username}).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if user.ID > 0 && user.IsDel == 0 {
		return true, nil
	} else if user.IsDel == 1 {
		//Free up username if user has been marked deleted in the DB
		if err = DeleteUser(user.ID); err != nil {
			return true, err
		}
	}
	return false, nil
}

//Same as the CheckUserExists(username string) but doesn't delete a user marked is_del=1
func CheckUserExistsNoDel(id int) (bool, error) {
	var user Users
	if err := db.Select("id").Where(Users{ID: id}).First(&user).Error; err != nil {
		return false, err
	}
	return true, nil
}

//Update user role
func UpdateUser(id, role int) error {
	exists, err := CheckUserExistsNoDel(id)
	if err != nil || !exists {
		return err
	}
	exists, err = CheckRoleExists(role)
	if err != nil || !exists {
		return err
	}
	if err := db.Model(&Users{}).Where("id=?", id).Update("role", role).Error; err != nil {
		return err
	}
	return nil
}

//UPDATE users SET is_del=1 WHERE id=?;
func MarkUserDeleted(id int) error {
	if err := db.Model(&Users{}).Where("id=?", id).Update("is_del", 1).Error; err != nil {
		return err
	}
	return nil
}

//DELETE from users WHERE id=?;
func DeleteUser(id int) error {
	if err := db.Where("id=?", id).Delete(&Users{}).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil
		default:
			return err
		}
	}
	return nil
}

//INSERT INTO users (username, password, role) VALUES (?, ?, ?)
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

type UserView struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Title    string `json:"title"`
}

func GetActiveUsers() ([]*UserView, error) {
	var users []*UserView
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

type Staff struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetStaff() ([]*Staff, error) {
	var staff []*Staff
	if err := db.Where("is_del=0").Find(&staff).Error; err != nil {
		return nil, err
	}
	return staff, nil
}
