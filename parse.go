package autodb

import (
	"log"
	"strings"
)

func Parse() {
	for _, line := range Tables {
		table := NewTableSt(line.Name, line.Comment, len(line.Fields))
		for _, field := range line.Fields {
			var column ColumnInterface
			if field.Type == "varchar" {
				column = NewStringColumn(field.Name, field.Size, field.Comment, field.Default)
			} else if field.Type == "datetime" {
				column = NewDateTimeColumn(field.Name, field.Comment)
			} else if strings.Contains(field.Type, "blob") {
				column = NewBlobColumn(field.Name, field.Type, field.Comment)
			} else if strings.Contains(field.Type, "int") {
				column = NewIntColumn(field.Name, field.Type, field.Unsigned, field.Comment, field.AutoIncrement, field.Default)
			} else {
				log.Fatalf("字段类型定义错误。 table:%s, field:%s", line.Name, field.Name)
			}
			table.AddColumn(column)
		}
		for _, key := range line.Keys {
			table.AddKey(KeyFuncMap[key.Type](key.Name))
		}
		tables[line.Name] = table
	}
}
