package bo

type TableLineageBo struct {
	Tables      []TableData `json:"tables"`
	FieldsEdges []FieldEdge `json:"fieldsEdges"`
}

func NewTableLineageBo() TableLineageBo {
	return TableLineageBo{make([]TableData, 0), make([]FieldEdge, 0)}
}

func (this *TableLineageBo) AddFieldEdge(fieldEdge *FieldEdge) {
	this.FieldsEdges = append(this.FieldsEdges, *fieldEdge)
}

func (this *TableLineageBo) AddTableData(tableData *TableData) {
	this.Tables = append(this.Tables, *tableData)
}

type TableData struct {
	Id     uint64      `json:"id"`
	Name   string      `json:"name"`
	Fields []FieldData `json:"fields"`
}

func NewTableData(id uint64, name string) TableData {
	return TableData{id, name, make([]FieldData, 0)}
}

func (this *TableData) AddFieldData(field *FieldData) {
	this.Fields = append(this.Fields, *field)
}

type FieldData struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

func NewFieldData(id uint64, name string) FieldData {
	return FieldData{id, name}
}

type FieldEdge struct {
	Origin uint64 `json:"origin"`
	Target uint64 `json:"target"`
}

func NewFieldEdge(origin, target uint64) FieldEdge {
	return FieldEdge{origin, target}
}
