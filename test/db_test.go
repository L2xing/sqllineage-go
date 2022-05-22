package test

import (
	"SqlLineage/src/services/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

var DB map[string]*sql.AstTable

//func init() {
//	DB = make(map[string]*sql.AstTable)
//	user := sql.NewAstTableByName("user")
//	teacher := sql.NewAstTableByName("teacher")
//	student := sql.NewAstTableByName("student")
//	DB[user.Name] = user
//	DB[teacher.Name] = teacher
//	DB[student.Name] = student
//
//	userId := sql.NewAstFiledByName("id", user)
//	username := sql.NewAstFiledByName("username", user)
//	password := sql.NewAstFiledByName("password", user)
//	user.AddFiled(userId)
//	user.AddFiled(username)
//	user.AddFiled(password)
//
//	teacherId := sql.NewAstFiledByName("id", teacher)
//	teacherName := sql.NewAstFiledByName("name", teacher)
//	teacherAge := sql.NewAstFiledByName("age", teacher)
//	teacher.AddFiled(teacherId)
//	teacher.AddFiled(teacherName)
//	teacher.AddFiled(teacherAge)
//
//	studentId := sql.NewAstFiledByName("id", student)
//	studentName := sql.NewAstFiledByName("name", student)
//	studentClass := sql.NewAstFiledByName("class", student)
//	student.AddFiled(studentId)
//	student.AddFiled(studentName)
//	student.AddFiled(studentClass)
//}

func TestGetStudent(t *testing.T) {
	if _, ok := DB["student"]; ok {
		println("存在student表")
	} else {
		println("没有student表")
	}
}

type DataBases struct {
	Name string
}

func TestShowDBS(t *testing.T) {
	str := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local&allowOldPasswords=1", "root", "123456", "127.0.0.1", "3306", "information_schema")
	println(str)
	db, _ := gorm.Open(mysql.Open(str), &gorm.Config{})
	var DBS []string
	db.Raw("SELECT `SCHEMA_NAME` AS Name FROM `information_schema`.`SCHEMATA`").Scan(&DBS)
	for _, v := range DBS {
		fmt.Print(v)
		var targets []string
		dsn := fmt.Sprintf("SELECT TABLE_NAME AS Name FROM  `information_schema`.`TABLES` WHERE TABLE_SCHEMA = '%s';", v)
		db.Raw(dsn).Scan(&targets)
		fmt.Print("————>")
		for _, v2 := range targets {
			var targets []string

			dsn = fmt.Sprintf("SELECT COLUMN_NAME AS Name FROM `information_schema`.`COLUMNS`  WHERE TABLE_SCHEMA = \"%s\" AND TABLE_NAME = \"%s\"", v, v2)
			db.Raw(dsn).Scan(&targets)
			for _, v3 := range targets {
				fmt.Print(v3)
				fmt.Print(", ")
			}
		}
		fmt.Println()
	}
}

func TestDBUtils(t *testing.T) {
	fetcher := sql.NewMysqlFetcher("127.0.0.1", "3306", "root", "123456")
	for _, db := range fetcher.FetchAllDataBases() {

		fmt.Println(db)
		for _, table := range fetcher.FetchAllTables(db) {
			fmt.Print("————+")
			fmt.Println(table)
			fmt.Print("————+————+")
			for _, filed := range fetcher.FetchAllFileds(db, table) {
				fmt.Print(filed)
				fmt.Print(", ")
			}
			fmt.Println()
		}
	}
}

// instance@database.table
