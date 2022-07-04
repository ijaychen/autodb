package autodb

import (
	"autodb/db"
	"fmt"
	"log"
)

func BuildTables() {
	for _, table := range tables {
		table.Build()
	}
	for _, procedure := range Procedures {
		procedure.Build()
	}

	Exec("call initdb", false)
}

func Exec(sql string, echo bool) {
	if len(sql) <= 0 {
		return
	}
	if echo {
		log.Printf("%s\n", sql)
	}
	_, err := db.OrmEngine.Exec(sql)
	if nil != err {
		log.Fatalf("%s", err)
	}
}

func DropColumnSQL(tblName string, columnName string) string {
	return fmt.Sprintf("alter table %s drop column %s;", tblName, columnName)
}

func DropKeySQL(tblName string, columnName string, kt string) string {
	switch kt {
	case MUL:
		return fmt.Sprintf("alter table %s drop index %s", tblName, columnName)
	case PRI:
		return fmt.Sprintf("alter table %s drop primary key", tblName)
	case UNI:
		return fmt.Sprintf("alter table %s drop index %s", tblName, columnName)
	default:
		log.Fatalf("drop key sql error")
	}
	return ""
}
