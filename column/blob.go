package column

import (
	"log"
)

const (
	TINYBLOB   = "tinyblob"
	BLOB       = "blob"
	MEDIUMBLOB = "mediumblob"
)

type BlobColumnSt struct {
	Base
}

func (st *BlobColumnSt) GetSize(ct string) int {
	switch ct {
	case TINYBLOB:
		return 256
	case BLOB:
		return 65535
	case MEDIUMBLOB:
		return 16776960
	default:
		return 999999999 //不是Blob类型，给一个很大的值
	}
}

func (st *BlobColumnSt) IsCompatible(info *MysqlColumn) bool {
	if !st.Base.IsCompatible(info) {
		return false
	}
	if st.GetSize(st.Type) < st.GetSize(info.Type) {
		return false
	}
	return true
}

func NewBlobColumn(name, t, comment string) IColumn {
	column := &BlobColumnSt{}
	column.Name = name
	column.Comment = comment
	if t != TINYBLOB && t != BLOB && t != MEDIUMBLOB {
		log.Fatalf("column type error!! name:%s t:%s comment:%s", name, t, comment)
	}
	column.Type = t
	column.Default = "null"
	return column
}
