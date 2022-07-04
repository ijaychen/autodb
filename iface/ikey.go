package iface

type IKey interface {
	GetType() string
	GetName() string
	CreateKeySQL(table ITable) string
	AddKeySQL(table ITable) string
	IsEqual(key string) bool
}
