package testfixtures

import (
	"database/sql"
	"xorm.io/xorm"
	xlog "xorm.io/xorm/log"
	"xorm.io/xorm/names"
)

func NewXOrm(o ...XOrmOption) (x *XOrm, err error) {
	x = new(XOrm)
	for i := range o {
		if err = o[i](x); err != nil {
			return
		}
	}
	if x.dialect == "" || x.dataSourceName == "" {
		x.dialect = "sqlite3"
		x.dataSourceName = "file::memory:?cache=shared"
	}

	var engine *xorm.Engine
	engine, err = xorm.NewEngine(x.dialect, x.dataSourceName)
	if err != nil {
		return
	}
	if len(x.tablePrefix) > 0 {
		engine.SetTableMapper(names.NewPrefixMapper(names.SnakeMapper{}, x.tablePrefix))
	}
	engine.ShowSQL(true)
	engine.SetLogLevel(xlog.LOG_DEBUG)

	return
}

type XOrm struct {
	dialect, dataSourceName string
	tablePrefix             string
	engine                  *xorm.Engine
}

func (x *XOrm) MigrationTableSchema(tables ...interface{}) (err error) {
	for _, v := range tables {
		err = x.engine.Sync2(v)
		if err != nil {
			return
		}
	}
	return
}

func (x *XOrm) RetDb() *sql.DB {
	return x.engine.DB().DB
}

func (x *XOrm) Dialect() string {
	return x.dialect
}

type XOrmOption func(x *XOrm) error

func Dialect(dialect string) XOrmOption {
	return func(x *XOrm) (err error) {
		x.dialect = dialect
		return
	}
}
func DataSourceName(dataSourceName string) XOrmOption {
	return func(x *XOrm) (err error) {
		x.dataSourceName = dataSourceName
		return
	}
}

func TablePrefix(tablePrefix string) XOrmOption {
	return func(x *XOrm) (err error) {
		x.tablePrefix = tablePrefix
		return
	}
}
