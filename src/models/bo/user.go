package bo

import "SqlLineage/src/models/do"

type UserInfo struct {
	Username string
	Password string
}

type AuthBo struct {
	Id       uint
	Username string
	NickName string
	Gender   int
	Mail     string
}

func UserDo2AuthBo(user *do.User) AuthBo {
	return AuthBo{
		Id:       user.ID,
		Username: user.Username,
		NickName: user.NickName,
		Gender:   user.Gender,
		Mail:     user.Mail,
	}
}
