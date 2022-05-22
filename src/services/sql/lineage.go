package sql

import (
	"SqlLineage/src/configs"
	"SqlLineage/src/models/bo"
	"SqlLineage/src/models/do"
	"github.com/Sirupsen/logrus"
	"github.com/gofrs/uuid"
	"sync"
)

var UserLineageMaps map[string]*LineageContext = make(map[string]*LineageContext, 0)
var mutex sync.RWMutex

func ParserSqlLineage(sql string) (*LineageContext, string) {
	v1, _ := uuid.NewV1()
	uuid := v1.String()
	if tparser, err := ASTparser(sql); err != nil {
		return nil, ""
	} else {
		mutex.Lock()
		// TODO 到时候使用redis进行替换
		UserLineageMaps[uuid] = tparser
		mutex.Unlock()
		return tparser, uuid
	}
}

func AddSqlLineage(sqlbo *bo.AddSqlBo) bool {
	var lineage *LineageContext
	if context, has := UserLineageMaps[sqlbo.Uuid]; has {
		lineage = context
		mutex.Lock()
		delete(UserLineageMaps, sqlbo.Uuid)
		mutex.Unlock()
	} else {
		return false
	}
	allOk := true
	tx := configs.GormDB.Begin()
	// 获取虚拟实例的虚拟数据库
	database := do.SelectDatabaseByInstanceIDAndName(uint(sqlbo.InstanceId), do.VirtualDatabase)
	if database == nil {
		allOk = false
		logrus.Errorln("不存在该实例或是该实例没有虚拟库，Id:%d", sqlbo.InstanceId)
		return false
	} else {
		// 添加VDataset
		dataset := do.ConstructorDataset(sqlbo.TableName, database.ID)
		if err := do.InsertVDataset(dataset, tx); err != nil {
			allOk = false
		} else {
			// 添加DataSet血缘
			_, tables := GetAstTables(lineage)
			for _, v := range *tables {
				if !allOk {
					break
				}
				edge := do.ConstructorEdge(v.VertexId, dataset.VertexId, do.Table2table)
				if err := do.InsertEdge(edge, tx); err != nil {
					allOk = false
				}
			}
		}
		// 添加Column
		_, columns := GetAstColumns(lineage)
		for k, v := range columns {
			if !allOk {
				break
			}
			column := do.ConstructorColumn(k, "oo", dataset.ID)
			if err := do.InsertColumn(column, tx); err != nil {
				allOk = false
			} else {
				edge := do.ConstructorEdge(v.VertexId, column.VertexId, do.Column2column)
				if err := do.InsertEdge(edge, tx); err != nil {
					allOk = false
				}
			}
		}
	}
	if allOk {
		tx.Commit()
	} else {
		tx.Rollback()
	}
	return allOk
}

func BuildColsGraph(cxt *LineageContext) []bo.GraphNode {
	var cols []bo.GraphNode
	columns, _ := GetAstColumns(cxt)
	for _, column := range columns {
		node := bo.NewGraphNode(column, column, 0, bo.ColumnNode)
		cols = append(cols, node)
	}
	return cols
}

func BuildTablesGraph(cxt *LineageContext) []bo.GraphNode {
	var tables []bo.GraphNode
	astTables, _ := GetAstTables(cxt)
	for _, table := range astTables {
		node := bo.NewGraphNode(table, table, 0, bo.DatasetNode)
		tables = append(tables, node)
	}
	return tables
}
