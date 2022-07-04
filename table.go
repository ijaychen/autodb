package autodb

import (
	"fmt"
	"github.com/ijaychen/autodb/column"
	"github.com/ijaychen/autodb/columnkey"
	"github.com/ijaychen/autodb/db"
	"github.com/ijaychen/autodb/iface"
	"log"
	"strings"
)

var tables = make(map[string]*Table)

type MysqlTable struct {
	Sequence []*column.MysqlColumn
	Columns  map[string]*column.MysqlColumn
}

type Table struct {
	Name     string
	Comment  string
	Columns  map[string]column.IColumn
	Sequence []column.IColumn
	Exists   bool
	pris     []string
	keys     map[string]iface.IKey
	sqlInfo  *MysqlTable
	hasData  int8
}

func NewTableSt(name, comment string, fieldCount int) *Table {
	table := &Table{
		Name: name, Comment: comment, hasData: -1,
	}
	table.Sequence = make([]column.IColumn, 0, fieldCount)
	table.Columns = make(map[string]column.IColumn)
	table.keys = make(map[string]iface.IKey)
	return table
}

func (st *Table) GetName() string {
	return st.Name
}

func (st *Table) AddKey(key iface.IKey) {
	name := key.GetName()
	if key.GetType() == columnkey.PRI {
		st.pris = append(st.pris, name)
	}
	st.keys[name] = key
}

func (st *Table) CreatePriKeySQL() string {
	if len(st.pris) <= 0 {
		return ""
	}
	names := strings.Join(st.pris, ", ")
	st.pris = make([]string, 0)
	return fmt.Sprintf("PRIMARY KEY (%s)", names)
}

func (st *Table) AddPriKeySQL() string {
	if len(st.pris) <= 0 {
		return ""
	}
	names := strings.Join(st.pris, ", ")
	st.pris = make([]string, 0)
	return fmt.Sprintf("ALTER TABLE %s ADD PRIMARY KEY (%s)", st.Name, names)
}

func (st *Table) AddColumn(column column.IColumn) {
	name := column.GetName()
	if _, exists := st.Columns[name]; exists {
		log.Fatalf("%s表中含有重复字段%s", st.Name, name)
	}
	st.Columns[name] = column
	st.Sequence = append(st.Sequence, column)
}

func (st *Table) CreateTableSQL() string {
	sqlVec := make([]string, 0, len(st.Sequence))
	for _, column := range st.Sequence {
		sqlVec = append(sqlVec, column.CreateColumnSQL())
	}

	for _, key := range st.keys {
		sql := key.CreateKeySQL(st)
		if len(sql) > 0 {
			sqlVec = append(sqlVec, sql)
		}
	}

	head := fmt.Sprintf(CreateSQLHead, st.Name)
	tail := fmt.Sprintf(CreateSQLTail, st.Comment)
	return head + strings.Join(sqlVec, ",\n") + tail
}

func (st *Table) HasTable() bool {
	if st.Exists {
		return true
	}
	var ret []string
	err := db.OrmEngine.SQL(fmt.Sprintf("show tables like '%s';", st.Name)).Find(&ret)
	if nil != err {
		log.Fatalf("查询失败！error:%s", err)
	}

	st.Exists = len(ret) > 0
	return st.Exists
}

func (st *Table) HasData() bool {
	if st.hasData != -1 {
		return st.hasData == 1
	}
	if !st.HasTable() {
		st.hasData = 0
		return false
	}
	rets, err := db.OrmEngine.QueryString(fmt.Sprintf("select * from %s limit 1;", st.Name))
	if nil != err {
		log.Fatalf("查询失败！%v", err)
	}
	if len(rets) > 0 {
		st.hasData = 1
	} else {
		st.hasData = 0
	}
	return st.hasData == 1
}

func (st *Table) Build() {
	st.Check()
	if !st.HasTable() {
		st.Create()
	} else {
		st.Change()
	}
}

