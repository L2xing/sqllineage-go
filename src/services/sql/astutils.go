package sql

import (
	"bytes"
	"fmt"
)

func GetAstTables(cxt *LineageContext) ([]string, *map[string]*AstTable) {
	sets := make(map[string]*AstTable, 0)
	getTabelFromLineageContext(cxt, &sets)
	strings := make([]string, 0)
	for k, _ := range sets {
		strings = append(strings, k)
	}
	return strings, &sets
}

func GetAstColumns(cxt *LineageContext) ([]string, map[string]*AstFiled) {
	cols := make(map[string]int, 0)
	colAstColumns := make(map[string]*AstFiled, 0)
	colsArray := make([]string, 0)
	for _, filed := range cxt.Fileds {
		tempName := ""
		if filed.Alias != "" {
			tempName = filed.Alias
		} else if filed.Name != "" {
			tempName = filed.Name
		}
		if tempName == "" {
			continue
		}
		var num int
		if v, isOk := cols[tempName]; isOk {
			num = v + 1
		} else {
			num = 0
		}
		cols[tempName] = num
		if num > 0 {
			tempName = fmt.Sprintf("%s(%d)", tempName, num)
		}
		colsArray = append(colsArray, tempName)
		colAstColumns[tempName] = filed
	}
	return colsArray, colAstColumns
}

func getTabelFromLineageContext(cxt *LineageContext, maps *map[string]*AstTable) {
	sets := *(maps)
	for _, table := range cxt.FTables {
		if _, ok := sets[table.Name]; !ok {
			sets[table.Name] = table
		}
	}
	for _, childCxt := range cxt.ChildCxts {
		getTabelFromLineageContext(childCxt, maps)
	}
}

func PrintAstTables(cxt *LineageContext) {
	tables, _ := GetAstTables(cxt)
	bufferString := bytes.NewBufferString("(")
	i := len(tables)
	for _, table := range tables {
		bufferString.WriteString(table)
		if i > 1 {
			bufferString.WriteRune(',')
		}
		i--
	}
	bufferString.WriteRune(')')
	fmt.Println(bufferString.String())
}

func PrintAstCols(cxt *LineageContext) {
	bufferString := bytes.NewBufferString("(")
	for _, filed := range cxt.Fileds {
		bufferString.WriteString(filed.Owner.Name + "." + filed.Name)
		bufferString.WriteRune(',')
	}
	bufferString.WriteString(")")
	fmt.Println(bufferString.String())
}

func PrintColsLineage(cxt *LineageContext) {

}
