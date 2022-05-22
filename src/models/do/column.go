package do

import (
	"SqlLineage/src/configs"
	"gorm.io/gorm"
)

type Column struct {
	Name      string
	Type      string
	DatasetId uint
	configs.Model
}

func (this *Column) TableName() string {
	return "column"
}

func ConstructorColumn(name, _type string, datasetId uint) *Column {
	return &Column{Name: name, Type: _type, DatasetId: datasetId}
}

func InsertColumn(column *Column, db *gorm.DB) error {
	if db == nil {
		db = configs.GormDB
	}
	vertex := ConstructorVertex(column.Name, DatabaseVertex)
	InsertVertex(vertex, db)
	column.VertexId = vertex.Id
	create := db.Create(&column)
	if create.Error != nil || create.RowsAffected != 1 {
		return create.Error
	}
	return nil
}

func SelectColumnsByDatasetID(datasetId uint) []Column {
	var cols []Column
	result := configs.GormDB.Where(&Column{DatasetId: datasetId}, "DatasetId").Find(&cols)
	if result.Error != nil {
		return nil
	} else {
		return cols
	}
}

func SelectDownStreamColumnWithCurVertexId(vertexId uint64, _type EdgeType) []Column {
	var columns []Column
	configs.GormDB.Raw("SELECT `column`.* FROM `column` LEFT JOIN edge ON `column`.vertex_id = edge.end_vertex WHERE edge.start_vertex = ? AND edge.type = ?", vertexId, _type).Scan(&columns)
	return columns
}

func SelectUpStreamColumnWithCurVertexId(vertexId uint64, _type EdgeType) []Column {
	var columns []Column
	configs.GormDB.Raw("SELECT `column`.* FROM `column` LEFT JOIN edge ON `column`.vertex_id = edge.start_vertex WHERE edge.end_vertex = ? AND edge.type = ?", vertexId, _type).Scan(&columns)
	return columns
}

func SelectColumnWithVertexId(vertexId uint64) Column {
	var column Column
	column.VertexId = vertexId
	configs.GormDB.Select(column, "VertexId").Find(&column)
	return column
}
