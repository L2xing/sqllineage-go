package configs

import (
	"SqlLineage/src/configs/properties"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/soft_delete"
	"time"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

var GormDB *gorm.DB

func SetUpGrom() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", properties.GolbalMySQL.Username, properties.GolbalMySQL.Password, properties.GolbalMySQL.Host, properties.GolbalMySQL.Port, properties.GolbalMySQL.Db)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	GormDB = db
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, _ := GormDB.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
}

type Model struct {
	gorm.Model
	VertexId  uint64
	IsDeleted soft_delete.DeletedAt `gorm:"softDelete:flag"`
}

//func (this *Model) SelectBatch(id uint) []interface{} {
//	result := GormDB.Find(this, id)
//	if result.RowsAffected == 1 {
//		return this
//	} else {
//		return nil
//	}
//}
//
//func DeleteUserByID(id uint) error {
//	var user User
//	result := configs.GormDB.Unscoped().Delete(&user, id)
//	return result.Error
//}
//
//func RecoverUserByID(id uint) error {
//	var user User
//	update := configs.GormDB.Model(user).Where("id = ?", id).Update("is_deleted", false)
//	return update.Error
//}
