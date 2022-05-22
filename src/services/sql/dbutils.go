package sql

import (
	"SqlLineage/src/models/do"
	"fmt"
	"github.com/Sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

const (
	pattern = "^db[0-9]+@[[:alpha:]]+$"
)

type DbFetcher interface {
	CanConnect() bool
	FetchAllDataBases() []string
	FetchAllTables(db string) []string
	FetchAllFileds(db, table string) []string
}

type MysqlFetcher struct {
	dsn         string
	db          *gorm.DB
	fetchDBS    string
	fetchTables string
	fetchFileds string
	canConnect  bool
}

func NewMysqlFetcher(ip, port, user, password string) *MysqlFetcher {
	str := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local&allowOldPasswords=1", user, password, ip, port, "information_schema")
	return &MysqlFetcher{
		dsn:         str,
		fetchDBS:    "SELECT `SCHEMA_NAME` AS Name FROM `information_schema`.`SCHEMATA`",
		fetchTables: "SELECT TABLE_NAME AS Name FROM  `information_schema`.`TABLES` WHERE TABLE_SCHEMA = '%s';",
		fetchFileds: "SELECT COLUMN_NAME AS Name FROM `information_schema`.`COLUMNS`  WHERE TABLE_SCHEMA = \"%s\" AND TABLE_NAME = \"%s\";",
	}
}

//TODO 有并发隐患
func (this *MysqlFetcher) CanConnect() bool {
	db, err := gorm.Open(mysql.Open(this.dsn), &gorm.Config{})
	if err != nil {
		this.canConnect = false
	} else {
		this.canConnect = true
		this.db = db
	}
	return this.canConnect
}
func (this *MysqlFetcher) FetchAllDataBases() []string {
	if !this.canConnect {
		if !this.CanConnect() {
			fmt.Errorf("无法连接目标数据库")
			return nil
		}
	}
	var dbs []string
	this.db.Raw(this.fetchDBS).Scan(&dbs)
	return dbs
}

func (this *MysqlFetcher) FetchAllTables(db string) []string {
	if !this.canConnect {
		if !this.CanConnect() {
			fmt.Errorf("无法连接目标数据库")
			return nil
		}
	}
	var dbs []string
	dsn := fmt.Sprintf(this.fetchTables, db)
	this.db.Raw(dsn).Scan(&dbs)
	return dbs
}

func (this *MysqlFetcher) FetchAllFileds(db, table string) []string {
	if !this.canConnect {
		if !this.CanConnect() {
			fmt.Errorf("无法连接目标数据库")
			return nil
		}
	}
	var dbs []string
	dsn := fmt.Sprintf(this.fetchFileds, db, table)
	this.db.Raw(dsn).Scan(&dbs)
	return dbs
}

func ParserDatabaseName(name string) (uint, string, bool) {
	if match, err := regexp.MatchString(pattern, name); err != nil || !match {
		return 0, "", false
	}
	split := strings.Split(name, "@")
	id := 0
	for i := 2; i < len(split[0]); i++ {
		id = id*10 + int(split[0][i]-'0')
	}
	return uint(id), split[1], true
}

func SelectDataset(instanceId uint, dbName, name string) *do.Dataset {
	database := do.SelectDatabaseByInstanceIDAndName(instanceId, dbName)
	if database == nil {
		logrus.Errorf("%d%s%s", instanceId, dbName, name)
		return nil
	}
	dataset := do.SelectDatasetByDatabaseIdAndName(database.ID, name)
	if dataset == nil {
		logrus.Errorf("%d%s%s", instanceId, dbName, name)
		return nil
	}
	return dataset
}
