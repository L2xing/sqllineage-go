package test

import (
	"SqlLineage/src/configs"
	"SqlLineage/src/models/do"
	"SqlLineage/src/services/app"
	"SqlLineage/src/utils/net"
	"fmt"
	"testing"
)

func SetUp() {
	// 连接数据库
	configs.SetUpGrom()
	// 初始化Casbin
}

func TestCreateUser(t *testing.T) {

}

func TestDeleteUser(t *testing.T) {
	SetUp()
}

func TestSelectUserById(t *testing.T) {
	SetUp()
	user, _ := do.SelectUserByUsernamePassword("user1", "ps")
	fmt.Println(user)
}

func TestSelectDatasetById(t *testing.T) {
	SetUp()
	databases := app.GetDatabaseByInstanceId(5)
	datasets := app.GetDatasetsByDatabaseId(1)
	columns := app.GetColumnsByDatasetId(1)
	fmt.Println(databases)
	fmt.Println(datasets)
	fmt.Println(columns)
}

func TestInsetDataset(t *testing.T) {
	SetUp()
	database := do.ConstructorDatabase("myDatabase", 0)
	do.InsertDatabase(database, nil)
	fmt.Println(database)
}

func TestInsetVertex(t *testing.T) {
	SetUp()
	vertex := do.ConstructorVertex("12312", do.DatabaseVertex)
	do.InsertVertex(vertex, configs.GormDB)
	fmt.Println(vertex)
}

func TestDownStreamDataset(t *testing.T) {
	SetUp()
	var id uint64 = 3811
	vertexId := do.SelectDownStreamDatasetWithCurVertexId(id, do.Table2table)
	id = 3821
	fmt.Println(vertexId)
	vertexId = do.SelectUpStreamDatasetWithCurVertexId(id, do.Table2table)
	fmt.Println(vertexId)

}

func TestSelectVertex(t *testing.T) {
	SetUp()
	var id uint64 = 3811
	withId := do.SelectVertexWithId(id)
	fmt.Println(withId)
}

func TestDeleteInstance(t *testing.T) {
	SetUp()
	var id uint = 6
	do.DeleteInstanceByID(id)
}

func TestSelectInstance(t *testing.T) {
	SetUp()
	SetUp()
	var id uint = 6
	byId := do.SelectInstanceById(id)
	fmt.Println(byId)
}

func TestSelectAllInstance(t *testing.T) {
	SetUp()
	pageInfo := net.NewPageInfo(100, 1, "")
	byId := do.SelectInstanceByUsername("admin", *pageInfo)
	fmt.Println(byId)
}
