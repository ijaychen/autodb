package column

type DateTimeColumnSt struct {
	Base
}

func (st *DateTimeColumnSt) IsEqual(info *MysqlColumn) bool {
	return st.Type == info.Type
}

func (st *DateTimeColumnSt) IsCompatible(info *MysqlColumn) bool {
	return st.IsEqual(info)
}

func NewDateTimeColumn(name string, comment string) IColumn {
	column := &DateTimeColumnSt{}
	column.Name = name
	column.Type = "datetime"
	column.Comment = comment
	column.Default = "'1970-01-01'"
	return column
}
