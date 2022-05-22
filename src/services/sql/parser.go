package sql

import (
	"SqlLineage/src/models/do"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/xwb1989/sqlparser"
)

// sql转Ast入口
func ASTparser(sql string) (*LineageContext, error) {
	parse, err := sqlparser.Parse(sql)
	if err != nil {
		logrus.Errorf("该sql无法解析： %s", sql)
		return nil, err
	}
	context := NewLineageContext()
	statementHandler(parse, context)
	return context, nil
}

func statementHandler(statement sqlparser.Statement, cxt *LineageContext) {
	switch statement.(type) {
	case *sqlparser.Select:
		selectHandler(statement.(*sqlparser.Select), cxt)
	case *sqlparser.Union:
		fmt.Println(456)
	}

}

// 单条Select解析
func selectHandler(sl *sqlparser.Select, cxt *LineageContext) {
	// 处理 from
	tablesExprHandler(sl.From, cxt)

	// 处理 col
	colsExprHandler(sl.SelectExprs, cxt)

	// 如果cxt存在Table,说明本SELECT是子查询的结果
	// 维护
	if cxt.Table != nil {
		for _, filed := range cxt.Fileds {
			cxt.Table.Fileds = append(cxt.Table.Fileds, filed)
		}
	}
}

// SELECT 列字段解析
func colsExprHandler(cols sqlparser.SelectExprs, cxt *LineageContext) {
	for _, col := range cols {
		// SelectExprs
		switch col.(type) {
		case *sqlparser.AliasedExpr:
			aliasedExprHandler(col.(*sqlparser.AliasedExpr), cxt)
		case *sqlparser.Nextval:
		case *sqlparser.StarExpr:
			starExprHandler(col.(*sqlparser.StarExpr), cxt)
		}
	}
}

// cols: 非通配符解析
func aliasedExprHandler(expr *sqlparser.AliasedExpr, cxt *LineageContext) {
	var filed *AstFiled
	// func
	// select
	innerExpr := expr.Expr
	switch innerExpr.(type) {
	case *sqlparser.ColName:
		// a.x
		// x
		filed = new(AstFiled)
		colExpr := innerExpr.(*sqlparser.ColName)
		filed.Name = colExpr.Name.String()
		if !colExpr.Qualifier.IsEmpty() {
			// 获取所属表
			if isContains, astTable := cxt.GetTable(colExpr.Qualifier.Name.String()); isContains {
				filed.Owner = astTable
				for _, astFiled := range astTable.Fileds {
					if astFiled.Name == filed.Name {
						filed.VertexId = astFiled.VertexId
						break
					}
				}
			}
			if filed.VertexId == 0 {
				panic("没有该数据")
			}
		} else {
			// 检验 From是否有唯一一个id
			temp := 0
			for _, table := range cxt.FTables {
				for _, tempFiled := range table.Fileds {
					if tempFiled.Name == filed.Name {
						temp++
						if temp == 1 {
							filed.Owner = table
							for _, astFiled := range table.Fileds {
								if astFiled.Name == filed.Name {
									filed.VertexId = astFiled.VertexId
									break
								}
							}
						} else {
							fmt.Errorf("存在多个FIlED, id：%s", filed.Name)
							return
						}
					}
				}
			}
		}
	case *sqlparser.FuncExpr:
		funcExpr := innerExpr.(*sqlparser.FuncExpr)
		fmt.Println(funcExpr)
	case *sqlparser.Subquery:
		// 子查询
		subTable := subQueryHandler(innerExpr.(*sqlparser.Subquery), cxt)
		if len(subTable.Fileds) != 1 {
			fmt.Errorf("子查询存在多个FIlED")
			return
		}
		filed = subTable.Fileds[0]
	case *sqlparser.ParenExpr:
		// (filed1 + filed2)
		parenExpr := innerExpr.(*sqlparser.ParenExpr)
		fmt.Println(parenExpr)
	case *sqlparser.SQLVal:
	}

	// 别名
	if !expr.As.IsEmpty() && filed != nil {
		filed.Alias = expr.As.String()
	}

	// 添加filed
	cxt.AddFiled(filed)
}

