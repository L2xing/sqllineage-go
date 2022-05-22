package app

import (
	"SqlLineage/src/configs"
	"SqlLineage/src/models/do"
	"SqlLineage/src/services/sql"
	"SqlLineage/src/utils/net"
	"github.com/Sirupsen/logrus"
)

func AddInstance(instance *do.Instance) bool {
	tx := configs.GormDB.Begin()
	allOk := true
	if err := do.InsertInstance(instance, tx); err != nil {
		allOk = false
	}
	if instance.Type == do.Virtual {
		database := do.ConstructorDatabase(do.VirtualDatabase, instance.ID)
		if err := do.InsertDatabase(database, tx); err != nil {
			allOk = false
		}
	}

	if allOk {
		tx.Commit()
		return true
	} else {
		tx.Rollback()
		logrus.Errorln("无法添加该数据库实例")
		return false
	}
}

func QueryInstancesPage(info net.PageInfo, username string) net.PageResult {
	byUsername := do.SelectInstanceByUsername(username, info)
	return *byUsername
}

func PullInstanceData(instanceId uint) bool {
	instance := do.SelectInstanceById(instanceId)
	if instance == nil {
		return false
	}
	if instance.Type == do.Virtual {
		logrus.Errorf("%s是虚拟表不可同步", instance.ID)
		return false
	}
	fetcher := sql.NewMysqlFetcher(instance.Host, instance.Port, instance.Username, instance.Password)
	if !fetcher.CanConnect() {
		return false
	}
	tx := configs.GormDB.Begin()
	isOk := true
	// 1. 同步instance下的databases
	databaseNames := fetcher.FetchAllDataBases()
a:
	for _, databaseName := range databaseNames {
		database := do.ConstructorDatabase(databaseName, instanceId)
		if err := do.InsertDatabase(database, tx); err != nil {
			isOk = false
			break a
		}
		// 2. 同步database下的datasets
		tableNames := fetcher.FetchAllTables(databaseName)
		for _, tableName := range tableNames {
			dataset := do.ConstructorDataset(tableName, database.ID)
			if err := do.InsertDataset(dataset, tx); err != nil {
				isOk = false
				break a
			}
			// 3. 同步dataset下的columns
			columns := fetcher.FetchAllFileds(databaseName, tableName)
			for _, column := range columns {
				column := do.ConstructorColumn(column, "oo", dataset.ID)
				if err := do.InsertColumn(column, tx); err != nil {
					isOk = false
					break a
				}
			}
		}
	}
	if isOk {
		tx.Commit()
	} else {
		tx.Rollback()
	}
	return isOk
}

func DeleteInstance(id uint, name string) {
	byId := do.SelectInstanceById(id)
	if byId == nil || byId.Owner != name {
		panic(net.ResponseData{400, "该用户无法删除数据源", nil})
	}
	do.DeleteInstanceByID(id)
}
