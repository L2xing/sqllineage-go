package test

import (
	"SqlLineage/src/services/sql"
	"fmt"
	"testing"
)

/*
文件名必须是_test.go结尾的，这样在执行go test的时候才会执行到相应的代码
你必须import testing这个包
所有的测试用例函数必须是Test开头
测试用例会按照源代码中写的顺序依次执行
测试函数TestXxx()的参数是testing.T，我们可以使用该类型来记录错误或者是测试状态
测试格式：func TestXxx (t *testing.T),Xxx部分可以为任意的字母数字的组合，但是首字母不能是小写字母[a-z]，例如Testintdiv是错误的函数名。
函数中通过调用testing.T的Error, Errorf, FailNow, Fatal, FatalIf方法，说明测试不通过，调用Log方法用来记录测试的信息。
*/

func TestInt(t *testing.T) {
	println("这是一个测试用例")
}

var simpleAllSelectSQL = "select * from T1"
var simpleMutiTablesSelectSQL = "select * from T1, T2"
var simpleThreeTableSQL = "SELECT a.*, b.*, c.* FROM T1 a,T2 b,T3 c"
var complexSelectCloms = "SELECT a.*, b.id, name n, *, max(a.id, b.id) as sum, (select id from t4) as t4Id from t1 a, t2 b"

//Alias
var simpleMutiTablesaAliasSelectSQL = "select * from T1 a, T2 b"
var simpleMutiTablesaAliasSelectSQL2 = "SELECT * from (instance, `column`)"
var complexFromSubQuery = "select * from (select a.id from instance a UNION SELECT b.id from instance b) c"

// 链接方式
var simpleLeftSelect = "SELECT column_name(s)\n" + "FROM table1\n" + "LEFT JOIN table2\n" + "ON table1.column_name=table2.column_name;"
var simpleRightSelect = "SELECT column_name(s)\n" + "FROM table1\n" + "RIGHT JOIN table2\n" + "ON table1.column_name=table2.column_name;"
var simpleInnerSelect = "SELECT column_name(s)\n" + "FROM table1\n" + "INNER JOIN table2\n" + "ON table1.column_name=table2.column_name;"
var complexInnerSelect = "SELECT * FROM `user` LEFT JOIN (select * from dataset) d ON `user`.id = d.id"

// 子查询
var simpleChildSelect = "select a.*, b.* from T1 a, (select id, name from T1) b"

var simpleSlecetSql = "Select id,name from T1 where id = 2"

var complexSelectOne = "select (select id from t1) as id, b.name from T3 a, (select id, name from T4) b;"
var complexSelectTwo = "select (select id from t3) as id, b.name from T3 a, (select id, name from T4) b;"
var complexSelectThree = "select (select id from t3) as id, b.name from T3 a, (select (select id from t1) as id, name from T4) b;"

// union
var selectUnion = "SELECT id, `name`, database_id from dataset UNION SELECT id, `name`, instance_id from `database` "
var selectUnion3 = "SELECT id, `name`, database_id from dataset \nUNION \nSELECT id, `name`, instance_id from `database`\nUNION \nSELECT id, `name`, dataset_id from `column`"

// col
var selectCols = "SELECT *, a.id, t.id, id as n, sum(id) AS i, (select id from user limit 1)as id2, a.* FROM user a, teacher t"
var selectCols1 = "SELECT sum(d1+d2) as a,(d1+(SELECT d1 FROM test LIMIT 1)), \"true\" FROM test\n"

// 展示
var realSelect = "SELECT * FROM sqllineage.`column`, sqllineage.dataset"
var realSelect2 = "SELECT id FROM sqllineage.`column`"

var realSelectMy1 = "SELECT id FROM db8@sqllineage.`column`"
var realSelectMy2 = "SELECT id FROM db8:sqllineage.`column`"
var realSelectMy3 = "SELECT id FROM db8_sqllineage.`column`"
var realSelectMy4 = "SELECT c.id, u.id FROM db8@sqllineage.`column` c, db8@sqllineage.`user` u"

func TestSql(t *testing.T) {
	//str := selectUnion
	//str = "select id, name from T1"s
	//query, err := sqlparser.Parse(str)
	//of := reflect.TypeOf(query)
	//println(of.Name())
	//println(of.String())
	//switch query.(type) {
	//	//case sqlparser.Select:
	//	//	fmt.Println("select 类型")
	//	//case sqlparser.Union:
	//	//	fmt.Println("union 类型")
	//	//}
	//	if err != nil {
	//		log.Fatal(err)
	//	}
}

func TestSql2(t *testing.T) {
	SetUp()
	str := realSelectMy4
	println(str)
	tparser, err := sql.ASTparser(str)
	if err != nil {
		println("解析失败")
	} else {
		sql.PrintAstTables(tparser)
		sql.PrintAstCols(tparser)
		columns, _ := sql.GetAstColumns(tparser)
		for _, v := range columns {
			fmt.Print(v)
			fmt.Print(" ")
		}
	}
	fmt.Println()
}
