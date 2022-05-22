package do

import (
	"SqlLineage/src/configs"
	"gorm.io/gorm"
)

type Database struct {
	Name       string
	InstanceId uint
	configs.Model
}

func (this *Database) TableName() string {
	return "database"
}

func ConstructorDatabase(name string, InstanceId uint) *Database {
	return &Database{Name: name, InstanceId: InstanceId}
}

func InsertDatabase(database *Database, db *gorm.DB) error {
	if db == nil {
		db = configs.GormDB
	}
	vertex := ConstructorVertex(database.Name, DatabaseVertex)
	InsertVertex(vertex, db)
	database.VertexId = vertex.Id
	create := db.Create(&database)
	if create.Error != nil || create.RowsAffected != 1 {
		return create.Error
	}
	return nil
}

func SelectDatabasesByInstanceID(instanceId uint) []Database {
	var databases []Database
	result := configs.GormDB.Where(&Database{InstanceId: instanceId}, "InstanceId").Find(&databases)
	if result.Error != nil {
		return nil
	}
	return databases
}

func SelectDatabaseByInstanceIDAndName(instanceId uint, dbName string) *Database {
	var database Database
	result := configs.GormDB.Where(&Database{InstanceId: instanceId, Name: dbName}, "InstanceId", "Name").First(&database)
	if result.Error != nil {
		return nil
	}
	return &database
}

//func DeleteDatabaseByID(id uint) error {
//	result := configs.GormDB.Unscoped().Delete(&this, id)
//	return result.Error
//}
