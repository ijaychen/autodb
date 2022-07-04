package autodb

type DateTimeColumnSt struct {
	ColumnBaseSt
}

func (st *DateTimeColumnSt) IsEqual(info *MysqlColumnSt) bool {
	return st.Type == info.Type
}

func (st *DateTimeColumnSt) IsCompatible(info *MysqlColumnSt) bool {
	return st.IsEqual(info)
}

func NewDateTimeColumn(name string, comment string) ColumnInterface {
	column := &DateTimeColumnSt{}
	column.Name = name
	column.Type = "datetime"
	column.Comment = comment
	column.Default = "'1970-01-01'"
	return column
}
