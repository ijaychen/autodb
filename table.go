package autodb

import (
	"fmt"
	"github.com/ijaychen/autodb/db"
	"log"
	"strings"
)

var tables = make(map[string]*TableSt)

type MysqlColumnSt struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}

type TableInfoSt struct {
	Sequence []*MysqlColumnSt
	Columns  map[string]*MysqlColumnSt
}

type TableSt struct {
	Name     string
	Comment  string
	Columns  map[string]ColumnInterface
	Sequence []ColumnInterface
	Exists   bool
	pris     []string
	keys     map[string]*FieldKey
	sqlInfo  *TableInfoSt
	hasData  int8
}

func NewTableSt(name, comment string, fieldCount int) *TableSt {
	table := &TableSt{
		Name: name, Comment: comment, hasData: -1,
	}
	table.Sequence = make([]ColumnInterface, 0, fieldCount)
	table.Columns = make(map[string]ColumnInterface)
	table.keys = make(map[string]*FieldKey)
	return table
}

func (st *TableSt) AddKey(key *FieldKey) {
	if key.Type == PRI {
		st.pris = append(st.pris, key.Name)
	}
	st.keys[key.Name] = key
}

func (st *TableSt) CreatePRIKeySQL() string {
	if len(st.pris) <= 0 {
		return ""
	}
	names := strings.Join(st.pris, ", ")
	st.pris = make([]string, 0)
	return fmt.Sprintf("PRIMARY KEY (%s)", names)
}

func (st *TableSt) AddPRIKeySQL() string {
	if len(st.pris) <= 0 {
		return ""
	}
	names := strings.Join(st.pris, ", ")
	st.pris = make([]string, 0)
	return fmt.Sprintf("ALTER TABLE %s ADD PRIMARY KEY (%s)", st.Name, names)
}

func (st *TableSt) AddColumn(column ColumnInterface) {
	name := column.GetName()
	if _, exists := st.Columns[name]; exists {
		log.Fatalf("%s表中含有重复字段%s", st.Name, name)
	}
	st.Columns[name] = column
	st.Sequence = append(st.Sequence, column)
}

func (st *TableSt) CreateTableSQL() string {
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

func (st *TableSt) HasTable() bool {
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

func (st *TableSt) HasData() bool {
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

func (st *TableSt) Build() {
	st.Check()
	if !st.HasTable() {
		st.Create()
	} else {
		st.Change()
	}
}

func (st *TableSt) Check() {
	if 0 == len(st.Columns) {
		log.Fatalf("%s表配置中没有字段", st.Name)
	}
	for name, key := range st.keys {
		if _, exists := st.Columns[key.Name]; !exists {
			log.Fatalf("%s %s 该key的对应字段不存在", st.Name, name)
		}
	}

	// there can be only one auto column and it must be defined as a key
	autoIncrement := false
	for name, column := range st.Columns {
		if column.IsAutoIncrement() {
			// only one
			if autoIncrement {
				log.Fatalf("%s只能有一个auto_increment字段", name)
			}
			// must be defined as a key
			if _, exists := st.keys[name]; !exists {
				log.Fatalf("%s auto_increment字段必须有索引", name)
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
		if nil != key && info.Key != NO {
			if !key.IsEqual(info) {
				log.Fatalf("%s %s不能修改字段约束", st.Name, info.Field)
			}
		}
	}
}

func (st *TableSt) Create() {
	Exec(st.CreateTableSQL(), true)
}

func (st *TableSt) Change() {
	bChange := false
	info := st.GetTableColumnInfo(false)
	columns := info.Columns

	for name, column := range st.Columns {
		//该列已经存在，检查是否需要修改
		if mysqlInfo, ok := columns[name]; ok {
			if !column.IsEqual(mysqlInfo) {
				Exec(column.ChangeColumnSQL(st.Name), true)
			}
		} else { //不存在则添加
			Exec(column.AddColumnSQL(st.Name), true)
			bChange = true
		}
	}
	//需要删除的column和key
	for _, mysqlInfo := range columns {
		if _, ok := st.Columns[mysqlInfo.Field]; !ok {
			Exec(DropColumnSQL(st.Name, mysqlInfo.Field), true)
			bChange = true
			//column不需要删除时再检测key (因为column删除了key也会删除)
		} else if mysqlInfo.Key != NO {
			//这个key以前有现在没了
			if _, ok := st.keys[mysqlInfo.Field]; !ok {
				Exec(DropKeySQL(st.Name, mysqlInfo.Field, mysqlInfo.Type), true)
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
			log.Fatalf("table %s add key error!", st.Name)
		}
		if mysqlInfo.Key == NO {
			Exec(key.AddKeySQL(st), true)
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
			Exec(column.ChangeColumnSQL(st.Name), true)
		}
	}
}

/////////////////////////////////
func (st *TableSt) GetTableColumnInfo(reset bool) *TableInfoSt {
	if nil != st.sqlInfo && !reset {
		return st.sqlInfo
	}
	st.sqlInfo = new(TableInfoSt)
	st.sqlInfo.Columns = make(map[string]*MysqlColumnSt)
	err := db.OrmEngine.SQL(fmt.Sprintf("show columns from %s;", st.Name)).Find(&st.sqlInfo.Sequence)
	if nil != err {
		log.Fatalf("查询失败！error:%s", err)
	}
	for _, ret := range st.sqlInfo.Sequence {
		st.sqlInfo.Columns[ret.Field] = ret
	}
	return st.sqlInfo
}
