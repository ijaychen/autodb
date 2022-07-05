package autodb

import (
	"fmt"
	"github.com/ijaychen/autodb/columnkey"
	"strings"
)

type MysqlKeyDesc struct {
	NonUnique  int
	KeyName    string
	SeqInIndex int
	ColumnName string
}

type TableKey struct {
	Name      string   // 索引名字
	Type      string   // 索引类型
	ColumnVec []string // 字段列表
}

func createTableKeyByMysqlData(desc *MysqlKeyDesc) *TableKey {
	key := &TableKey{
		Name:      desc.KeyName,
		ColumnVec: nil,
	}
	switch desc.KeyName {
	case columnkey.PRI:
		key.Type = columnkey.PRI
	default:
		if 0 == desc.NonUnique {
			key.Type = columnkey.UNI
		} else {
			key.Type = columnkey.MUL
		}
	}
	key.ColumnVec = append(key.ColumnVec, desc.ColumnName)
	return key
}

func (key *TableKey) Copy() *TableKey {
	ret := &TableKey{
		Name:      key.Name,
		Type:      key.Type,
		ColumnVec: make([]string, len(key.ColumnVec)),
	}
	copy(ret.ColumnVec, key.ColumnVec)
	return ret
}

func (key *TableKey) CreateKeySQL() string {
	size := len(key.ColumnVec)
	if size <= 0 {
		return ""
	}
	switch key.Type {
	case columnkey.PRI:
		return key.createPriKeySQL()
	case columnkey.UNI:
		return key.createUniKeySQL()
	case columnkey.MUL:
		return key.createMulKeySQL()
	}
	return ""
}

func (key *TableKey) AddKeySQL(tblName string) string {
	size := len(key.ColumnVec)
	if size <= 0 {
		return ""
	}
	switch key.Type {
	case columnkey.PRI:
		return key.addPriKeySQL(tblName)
	case columnkey.UNI:
		return key.addUniKeySQL(tblName)
	case columnkey.MUL:
		return key.addMulKeySQL(tblName)
	}
	return ""
}

func (key *TableKey) createPriKeySQL() string {
	if size := len(key.ColumnVec); size > 0 {
		columns := strings.Join(key.ColumnVec, ", ")
		return fmt.Sprintf("PRIMARY KEY (%s)", columns)
	}
	return ""
}

func (key *TableKey) addPriKeySQL(tblName string) string {
	if size := len(key.ColumnVec); size > 0 {
		columns := strings.Join(key.ColumnVec, ", ")
		return fmt.Sprintf("ALTER TABLE %s ADD PRIMARY KEY (%s)", tblName, columns)
	}
	return ""
}

func (key *TableKey) createUniKeySQL() string {
	if size := len(key.ColumnVec); size > 0 {
		columns := strings.Join(key.ColumnVec, ", ")
		return fmt.Sprintf("UNIQUE KEY %s(%s)", key.Name, columns)
	}
	return ""
}

func (key *TableKey) addUniKeySQL(tblName string) string {
	if size := len(key.ColumnVec); size > 0 {
		columns := strings.Join(key.ColumnVec, ", ")
		return fmt.Sprintf("ALTER TABLE %s ADD UNIQUE %s (%s)", tblName, key.Name, columns)
	}
	return ""
}

func (key *TableKey) createMulKeySQL() string {
	if size := len(key.ColumnVec); size > 0 {
		columns := strings.Join(key.ColumnVec, ", ")
		return fmt.Sprintf("KEY %s(%s)", key.Name, columns)
	}
	return ""
}

func (key *TableKey) addMulKeySQL(tblName string) string {
	if size := len(key.ColumnVec); size > 0 {
		columns := strings.Join(key.ColumnVec, ", ")
		return fmt.Sprintf("ALTER TABLE %s ADD INDEX %s (%s)", tblName, key.Name, columns)
	}
	return ""
}
