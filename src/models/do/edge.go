package do

import (
	"SqlLineage/src/configs"
	"gorm.io/gorm"
)

type Edge struct {
	Id          uint64
	StartVertex uint64
	EndVertex   uint64
	Type        EdgeType
}

type EdgeType uint8

// graph 边类型
const (
	Table2table   EdgeType = 1
	Table2column  EdgeType = 2
	Column2column EdgeType = 3
)

func (this *Edge) TableName() string {
	return "edge"
}

func ConstructorEdge(startVertex, endVertex uint64, _type EdgeType) *Edge {
	return &Edge{
		StartVertex: startVertex,
		EndVertex:   endVertex,
		Type:        _type,
	}
}

func InsertEdge(edge *Edge, db *gorm.DB) error {
	if db == nil {
		db = configs.GormDB
	}
	create := db.Create(edge)
	if create.Error != nil || create.RowsAffected != 1 {
		return create.Error
	}
	return nil
}
