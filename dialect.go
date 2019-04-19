package gowhere

import "strings"

// Dialect represents the interface for a dialect
type Dialect interface {
	GetName() string
	QuoteIdentifier(string) string
}

type mysqlDialect struct{}
type postgresqlDialect struct{}

const (
	// DialectMySQLName defines the MySQL dialect name
	DialectMySQLName = "mysql"
	// DialectPostgreSQLName defines the PostgreSQL dialect name
	DialectPostgreSQLName = "postgres"
)

var (
	// DialectMySQL predefines the MySQL dialect
	DialectMySQL = &mysqlDialect{}
	// DialectPostgreSQL predefines the PostgreSQL dialect
	DialectPostgreSQL = &postgresqlDialect{}
)

func (md *mysqlDialect) GetName() string {
	return DialectMySQLName
}

func (md *mysqlDialect) QuoteIdentifier(name string) string {
	end := strings.IndexRune(name, 0)
	if end > -1 {
		name = name[:end]
	}
	name = "`" + strings.Replace(name, "`", "``", -1) + "`"
	return strings.Replace(name, ".", "`.`", -1)
}

func (pd *postgresqlDialect) GetName() string {
	return DialectPostgreSQLName
}

func (pd *postgresqlDialect) QuoteIdentifier(name string) string {
	end := strings.IndexRune(name, 0)
	if end > -1 {
		name = name[:end]
	}
	name = `"` + strings.Replace(name, `"`, `""`, -1) + `"`
	return strings.Replace(name, `.`, `"."`, -1)
}
