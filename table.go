package autodb

import (
	"fmt"
	"github.com/ijaychen/autodb/column"
	"github.com/ijaychen/autodb/columnkey"
	"github.com/ijaychen/autodb/db"
	"log"
	"sort"
	"strings"
)

var tables = make(map[string]*Table)

type MysqlTable struct {
	Columns map[string]*column.MysqlColumn
	Keys    map[string]*TableKey
}

type Table struct {
	Name     string
	Comment  string
	Columns  map[string]column.IColumn
	Sequence []column.IColumn
	Exists   bool
	tblKeys  map[string]*TableKey
	sqlInfo  *MysqlTable
	hasData  int8
}

func NewTableSt(name, comment string, fieldCount int) *Table {
	table := &Table{
		Name: name, Comment: comment, hasData: -1,
	}
	table.Sequence = make([]column.IColumn, 0, fieldCount)
	table.Columns = make(map[string]column.IColumn)
	table.tblKeys = make(map[string]*TableKey)
	return table
}

func (st *Table) GetName() string {
	return st.Name
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

	for _, key := range st.tblKeys {
		sql := key.CreateKeySQL()
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
	for name, key := range st.tblKeys {
		for _, column := range key.ColumnVec {
			if _, exists := st.Columns[column]; !exists {
				log.Fatalf("table[%s] column is not exists! key[%s]", st.Name, name)
			}
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

			found := false
			if keys, ok := st.tblKeys[columnkey.PRI]; ok {
				for _, line := range keys.ColumnVec {
					if line == name {
						found = true
						break
					}
				}
			}
			// must be defined as a columnKey
			if !found {
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

	// 检查字段
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
	}
	// 检查索引， 不删除已存在的索引
	for _, info := range sqlInfo.Keys {
		newKey, ok := st.tblKeys[info.Name]
		if !ok {
			log.Fatalf("%s %s不能删除已存在的索引", st.Name, info.Name)
		}
		if info.Type != newKey.Type {
			log.Fatalf("%s %s不能修改已存在的索引类型", st.Name, info.Name)
		}
		for _, line1 := range info.ColumnVec {
			for _, line2 := range newKey.ColumnVec {
				if line1 != line2 {
					log.Fatalf("%s %s不能修改已存在的索引", st.Name, info.Name)
				}
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
	//需要删除的column
	for _, mysqlInfo := range columns {
		if _, ok := st.Columns[mysqlInfo.Field]; !ok {
			st.CreateDropColumnSQL(mysqlInfo.Field)
			bChange = true
		}
	}

	// 表变化了重新查询一次数据库表信息
	if bChange {
		info = st.GetTableColumnInfo(true)
		columns = info.Columns
	}
	// key的添加放在删除后，否则可能会冲突
	// 已经检查过key的合法性，所以对不存在的key直接添加，这里不处理 PRIMARY KEY之前只有一个，现在有两个的问题
	for _, key := range st.tblKeys {
		if _, exist := info.Keys[key.Name]; !exist {
			execSQL(key.AddKeySQL(st.Name), true)
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

	columns := make([]*column.MysqlColumn, 0)
	err := db.OrmEngine.SQL(fmt.Sprintf("show columns from %s;", st.Name)).Find(&columns)
	if nil != err {
		log.Fatalf("查询失败！error:%s", err)
	}

	for _, ret := range columns {
		st.sqlInfo.Columns[ret.Field] = ret
	}

	// keys
	vec := make([]*MysqlKeyDesc, 0)
	db.OrmEngine.SQL(fmt.Sprintf("show index from %s;", st.Name)).Find(&vec)
	sort.Slice(vec, func(i, j int) bool {
		return vec[i].SeqInIndex < vec[j].SeqInIndex
	})
	st.sqlInfo.Keys = make(map[string]*TableKey)
	for _, line := range vec {
		if line.KeyName == "PRIMARY" {
			line.KeyName = columnkey.PRI
		}
		cur, exist := st.sqlInfo.Keys[line.KeyName]
		if !exist {
			cur = createTableKeyByMysqlData(line)
			st.sqlInfo.Keys[line.KeyName] = cur
		} else {
			cur.ColumnVec = append(cur.ColumnVec, line.ColumnName)
		}
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
