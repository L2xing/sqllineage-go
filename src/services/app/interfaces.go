package app

import "SqlLineage/src/models/bo"

type UserService interface {
	verify(username, password string) bo.AuthBo
	register(username, password string) bo.AuthBo
}

type DatabaseService interface {
}

type DatasetService interface {
}

type InstanceService interface {
}

type ColumnService interface {
}
