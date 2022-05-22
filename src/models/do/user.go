package do

import (
	"SqlLineage/src/configs"
	"gorm.io/gorm"
)

type User struct {
	Username string
	Password string
	NickName string
	Gender   int
	Mail     string
	gorm.Model
}

func (this *User) TableName() string {
	return "user"
}

func ConstructorUser(username, password string) *User {
	return &User{Username: username, Password: password}
}

func InsertUser(user *User) error {
	create := configs.GormDB.Omit("nick_name").Create(&user)
	if create.Error != nil || create.RowsAffected != 1 {
		return create.Error
	}
	return nil
}

func SelectUserByID(id uint) *User {
	var user User
	result := configs.GormDB.Find(&user, id)
	if result.RowsAffected == 1 {
		return &user
	} else {
		return nil
	}
}

func DeleteUserByID(id uint) error {
	var user User
	result := configs.GormDB.Unscoped().Delete(user, id)
	return result.Error
}

func SelectUserByUsernamePassword(username, password string) (*User, error) {
	var user User
	first := configs.GormDB.Where(&User{
		Username: username,
		Password: password,
	}, "username", "password").First(&user)
	if first.RowsAffected == 1 {
		return &user, nil
	} else {
		return nil, nil
	}
}

//func RecoverUserByID(id uint) error {
//	var user User
//	update := configs.GormDB.Model(user).Where("id = ?", id).Update("is_deleted", false)
//	return update.Error
//}
