package autodb

import (
	"encoding/json"
	"github.com/ijaychen/autodb/column"
	"github.com/ijaychen/autodb/columnkey"
	"io/ioutil"
	"log"
	"strings"
)

var (
	Tables     []*TableConf
	Procedures []*ProcedureSt
)

type (
	ColumnConf struct {
		Name          string
		Type          string
		Comment       string
		AutoIncrement bool
		Size          int
		Unsigned      bool
		Default       string
	}

	KeyConf struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}

	TableConf struct {
		Name    string
		Comment string
		Columns []*ColumnConf
		Keys    []*KeyConf
	}
)

func loadTableConf(file string) bool {
	data, err := ioutil.ReadFile(file)
	if nil != err {
		log.Printf("load table conf error! [%s] %s\n", file, err)
		return false
	}

	if err = json.Unmarshal(data, Tables); err != nil {
		log.Fatalf("load %s Unmarshal json error:%s", file, err)
		return false
	}
	return true
}

func ParseTableConf() {
	for _, line := range Tables {
		table := NewTableSt(line.Name, line.Comment, len(line.Columns))
		for _, field := range line.Columns {
			var col column.IColumn
			if field.Type == "varchar" {
				col = column.NewStringColumn(field.Name, field.Size, field.Comment, field.Default)
			} else if field.Type == "datetime" {
				col = column.NewDateTimeColumn(field.Name, field.Comment)
			} else if strings.Contains(field.Type, "blob") {
				col = column.NewBlobColumn(field.Name, field.Type, field.Comment)
			} else if strings.Contains(field.Type, "int") {
				col = column.NewIntColumn(field.Name, field.Type, field.Unsigned, field.Comment, field.AutoIncrement, field.Default)
			} else {
				log.Fatalf("字段类型定义错误。 table:%s, field:%s", line.Name, field.Name)
			}
			table.AddColumn(col)
		}
		for _, key := range line.Keys {
			if ik := columnkey.CreateKey(key.Type, key.Name); nil != ik {
				table.AddKey(ik)
			}
		}
		tables[line.Name] = table
	}
}
