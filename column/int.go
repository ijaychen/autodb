package column

import (
	"github.com/ijaychen/autodb"
	"strings"
)

type IntColumnSt struct {
	Base
}

func (st *IntColumnSt) GetSize(ct string) int {
	switch ct {
	case autodb.Int:
		return 10
	case autodb.TinyInt:
		return 3
	case autodb.SmallInt:
		return 5
	case autodb.BigInt:
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
		column.Type = autodb.TinyInt
	} else if strings.Contains(t, "small") {
		column.Type = autodb.SmallInt
	} else if strings.Contains(t, "big") {
		column.Type = autodb.BigInt
	} else {
		column.Type = autodb.Int
	}

	if unsigned {
		column.Type += " unsigned"
	}
	return column
}
