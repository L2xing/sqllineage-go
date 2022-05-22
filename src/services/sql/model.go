package sql

type Obj interface {
	isNull() bool
	isNotNull() bool
}

// 字段节点
type AstFiled struct {
	Name     string
	Alias    string
	Owner    *AstTable
	VertexId uint64
}

func NewAstFiledByName(name string, vertexId uint64, owner *AstTable) *AstFiled {
	return NewAstFiledByNameAlias(name, "", vertexId, owner)
}
func NewAstFiledByNameAlias(name, alias string, vertexId uint64, owner *AstTable) *AstFiled {
	return &AstFiled{Name: name, Alias: alias, Owner: owner, VertexId: vertexId}
}

// 表节点
type AstTable struct {
	Name     string
	Alias    string
	Fileds   []*AstFiled
	Owners   []*AstTable
	VertexId uint64
}

func NewAstTableByName(name string) *AstTable {
	return NewAstTableByNameAlias(name, "")
}
func NewAstTableByNameAlias(name, alias string) *AstTable {
	return &AstTable{Name: name, Alias: alias}
}

func (this *AstTable) AddFiled(filed *AstFiled) {
	this.Fileds = append(this.Fileds, filed)
}

type AstBase struct {
	Name string
}

func NewAstBase(name string) *AstBase {
	return &AstBase{Name: name}
}

// select *** from ***
// 记录每一个Select的items和froms
type LineageContext struct {
	Sql string

	Table *AstTable

	// cols
	Fileds []*AstFiled

	// cols - tables
	CTables map[string]*AstTable

	// froms - tables
	FTables map[string]*AstTable

	// 子SQL
	ChildCxts []*LineageContext
}

func NewLineageContext() *LineageContext {
	cxt := new(LineageContext)
	cxt.CTables = make(map[string]*AstTable, 0)
	cxt.FTables = make(map[string]*AstTable, 0)
	cxt.ChildCxts = make([]*LineageContext, 0)
	cxt.Fileds = make([]*AstFiled, 0)
	return cxt
}

func (this *LineageContext) containsCTable(nameOrAlias string) bool {
	if _, ok := this.CTables[nameOrAlias]; ok {
		return true
	}
	return false
}

func (this *LineageContext) containsFTable(nameOrAlias string) bool {
	if _, ok := this.FTables[nameOrAlias]; ok {
		return true
	}
	return false
}

func (this *LineageContext) GetTable(nameOrAlias string) (bool, *AstTable) {
	if this.containsFTable(nameOrAlias) {
		return true, this.FTables[nameOrAlias]
	} else {
		return false, nil
	}
}

func (this *LineageContext) AddTable(table *AstTable) bool {
	if table == nil || table.Name == "" {
		return false
	}
	if this.containsFTable(table.Name) || this.containsFTable(table.Alias) {
		return false
	}
	this.FTables[table.Name] = table
	if table.Alias != "" {
		this.FTables[table.Alias] = table
	}
	return true
}

func (this *LineageContext) AddLineageContext(cxt *LineageContext) {
	this.ChildCxts = append(this.ChildCxts, cxt)
}

func (this *LineageContext) AddFiled(filed *AstFiled) {
	this.Fileds = append(this.Fileds, filed)
}
