package testfixtures

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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
	gormConfig := new(gorm.Config)
	gormConfig.NamingStrategy = schema.NamingStrategy{
		TablePrefix:   x.tablePrefix,
		SingularTable: true,
	}

	// gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	x.db, err = gorm.Open(sqlite.Open(x.dataSourceName), gormConfig)
	if err != nil {
		return
	}
	if x.isOpen {
		x.db = x.db.Debug()
	}
	return
}

type GOrm struct {
	dialect, dataSourceName string
	tablePrefix             string
	db                      *gorm.DB
	isOpen                  bool
}

func (x *GOrm) MigrationTableSchema(tables ...interface{}) (err error) {
	for _, v := range tables {
		err = x.db.Migrator().DropTable(v)
		err = x.db.AutoMigrate(v)
		if err != nil {
			return
		}
	}
	return
}

func (x *GOrm) RetDb() DB {
	return DB{
		gOrmDB: x.db,
		__type: GORM,
	}
}
func (x *GOrm) Name() OrmType {
	return GORM
}

func (x *GOrm) Dialect() string {
	return x.dialect
}

type GOrmOption func(x *GOrm) error

func GOrmOptionDialect(dialect string) GOrmOption {
	return func(x *GOrm) (err error) {
		x.dialect = dialect
		return
	}
}
func GOrmOptionDataSourceName(dataSourceName string) GOrmOption {
	return func(x *GOrm) (err error) {
		x.dataSourceName = dataSourceName
		return
	}
}

func GOrmOptionTablePrefix(tablePrefix string) GOrmOption {
	return func(x *GOrm) (err error) {
		x.tablePrefix = tablePrefix
		return
	}
}

func GOrmOptionOpenDebug(isOpen bool) GOrmOption {
	return func(x *GOrm) (err error) {
		x.isOpen = isOpen
		return
	}
}
