package yaorm

import (
	"github.com/geoffreybauduin/yaorm/_vendor/github.com/lann/squirrel"
	"github.com/geoffreybauduin/yaorm/_vendor/github.com/loopfz/gadgeto/zesty"
	"github.com/geoffreybauduin/yaorm/tools"
	"github.com/go-gorp/gorp"
)

// DBProvider provides an abstracted way of accessing the database
type DBProvider interface {
	zesty.DBProvider
	EscapeValue(value string) string
	CanSelectForUpdate() bool
	getStatementGenerator() squirrel.StatementBuilderType
}

type dbprovider struct {
	zesty.DBProvider
	name string
}

// NewDBProvider creates a new db provider
func NewDBProvider(name string) (DBProvider, error) {
	dblock.RLock()
	defer dblock.RUnlock()
	dbp, err := zesty.NewDBProvider(name)
	if err != nil {
		return nil, err
	}
	return &dbprovider{DBProvider: dbp, name: name}, nil
}

// DB returns a SQL Executor interface
func (dbp *dbprovider) DB() gorp.SqlExecutor {
	db := registry[dbp.name]
	return &SqlExecutor{DB: db}
}

// EscapeValue escapes the value sent according to the dialect
func (dbp *dbprovider) EscapeValue(value string) string {
	return dbp.getDialect().QuoteField(value)
}

// CanSelectForUpdate returns true if the current dialect can perform select for update statements
func (dbp *dbprovider) CanSelectForUpdate() bool {
	db := registry[dbp.name]
	switch db.System() {
	case DatabasePostgreSQL:
		return true
	}
	return false
}

func (dbp *dbprovider) getDb() DB {
	return registry[dbp.name]
}

func (dbp *dbprovider) getDialect() gorp.Dialect {
	v := tools.GetNonPtrValue(dbp.DB())
	dbField := tools.GetNonPtrValue(v.FieldByName("DB").Interface()).FieldByName("DB")
	realValue := tools.GetNonPtrValue(dbField.Interface())
	field := realValue.FieldByName("DbMap")
	s := field.Interface().(*gorp.DbMap)
	return s.Dialect
}

func (dbp *dbprovider) getStatementGenerator() squirrel.StatementBuilderType {
	switch dbp.getDb().System() {
	case DatabaseMySQL:
		return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	}
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
