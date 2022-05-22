package interfaces

import (
	"SqlLineage/src/models/do"
)

type IUseApi interface {
	insertUser(user do.User)
	updateUser(user do.User)
	deleteUser(id int64)
	getUser(id int64)
}
