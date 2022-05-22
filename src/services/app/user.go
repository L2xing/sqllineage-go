package app

import (
	"SqlLineage/src/models/bo"
	"SqlLineage/src/models/do"
	"github.com/Sirupsen/logrus"
)

// User登录
func Verify(username, password string) *bo.AuthBo {
	user, _ := do.SelectUserByUsernamePassword(username, password)
	if user == nil {
		return nil
	}
	authBo := bo.UserDo2AuthBo(user)
	return &authBo
}

// User注册
func Register(username, password string) *bo.AuthBo {
	user, _ := do.SelectUserByUsernamePassword(username, password)
	if user != nil {
		logrus.Errorf("user已存在；user:%s", username)
		return nil
	}
	puser := &do.User{Username: username, Password: password}
	if err := do.InsertUser(puser); err != nil {
		logrus.Errorf("user添加失败;user:%s", username)
		return nil
	}
	authBo := bo.UserDo2AuthBo(puser)
	return &authBo
}