func (st *Table) Check() {
	if 0 == len(st.Columns) {
		log.Fatalf("table[%s] is empty!", st.Name)
	}
	for name, key := range st.keys {
		if _, exists := st.Columns[key.GetName()]; !exists {
			log.Fatalf("table[%s] column is not exists! key[%s]", st.Name, name)
		}
	}

	// there can be only one auto column, and it must be defined as a columnKey
	autoIncrement := false
	for name, column := range st.Columns {
		if column.IsAutoIncrement() {
			// only one
			if autoIncrement {
				log.Fatalf("table[%s] %s can be only one auto column", st.Name, name)
			}
			// must be defined as a columnKey
			if _, exists := st.keys[name]; !exists {
				log.Fatalf("column %s auto_increment must be index", name)
			}
			autoIncrement = true
		}
	}

	// 如果表不存在，不用做后续的兼容判定
	if !st.HasTable() {
		return
	}

	sqlInfo := st.GetTableColumnInfo(false)
	if nil == sqlInfo {
		log.Fatalf("%s table column info is nil", st.Name)
	}
	for _, info := range sqlInfo.Columns {
		column, exists := st.Columns[info.Field]
		//以前有现在也有的column需要兼容判定
		if exists {
			if !column.IsEqual(info) {
				if !column.IsCompatible(info) {
					log.Fatalf("%s %s字段已存在，新旧类型不兼容", st.Name, column.GetName())

				}
			}
		} else { //以前有现在没有，需要删除
			//只有没有数据的表可以删
			if st.HasData() {
				log.Fatalf("%s 表已使用，不允许删除字段", st.Name)
			}
		}
		key := st.keys[info.Field]
		//对于已经存在的约束不能修改, 断言约束相同
		if nil != key && info.Key != columnkey.NO {
			if !key.IsEqual(info.Key) {
				log.Fatalf("%s %s不能修改字段约束", st.Name, info.Field)
			}
		}
	}
}

func (st *Table) Create() {
	execSQL(st.CreateTableSQL(), true)
}

func (st *Table) Change() {
	bChange := false
	info := st.GetTableColumnInfo(false)
	columns := info.Columns

	for name, column := range st.Columns {
		//该列已经存在，检查是否需要修改
		if mysqlInfo, ok := columns[name]; ok {
			if !column.IsEqual(mysqlInfo) {
				execSQL(column.ChangeColumnSQL(st.Name), true)
			}
		} else { //不存在则添加
			execSQL(column.AddColumnSQL(st.Name), true)
			bChange = true
		}
	}
	//需要删除的column和key
	for _, mysqlInfo := range columns {
		if _, ok := st.Columns[mysqlInfo.Field]; !ok {
			st.CreateDropColumnSQL(mysqlInfo.Field)
			bChange = true
			//column不需要删除时再检测key (因为column删除了key也会删除)
		} else if mysqlInfo.Key != columnkey.NO {
			//这个key以前有现在没了
			if _, ok := st.keys[mysqlInfo.Field]; !ok {
				st.DropKeySQL(mysqlInfo.Field, mysqlInfo.Type)
			}
		}
	}

	if bChange {
		info = st.GetTableColumnInfo(true)
		columns = info.Columns
	}
	// key的添加放在删除后，否则可能会冲突
	// 已经检查过key的合法性，所以对不存在的key直接添加，这里不处理 PRIMARY KEY之前只有一个，现在有两个的问题
	for name, key := range st.keys {
		mysqlInfo, ok := columns[name]
		if !ok {
			log.Fatalf("table %s add columnkey error!", st.Name)
		}
		if mysqlInfo.Key == columnkey.NO {
			execSQL(key.AddKeySQL(st), true)
		}
	}

	//当添加或删除过字段时,认为字段顺序可能不一致
	//比对顺序太麻烦了，直接全部change一遍
	if bChange {
		for i, column := range st.Sequence {
			if i == 0 {
				column.SetFirst()
			} else {
				column.SetPlace(st.Sequence[i-1].GetName())
			}
			execSQL(column.ChangeColumnSQL(st.Name), true)
		}
	}
}

func (st *Table) GetTableColumnInfo(reset bool) *MysqlTable {
	if nil != st.sqlInfo && !reset {
		return st.sqlInfo
	}
	st.sqlInfo = new(MysqlTable)
	st.sqlInfo.Columns = make(map[string]*column.MysqlColumn)
	err := db.OrmEngine.SQL(fmt.Sprintf("show columns from %s;", st.Name)).Find(&st.sqlInfo.Sequence)
	if nil != err {
		log.Fatalf("查询失败！error:%s", err)
	}
	for _, ret := range st.sqlInfo.Sequence {
		st.sqlInfo.Columns[ret.Field] = ret
	}
	return st.sqlInfo
}

func (st *Table) CreateDropColumnSQL(columnName string) {
	execSQL(fmt.Sprintf("alter table %s drop column %s;", st.Name, columnName), true)
}

func (st *Table) DropKeySQL(columnName, kt string) {
	var sql string
	switch kt {
	case columnkey.MUL:
		sql = fmt.Sprintf("alter table %s drop index %s", st.Name, columnName)
	case columnkey.PRI:
		sql = fmt.Sprintf("alter table %s drop primary columnkey", st.Name)
	case columnkey.UNI:
		sql = fmt.Sprintf("alter table %s drop index %s", st.Name, columnName)
	default:
		log.Fatalf("drop columnkey sql error")
	}
	if "" != sql {
		execSQL(sql, true)
	}
}
