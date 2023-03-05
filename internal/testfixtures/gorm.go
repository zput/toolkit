package testfixtures

import (
	"database/sql"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewGOrm(o ...GOrmOption) (x *GOrm, err error) {
	x = new(GOrm)
	for i := range o {
		if err = o[i](x); err != nil {
			return
		}
	}
	if x.dialect == "" || x.dataSourceName == "" {
		x.dialect = "sqlite3"
		x.dataSourceName = "file::memory:?cache=shared"
	}

	// gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	x.db, err = gorm.Open(sqlite.Open(x.dataSourceName), &gorm.Config{})
	if err != nil {
		return
	}
	x.db = x.db.Debug()
	return
}

type GOrm struct {
	dialect, dataSourceName string
	tablePrefix             string
	db                      *gorm.DB
}

func (x *GOrm) MigrationTableSchema(tables ...interface{}) (err error) {
	for _, v := range tables {
		err = x.db.Migrator().CreateTable(v)
		if err != nil {
			return
		}
	}
	return
}

func (x *GOrm) RetDb() (db *sql.DB, err error) {
	return x.db.Debug().DB()
}

func (x *GOrm) Dialect() string {
	return x.dialect
}

type GOrmOption func(x *GOrm) error

func DialectByGOrm(dialect string) GOrmOption {
	return func(x *GOrm) (err error) {
		x.dialect = dialect
		return
	}
}
func DataSourceNameByGOrm(dataSourceName string) GOrmOption {
	return func(x *GOrm) (err error) {
		x.dataSourceName = dataSourceName
		return
	}
}

func TablePrefixByGOrm(tablePrefix string) GOrmOption {
	return func(x *GOrm) (err error) {
		x.tablePrefix = tablePrefix
		return
	}
}