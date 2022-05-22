package bo

import "SqlLineage/src/models/do"

type Sql struct {
	Sql string `json:"sql" form:"sql"`
}

type AddSqlBo struct {
	Uuid       string `json:"uuid" form:"uuid"`
	TableName  string `json:"tableName" form:"tableName"`
	InstanceId uint64 `json:"instanceId" form:"instanceId"`
}

type DirType uint8

const (
	IN   DirType = 0
	OUT  DirType = 1
	BOTH DirType = 2
)

type SqlNodeQuery struct {
	NodeId   uint64        `json:"uuid" form:"uuid"`
	NodeType do.VertexType `json:"nodeType" form:"nodeType"`
	DirType  DirType       `json:"dirType" form:"dirType"`
	EdgeType do.EdgeType   `json:"edgeType" form:"edgeType"`
}
