package autodb

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var (
	Tables     []*TableConf
	Procedures []*ProcedureSt
)

type (
	FieldConf struct {
		Name          string
		Type          string
		Comment       string
		AutoIncrement bool
		Size          int
		Unsigned      bool
		Default       string
	}

	TableConf struct {
		Name    string
		Comment string
		Fields  []*FieldConf
		Keys    []*FieldKey
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
