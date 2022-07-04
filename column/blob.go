package column

import (
	"github.com/ijaychen/autodb"
	"log"
)

type BlobColumnSt struct {
	Base
}

func (st *BlobColumnSt) GetSize(ct string) int {
	switch ct {
	case autodb.TINYBLOB:
		return 256
	case autodb.BLOB:
		return 65535
	case autodb.MEDIUMBLOB:
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
	if t != autodb.TINYBLOB && t != autodb.BLOB && t != autodb.MEDIUMBLOB {
		log.Fatalf("column type error!! name:%s t:%s comment:%s", name, t, comment)
	}
	column.Type = t
	column.Default = "null"
	return column
}
