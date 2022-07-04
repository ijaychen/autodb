package columnkey

import (
	"fmt"
	"github.com/ijaychen/autodb/iface"
	"strings"
)

var (
	_ iface.IKey = new(PriKey)
	_ iface.IKey = new(MulKey)
	_ iface.IKey = new(UniKey)
)

type (
	Base struct {
		Type string
		Name string
	}

	PriKey struct {
		Base
	}

	MulKey struct {
		Base
	}

	UniKey struct {
		Base
	}
)

func CreateKey(kt, name string) iface.IKey {
	switch strings.ToUpper(kt) {
	case PRI:
		return &PriKey{Base{Type: kt, Name: name}}
	case MUL:
		return &MulKey{Base{Type: kt, Name: name}}
	case UNI:
		return &UniKey{Base{Type: kt, Name: name}}
	}
	return nil
}

func (base *Base) GetType() string {
	return base.Type
}

func (base *Base) GetName() string {
	return base.Name
}

func (base *Base) IsEqual(key string) bool {
	return base.Type == key
}

func (key *PriKey) CreateKeySQL(table iface.ITable) string {
	return table.CreatePriKeySQL()
}

func (key *PriKey) AddKeySQL(table iface.ITable) string {
	return table.AddPriKeySQL()
}

func (key *MulKey) CreateKeySQL(iface.ITable) string {
	return fmt.Sprintf("KEY %s(%s)", key.Name, key.Name)
}

func (key *MulKey) AddKeySQL(table iface.ITable) string {
	return fmt.Sprintf("ALTER TABLE %s ADD INDEX %s(%s)", table.GetName(), key.Name, key.Name)
}

func (key *UniKey) CreateKeySQL(iface.ITable) string {
	return fmt.Sprintf("UNIQUE KEY %s(%s)", key.Name, key.Name)
}

func (key *UniKey) AddKeySQL(table iface.ITable) string {
	return fmt.Sprintf("ALTER TABLE %s ADD UNIQUE (%s)", table.GetName(), key.Name)
}
