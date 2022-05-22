package do

import (
	"SqlLineage/src/configs"
	"SqlLineage/src/utils/net"
	"github.com/Sirupsen/logrus"
	"gorm.io/gorm"
)

type Instance struct {
	Host     string
	Port     string
	Username string
	Password string
	Type     string
	Owner    string
	configs.Model
}

const (
	Virtual         = "virtual"
	VirtualDatabase = "virDatabase"
)

func (this *Instance) TableName() string {
	return "instance"
}

func ConstructorInstance(host, port, username, password, _type, owner string) *Instance {
	return &Instance{Host: host, Port: port, Username: username, Password: password, Type: _type, Owner: owner}
}

func InsertInstance(instance *Instance, db *gorm.DB) error {
	if db == nil {
		db = configs.GormDB
	}
	create := db.Create(&instance)
	if create.Error != nil || create.RowsAffected != 1 {
		return create.Error
	}
	return nil
}

func SelectInstanceById(instanceId uint) *Instance {
	var instance Instance
	result := configs.GormDB.First(&instance, instanceId)
	if result.Error != nil {
		return nil
	}
	return &instance
}

func SelectInstanceByUsername(username string, page net.PageInfo) *net.PageResult {
	var instances []Instance
	var total int64
	if result := configs.GormDB.Where(&Instance{
		Owner: username,
	}, "Owner").Offset((page.Current - 1) * page.Size).Limit(page.Size).Find(&instances); result.Error != nil {
		logrus.Errorf("查询%s用户数据失败", username)
		return nil
	}
	if result := configs.GormDB.Model(&Instance{}).Where(&Instance{
		Owner: username,
	}, "Owner").Count(&total); result.Error != nil {
		logrus.Errorf("查询%s用户数据失败", username)
		return nil
	}
	if total%int64(page.Size) == 0 {
		total = total / int64(page.Size)
	} else {
		total = total/int64(page.Size) + 1
	}
	return net.NewPageResult(page.Size, int(total), page.Current, page.Criteria, instances)
}

func DeleteInstanceByID(id uint) error {
	var instance Instance
	result := configs.GormDB.Delete(&instance, id)
	return result.Error
}
