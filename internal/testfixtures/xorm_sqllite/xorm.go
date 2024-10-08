package xorm_sqllite

import (
	. "github.com/zput/toolkit/internal/testfixtures/common"
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

	x.engine, err = xorm.NewEngine(x.dialect, x.dataSourceName)
	if err != nil {
		return
	}
	if len(x.tablePrefix) > 0 {
		x.engine.SetTableMapper(names.NewPrefixMapper(names.SnakeMapper{}, x.tablePrefix))
	}
	x.engine.ShowSQL(true)
	x.engine.SetLogLevel(xlog.LOG_DEBUG)

	return
}

type XOrm struct {
	dialect, dataSourceName string
	tablePrefix             string
	engine                  *xorm.Engine
}

func (x *XOrm) Name() OrmType {
	return XORM
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

func (x *XOrm) RetDb() DB {
	return NewDB(nil, x.engine, XORM)
}

func (x *XOrm) Dialect() string {
	return x.dialect
}

type XOrmOption func(x *XOrm) error

func DialectByXorm(dialect string) XOrmOption {
	return func(x *XOrm) (err error) {
		x.dialect = dialect
		return
	}
}
func DataSourceNameByXorm(dataSourceName string) XOrmOption {
	return func(x *XOrm) (err error) {
		x.dataSourceName = dataSourceName
		return
	}
}

func TablePrefixByXorm(tablePrefix string) XOrmOption {
	return func(x *XOrm) (err error) {
		x.tablePrefix = tablePrefix
		return
	}
}
