package autodb

import (
	"fmt"
)

type ColumnInterface interface {
	GetName() string
	CreateColumnSQL() string
	AddColumnSQL(tblName string) string
	ChangeColumnSQL(tblName string) string
	SetPlace(before string)
	SetFirst()
	IsEqual(info *MysqlColumnSt) bool
	IsCompatible(info *MysqlColumnSt) bool
	IsAutoIncrement() bool
}

type ColumnBaseSt struct {
	Name          string //字段名字
	Type          string //字段类型
	Null          string //是否默认为空
	Key           string //key
	Extra         string //extra
	Default       string //默认值
	Comment       string //注释
	Size          int
	First         bool
	Before        string
	AutoIncrement bool
}

func (st *ColumnBaseSt) GetName() string {
	return st.Name
}

func (st *ColumnBaseSt) IsAutoIncrement() bool {
	return st.AutoIncrement
}

func (st *ColumnBaseSt) CreateColumnSQL() string {
	def := "not null" + st.Extra
	if len(st.Default) > 0 {
		def = "default " + st.Default
	}
	return fmt.Sprintf("%s %s %s comment '%s'", st.Name, st.Type, def, st.Comment)
}

func (st *ColumnBaseSt) AddColumnSQL(tblName string) string {
	def := "not null" + st.Extra
	if len(st.Default) > 0 {
		def = "default " + st.Default
	}
	return fmt.Sprintf("alter table %s add %s %s %s comment '%s';", tblName, st.Name, st.Type, def, st.Comment)
}

func (st *ColumnBaseSt) ChangeColumnSQL(tblName string) string {
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

func (st *ColumnBaseSt) SetPlace(before string) {
	st.Before = before
}

func (st *ColumnBaseSt) SetFirst() {
	st.First = true
}

func (st *ColumnBaseSt) IsEqual(info *MysqlColumnSt) bool {
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
func (st *ColumnBaseSt) IsCompatible(info *MysqlColumnSt) bool {
	if st.Type != info.Type {
		return false
	}

	return true
}
