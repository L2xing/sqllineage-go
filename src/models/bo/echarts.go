package bo

const (
	InstanceNode GraphNodeType = 0
	DatabaseNode GraphNodeType = 1
	DatasetNode  GraphNodeType = 2
	VDatasetNode GraphNodeType = 3
	ColumnNode   GraphNodeType = 4
)

var GraphCategory []string = []string{"Instance", "Database", "Dataset", "VDataset", "Column"}

type GraphNodeType int

type GraphNode struct {
	// 节点id
	Id string `json:"id"`
	// 节点名称
	Name string `json:"name"`
	// 权值
	Value uint `json:"value"`
	// echarts节点显示分类
	Category int `json:"category"`
	// echarts权重
	SymbolSize int `json:"symbolSize"`
	// 节点类型
	Type string `json:"type"`
}

type GraphEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type GraphModel struct {
	Nodes []GraphNode
	Edges []GraphEdge
}

func NewGraphEdge(source, target string) GraphEdge {
	return GraphEdge{source, target}
}

func NewGraphNode(id, name string, value uint, _type GraphNodeType) GraphNode {
	node := GraphNode{Id: id, Name: name, Value: value, Type: GraphCategory[_type]}
	switch _type {
	case InstanceNode:
		node.SymbolSize = 40
		node.Category = 0
	case DatabaseNode:
		node.SymbolSize = 35
		node.Category = 1
	case DatasetNode:
		node.SymbolSize = 30
		node.Category = 2
	case VDatasetNode:
		node.SymbolSize = 20
		node.Category = 3
	case ColumnNode:
		node.SymbolSize = 15
		node.Category = 4
	}
	return node
}