// cols： * 通配符解析
func starExprHandler(star *sqlparser.StarExpr, cxt *LineageContext) {
	if star.TableName.IsEmpty() {
		// *
		for _, table := range cxt.FTables {
			for _, filed := range table.Fileds {
				cxt.AddFiled(filed)
			}
		}
	} else {
		// a.*
		if isContain, astTable := cxt.GetTable(star.TableName.Name.String()); isContain {
			for _, filed := range astTable.Fileds {
				cxt.AddFiled(filed)
			}
		}
	}

}

// SELECT from 字段解析
func tablesExprHandler(tables sqlparser.TableExprs, cxt *LineageContext) {
	for _, table := range tables {
		tableExprHandler(&table, cxt)
	}
}

func tableExprHandler(expr *sqlparser.TableExpr, cxt *LineageContext) {
	tableExpr := *expr
	switch tableExpr.(type) {
	case *sqlparser.AliasedTableExpr:
		aliasedTableHandler(tableExpr.(*sqlparser.AliasedTableExpr), cxt)
	case *sqlparser.ParenTableExpr:
		parenTableHandler(tableExpr.(*sqlparser.ParenTableExpr), cxt)
	case *sqlparser.JoinTableExpr:
		joinTableHandler(tableExpr.(*sqlparser.JoinTableExpr), cxt)
	}
}

// tables: table 和 子查询
func aliasedTableHandler(expr *sqlparser.AliasedTableExpr, cxt *LineageContext) {
	t := new(AstTable)
	simpleTableExpr := expr.Expr
	switch simpleTableExpr.(type) {
	case sqlparser.TableName:
		name := simpleTableExpr.(sqlparser.TableName)
		t = tableNameHandler(&name, cxt)
	case *sqlparser.Subquery:
		t = subQueryHandler(simpleTableExpr.(*sqlparser.Subquery), cxt)
	}
	//tableExpr := expr.Expr.(sqlparser.TableName)
	//t.Name = tableExpr.Name.String()
	if !expr.As.IsEmpty() {
		t.Alias = expr.As.String()
	}
	cxt.AddTable(t)
}

// froms: 括号括起来多个
func parenTableHandler(expr *sqlparser.ParenTableExpr, cxt *LineageContext) {
	exprs := expr.Exprs
	tablesExprHandler(exprs, cxt)
}

// tables: 连接查询
func joinTableHandler(expr *sqlparser.JoinTableExpr, cxt *LineageContext) {
	tableExprHandler(&expr.LeftExpr, cxt)
	tableExprHandler(&expr.RightExpr, cxt)
}

func tableNameHandler(tn *sqlparser.TableName, cxt *LineageContext) *AstTable {
	table := new(AstTable)
	if tn.Qualifier.IsEmpty() {
		logrus.Errorf("%s必须要有所有数据库", tn.Name)
		return nil
	}
	id, db, isMatch := ParserDatabaseName(tn.Qualifier.String())
	if !isMatch {
		logrus.Errorf("%s不合法", tn.Qualifier.String())
		return nil
	}
	dataset := SelectDataset(id, db, tn.Name.String())
	table.VertexId = dataset.VertexId
	if dataset == nil {
		logrus.Errorf("%s不存在该表", tn.Name)
		return nil
	}
	columns := do.SelectColumnsByDatasetID(dataset.ID)
	table.Name = tn.Name.String()
	for _, column := range columns {
		table.AddFiled(NewAstFiledByName(column.Name, column.VertexId, table))
	}
	return table
}

func subQueryHandler(sb *sqlparser.Subquery, cxt *LineageContext) *AstTable {
	table := new(AstTable)
	statement := sb.Select
	context := NewLineageContext()
	context.Table = table
	cxt.AddLineageContext(context)
	statementHandler(statement, context)
	return table
}
