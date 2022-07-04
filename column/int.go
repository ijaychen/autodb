package column

import (
	"strings"
)

const (
	Int      = "int(10)"
	TinyInt  = "tinyint(3)"
	SmallInt = "smallint(5)"
	BigInt   = "bigint(20)"
)

type IntColumnSt struct {
	Base
}

func (st *IntColumnSt) GetSize(ct string) int {
	switch ct {
	case Int:
		return 10
	case TinyInt:
		return 3
	case SmallInt:
		return 5
	case BigInt:
		return 20
	default:
		return 999 //不是int类型，给一个很大的值
	}
}

func (st *IntColumnSt) IsCompatible(info *MysqlColumn) bool {
	if !st.Base.IsCompatible(info) {
		return false
	}
	if st.GetSize(st.Type) < st.GetSize(info.Type) {
		return false
	}

	return true
}

func NewIntColumn(name, t string, unsigned bool, comment string, increment bool, def string) IColumn {
	column := &IntColumnSt{}
	column.Name = name
	column.Comment = comment

	if increment {
		column.AutoIncrement = true
		column.Extra = " auto_increment"
	} else {
		if len(def) > 0 {
			column.Default = def
		} else {
			column.Default = "0"
		}
	}

	t = strings.ToLower(t)

	if strings.Contains(t, "tiny") {
		column.Type = TinyInt
	} else if strings.Contains(t, "small") {
		column.Type = SmallInt
	} else if strings.Contains(t, "big") {
		column.Type = BigInt
	} else {
		column.Type = Int
	}

	if unsigned {
		column.Type += " unsigned"
	}
	return column
}
