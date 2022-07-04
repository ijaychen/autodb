package autodb

const (
	MySqlEngine = "INNODB"

	CreateSQLHead = "create table %s \n(\n"
	CreateSQLTail = "\n)\nENGINE=" + MySqlEngine + " DEFAULT CHARSET=utf8 COMMENT '%s'"

	ProcedureDropTemplate = "drop procedure if exists %s;"
)
