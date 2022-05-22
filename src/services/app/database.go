package app

import "SqlLineage/src/models/do"

func QueryDatabaseByInstanceId(instanceId uint) []do.Database {
	// TODO 判断权限
	return do.SelectDatabasesByInstanceID(instanceId)
}
