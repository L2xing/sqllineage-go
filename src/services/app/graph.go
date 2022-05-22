package app

import (
	"SqlLineage/src/models/bo"
	"SqlLineage/src/models/do"
	"SqlLineage/src/utils"
	"SqlLineage/src/utils/net"
	"fmt"
	"go/types"
)

func GetInstanceByUsername(username string) []bo.GraphNode {
	graph := make([]bo.GraphNode, 0)
	info := net.NewPageInfo(10, 1, "")
	result := do.SelectInstanceByUsername(username, *info)
	if result.Records == nil {
		return nil
	}
	instances := result.Records.([]do.Instance)
	for _, instance := range instances {
		id := fmt.Sprintf("instance_%v", instance.ID)
		graph = append(graph, bo.NewGraphNode(id, instance.Host+":"+instance.Port, instance.ID, bo.InstanceNode))
	}
	return graph
}

func GetDatabaseByInstanceId(id uint) []bo.GraphNode {
	// 获取 instance下的Database
	graph := make([]bo.GraphNode, 0)
	databases := do.SelectDatabasesByInstanceID(id)
	for _, database := range databases {
		id := fmt.Sprintf("database_%v", database.ID)
		graph = append(graph, bo.NewGraphNode(id, database.Name, database.ID, bo.DatabaseNode))
	}
	return graph
}

func GetDatasetsByDatabaseId(databaseId uint) []bo.GraphNode {
	graph := make([]bo.GraphNode, 0)
	// 获取database下的dataset
	datasets := do.SelectDatasetsByDatabaseId(databaseId)
	for _, dataset := range datasets {
		id := fmt.Sprintf("dataset_%v", dataset.ID)
		graph = append(graph, bo.NewGraphNode(id, dataset.Name, dataset.ID, bo.DatasetNode))
	}
	return graph
}

func GetColumnsByDatasetId(datasetId uint) []bo.GraphNode {
	graph := make([]bo.GraphNode, 0)
	colums := do.SelectColumnsByDatasetID(datasetId)
	// 获取 dataset下的column
	for _, column := range colums {
		id := fmt.Sprintf("column_%v", column.ID)
		graph = append(graph, bo.NewGraphNode(id, column.Name, column.ID, bo.ColumnNode))
	}
	return graph
}

