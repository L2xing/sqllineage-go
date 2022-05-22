package do

import (
	"SqlLineage/src/configs"
	"gorm.io/gorm"
)

type Vertex struct {
	Id   uint64
	Name string
	Type VertexType
}

func (this *Vertex) TableName() string {
	return "vertex"
}

type VertexType uint8

// graph 节点类型
const (
	DatasetVertex  VertexType = 1
	DatabaseVertex VertexType = 2
	ColumnVertex   VertexType = 3
	InstanceVertex VertexType = 4
	VDatasetVertex VertexType = 5
)

func ConstructorVertex(vertexName string, _type VertexType) *Vertex {
	return &Vertex{
		Name: vertexName,
		Type: _type,
	}
}

func InsertVertex(vertex *Vertex, db *gorm.DB) error {
	if db == nil {
		db = configs.GormDB
	}
	create := db.Create(vertex)
	if create.Error != nil || create.RowsAffected != 1 {
		return create.Error
	}
	return nil
}

func SelectVertexWithId(id uint64) Vertex {
	var vertex Vertex
	vertex.Id = id
	configs.GormDB.Find(&vertex)
	return vertex
}

func SelectEndVertexWithStartVertex(startVertexId uint64, _type EdgeType) []Vertex {
	//SELECT dataset.* FROM dataset LEFT JOIN edge ON dataset.vertex_id = edge.end_vertex WHERE edge.start_vertex = 00003811 AND edge.type = 1
	var vertexs []Vertex
	configs.GormDB.Model(&Vertex{}).Joins("")
	return vertexs
}

func SelectStartVertexWithEndVertex(endVertexId uint64, _type EdgeType) []Vertex {
	var vertexs []Vertex
	configs.GormDB.Where(&Edge{
		EndVertex: endVertexId,
	}).Find(&vertexs)
	return vertexs
}
