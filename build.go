package autodb

import (
	"github.com/ijaychen/autodb/db"
	"log"
)

func BuildTables() {
	for _, table := range tables {
		table.Build()
	}
	for _, procedure := range Procedures {
		procedure.Build()
	}

	execSQL("call initdb", false)
}

func execSQL(sql string, echo bool) {
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