func GetGraphNode(query bo.SqlNodeQuery) bo.GraphModel {
	nodes := make([]bo.GraphNode, 0)
	edges := make([]bo.GraphEdge, 0)
	vertex := do.SelectVertexWithId(query.NodeId)
	switch query.NodeType {
	case do.DatasetVertex:
		source := do.SelectDatasetByVertexId(vertex.Id)
		if vertex.Type == do.VDatasetVertex {
			nodes = append(nodes, bo.NewGraphNode(utils.Uint642str(source.VertexId), source.Name, uint(source.VertexId), bo.VDatasetNode))
		} else if vertex.Type == do.DatabaseVertex {
			nodes = append(nodes, bo.NewGraphNode(utils.Uint642str(source.VertexId), source.Name, uint(source.VertexId), bo.DatasetNode))
		}
		if query.DirType == bo.BOTH || query.DirType == bo.IN {
			datases := do.SelectUpStreamDatasetWithCurVertexId(query.NodeId, do.Table2table)
			for _, dataset := range datases {
				vertex := do.SelectVertexWithId(dataset.VertexId)
				if vertex.Type == do.VDatasetVertex {
					nodes = append(nodes, bo.NewGraphNode(utils.Uint642str(dataset.VertexId), dataset.Name, uint(dataset.VertexId), bo.VDatasetNode))
				} else if vertex.Type == do.DatasetVertex {
					nodes = append(nodes, bo.NewGraphNode(utils.Uint642str(dataset.VertexId), dataset.Name, uint(dataset.VertexId), bo.DatasetNode))
				}
				edges = append(edges, bo.NewGraphEdge(utils.Uint642str(dataset.VertexId), utils.Uint642str(query.NodeId)))
			}
		}
		if query.DirType == bo.BOTH || query.DirType == bo.OUT {
			datases := do.SelectDownStreamDatasetWithCurVertexId(query.NodeId, do.Table2table)
			for _, dataset := range datases {
				vertex := do.SelectVertexWithId(dataset.VertexId)
				if vertex.Type == do.VDatasetVertex {
					nodes = append(nodes, bo.NewGraphNode(utils.Uint642str(dataset.VertexId), dataset.Name, uint(dataset.VertexId), bo.VDatasetNode))
				} else if vertex.Type == do.DatasetVertex {
					nodes = append(nodes, bo.NewGraphNode(utils.Uint642str(dataset.VertexId), dataset.Name, uint(dataset.VertexId), bo.DatasetNode))
				}
				edges = append(edges, bo.NewGraphEdge(utils.Uint642str(dataset.VertexId), utils.Uint642str(query.NodeId)))
			}
		}
	case do.ColumnVertex:
		source := do.SelectColumnWithVertexId(vertex.Id)
		nodes = append(nodes, bo.NewGraphNode(utils.Uint642str(source.VertexId), source.Name, uint(source.VertexId), bo.ColumnNode))
		if query.DirType == bo.BOTH || query.DirType == bo.IN {
			columns := do.SelectUpStreamColumnWithCurVertexId(query.NodeId, do.Column2column)
			for _, column := range columns {
				nodes = append(nodes, bo.NewGraphNode(utils.Uint642str(column.VertexId), column.Name, uint(column.VertexId), bo.ColumnNode))
				edges = append(edges, bo.NewGraphEdge(utils.Uint642str(column.VertexId), utils.Uint642str(query.NodeId)))
			}
		}
		if query.DirType == bo.BOTH || query.DirType == bo.OUT {
			columns := do.SelectDownStreamColumnWithCurVertexId(query.NodeId, do.Column2column)
			for _, column := range columns {
				nodes = append(nodes, bo.NewGraphNode(utils.Uint642str(column.VertexId), column.Name, uint(column.VertexId), bo.ColumnNode))
				edges = append(edges, bo.NewGraphEdge(utils.Uint642str(query.NodeId), utils.Uint642str(column.VertexId)))
			}
		}
	}
	return bo.GraphModel{nodes, edges}
}

func GetLineageNode(id uint64) bo.TableLineageBo {
	lineageBo := bo.NewTableLineageBo()
	otable := do.SelectDatasetByVertexId(id)
	table := fillTableFields(&otable)
	lineageBo.AddTableData(&table)
	upTables := make(map[uint]types.Nil, 0)
	downTables := make(map[uint]types.Nil, 0)
	for _, field := range table.Fields {
		// 1. 寻找上游
		upstreams := do.SelectUpStreamColumnWithCurVertexId(field.Id, do.Column2column)
		for _, upstream := range upstreams {
			upTables[upstream.DatasetId] = types.Nil{}
			edge := bo.NewFieldEdge(upstream.VertexId, field.Id)
			lineageBo.AddFieldEdge(&edge)
		}
		// 2. 寻找下游
		downstreams := do.SelectDownStreamColumnWithCurVertexId(field.Id, do.Column2column)
		for _, downstream := range downstreams {
			downTables[downstream.DatasetId] = types.Nil{}
			edge := bo.NewFieldEdge(field.Id, downstream.VertexId)
			lineageBo.AddFieldEdge(&edge)
		}
	}
	for k, _ := range upTables {
		otable := do.SelectDatasetByDatasetId(k)
		table := fillTableFields(&otable)
		lineageBo.AddTableData(&table)
	}
	for k, _ := range downTables {
		otable := do.SelectDatasetByDatasetId(k)
		table := fillTableFields(&otable)
		lineageBo.AddTableData(&table)
	}
	return lineageBo
}

func fillTableFields(otable *do.Dataset) bo.TableData {
	table := bo.NewTableData(otable.VertexId, otable.Name)
	cols := do.SelectColumnsByDatasetID(otable.ID)
	for _, col := range cols {
		data := bo.NewFieldData(col.VertexId, col.Name)
		table.AddFieldData(&data)
	}
	return table
}
