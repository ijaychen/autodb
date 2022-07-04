package column

import (
	"fmt"
	"strconv"
	"strings"
)

type StringColumnSt struct {
	Base
}

func find(s string) string {
	start := strings.Index(s, "(")
	end := strings.Index(s, ")")
	if start >= 0 && end > start {
		return s[start+1 : end]
	}
	return ""
}

//是否兼容
func (st *StringColumnSt) IsCompatible(info *MysqlColumn) bool {
	size := find(info.Type)
	//r := regexp.MustCompile(`\(.*?\)`)
	//size := r.FindStringSubmatch(info.Type)
	if len(size) <= 0 {
		return false
	}

	if i, err := strconv.Atoi(size); nil == err {
		if st.Size < i {
			return false
		}
	}

	return true
}

func NewStringColumn(name string, size int, comment string, def string) IColumn {
	column := &StringColumnSt{}
	column.Name = name
	column.Type = fmt.Sprintf("varchar(%d)", size)
	column.Size = size
	column.Comment = comment
	column.Default = def
	return column
}
