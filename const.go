package autodb

const (
	TINYBLOB   = "tinyblob"
	BLOB       = "blob"
	MEDIUMBLOB = "mediumblob"

	Int      = "int(10)"
	TinyInt  = "tinyint(3)"
	SmallInt = "smallint(5)"
	BigInt   = "bigint(20)"

	MySqlEngine = "INNODB"

	CreateSQLHead = "create table %s \n(\n"
	CreateSQLTail = "\n)\nENGINE=" + MySqlEngine + " DEFAULT CHARSET=utf8 COMMENT '%s'"

	ProcedureDropTemplate = "drop procedure if exists %s;"
)
