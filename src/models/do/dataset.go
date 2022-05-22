package do

import (
	"SqlLineage/src/configs"
	"gorm.io/gorm"
)

type Dataset struct {
	Name       string
	DatabaseId uint
	configs.Model
}

func (this *Dataset) TableName() string {
	return "dataset"
}

func ConstructorDataset(name string, databaseId uint) *Dataset {
	return &Dataset{Name: name, DatabaseId: databaseId}
}

func InsertDataset(dataset *Dataset, db *gorm.DB) error {
	if db == nil {
		db = configs.GormDB
	}
	vertex := ConstructorVertex(dataset.Name, DatasetVertex)
	InsertVertex(vertex, db)
	dataset.VertexId = vertex.Id
	create := db.Create(&dataset)
	if create.Error != nil || create.RowsAffected != 1 {
		return create.Error
	}
	return nil
}

func InsertVDataset(dataset *Dataset, db *gorm.DB) error {
	if db == nil {
		db = configs.GormDB
	}
	vertex := ConstructorVertex(dataset.Name, VDatasetVertex)
	InsertVertex(vertex, db)
	dataset.VertexId = vertex.Id
	create := db.Create(&dataset)
	if create.Error != nil || create.RowsAffected != 1 {
		return create.Error
	}
	return nil
}

func SelectDatasetsByDatabaseId(databaseId uint) []Dataset {
	var datasets []Dataset
	result := configs.GormDB.Where(&Dataset{DatabaseId: databaseId}, "DatabaseId").Find(&datasets)
	if result.Error != nil {
		return nil
	} else {
		return datasets
	}
}

func SelectDatasetByVertexId(id uint64) Dataset {
	var dataset Dataset
	dataset.VertexId = id
	configs.GormDB.Where(&dataset, "VertexId").Find(&dataset)
	return dataset
}

func SelectDatasetByDatasetId(id uint) Dataset {
	var dataset Dataset
	dataset.ID = id
	configs.GormDB.Where(&dataset, "ID").Find(&dataset)
	return dataset
}

func SelectDatasetByDatabaseIdAndName(databaseId uint, name string) *Dataset {
	var dataset Dataset
	result := configs.GormDB.Where(&Dataset{DatabaseId: databaseId, Name: name}, "DatabaseId", "Name").Find(&dataset)
	if result.Error != nil {
		return nil
	} else {
		return &dataset
	}
}

func SelectDownStreamDatasetWithCurVertexId(vertexId uint64, _type EdgeType) []Dataset {
	var datasets []Dataset
	configs.GormDB.Raw("SELECT dataset.* FROM dataset LEFT JOIN edge ON dataset.vertex_id = edge.end_vertex WHERE edge.start_vertex = ? AND edge.type = ?", vertexId, _type).Scan(&datasets)
	return datasets
}

func SelectUpStreamDatasetWithCurVertexId(vertexId uint64, _type EdgeType) []Dataset {
	var datasets []Dataset
	configs.GormDB.Raw("SELECT dataset.* FROM dataset LEFT JOIN edge ON dataset.vertex_id = edge.start_vertex WHERE edge.end_vertex = ? AND edge.type = ?", vertexId, _type).Scan(&datasets)
	return datasets
}
