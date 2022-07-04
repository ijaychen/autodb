package column

import (
	"fmt"
)

type IColumn interface {
	GetName() string
	CreateColumnSQL() string
	AddColumnSQL(tblName string) string
	ChangeColumnSQL(tblName string) string
	SetPlace(before string)
	SetFirst()
	IsEqual(info *MysqlColumn) bool
	IsCompatible(info *MysqlColumn) bool
	IsAutoIncrement() bool
}

type MysqlColumn struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}

type Base struct {
	Name          string //字段名字
	Type          string //字段类型
	Null          string //是否默认为空
	Key           string //columnkey
	Extra         string //extra
	Default       string //默认值
	Comment       string //注释
	Size          int
	First         bool
	Before        string
	AutoIncrement bool
}

func (st *Base) GetName() string {
	return st.Name
}

func (st *Base) IsAutoIncrement() bool {
	return st.AutoIncrement
}

func (st *Base) CreateColumnSQL() string {
	def := "not null" + st.Extra
	if len(st.Default) > 0 {
		def = "default " + st.Default
	}
	return fmt.Sprintf("%s %s %s comment '%s'", st.Name, st.Type, def, st.Comment)
}

func (st *Base) AddColumnSQL(tblName string) string {
	def := "not null" + st.Extra
	if len(st.Default) > 0 {
		def = "default " + st.Default
	}
	return fmt.Sprintf("alter table %s add %s %s %s comment '%s';", tblName, st.Name, st.Type, def, st.Comment)
}

func (st *Base) ChangeColumnSQL(tblName string) string {
	def := "not null" + st.Extra
	if len(st.Default) > 0 {
		def = "default " + st.Default
	}
	afterSQL := ""
	if len(st.Before) > 0 {
		afterSQL = fmt.Sprintf("after %s", st.Before)
	} else if st.First {
		afterSQL = "first"
	}
	return fmt.Sprintf("alter table %s modify %s %s %s comment '%s' %s;", tblName, st.Name, st.Type, def, st.Comment, afterSQL)
}

func (st *Base) SetPlace(before string) {
	st.Before = before
}

func (st *Base) SetFirst() {
	st.First = true
}

func (st *Base) IsEqual(info *MysqlColumn) bool {
	if st.Default != info.Default {
		if !(st.Default == "null" && info.Null == "YES") {
			return false
		}
	}
	if st.Type != info.Type {
		return false
	}
	return true
}

//是否兼容
func (st *Base) IsCompatible(info *MysqlColumn) bool {
	if st.Type != info.Type {
		return false
	}

	return true
}
