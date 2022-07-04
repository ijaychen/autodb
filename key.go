package autodb

import (
	"fmt"
	"log"
)

type FieldKey struct {
	Type string
	Name string
}

func (st *FieldKey) CreateKeySQL(table *TableSt) string {
	if st.Type == MUL {
		return fmt.Sprintf("KEY %s(%s)", st.Name, st.Name)
	} else if st.Type == PRI {
		return table.CreatePRIKeySQL()
	} else if st.Type == UNI {
		return fmt.Sprintf("UNIQUE KEY %s(%s)", st.Name, st.Name)
	}
	log.Fatalf("CreateKeySQL error!!! name:%s, type:%s, table:%s", st.Name, st.Type, table.Name)
	return ""
}

func (st *FieldKey) AddKeySQL(table *TableSt) string {
	switch st.Type {
	case MUL:
		return fmt.Sprintf("ALTER TABLE %s ADD INDEX %s(%s)", table.Name, st.Name, st.Name)
	case PRI:
		return table.AddPRIKeySQL()
	case UNI:
		return fmt.Sprintf("ALTER TABLE %s ADD UNIQUE (%s)", table.Name, st.Name)
	}
	log.Fatalf("%s %s add key sql error", table.Name, st.Name)
	return ""
}

func (st *FieldKey) IsEqual(info *MysqlColumnSt) bool {
	return st.Type == info.Key
}

var KeyFuncMap = map[string]func(name string) *FieldKey{}

func init() {
	KeyFuncMap["pri"] = func(name string) *FieldKey {
		return &FieldKey{Type: PRI, Name: name}
	}
	KeyFuncMap["mul"] = func(name string) *FieldKey {
		return &FieldKey{Type: MUL, Name: name}
	}
	KeyFuncMap["uni"] = func(name string) *FieldKey {
		return &FieldKey{Type: UNI, Name: name}
	}
}
