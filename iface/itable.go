package iface

type ITable interface {
	CreatePriKeySQL() string
	AddPriKeySQL() string
	GetName() string
}
