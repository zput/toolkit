package testfixtures

import (
	"fmt"
	"github.com/go-testfixtures/testfixtures/v3"
	"io/ioutil"
)

type IOrm interface {
	Name() OrmType
	MigrationTableSchema(tables ...interface{}) error
	RetDb() DB
	Dialect() string // // Possible options are "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver".
}

type IFixture interface {
	MigrationTableSchema(tables ...interface{}) error // 表初始化到数据库
	LoadMockData() error
	RetDb() DB
}

/*
s := ztTest.New() // ---configuration---
s.Construct()
s.Load()
s.Engine()
*/

func NewFixture(o ...Option) (f *Fixture, err error) {
	f = new(Fixture)
	for i := range o {
		if err = o[i](f); err != nil {
			return
		}
	}
	var (
		db      = f.orm.RetDb()
		dialect = f.orm.Dialect()
		path    = f.mockDataPath
	)

	var para = []func(*testfixtures.Loader) error{
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.Database(db.StandDB()),
		testfixtures.Dialect(dialect),
	}
	if len(path) > 0 {
		para = append(para, testfixtures.Directory(path))
	}

	f.f, err = testfixtures.New(para...)
	if err != nil {
		return nil, err
	}
	return

}

// Fixture 脚手架 夹具
type Fixture struct {
	orm          IOrm
	mockDataPath string // path in order to load YAML files from a given directory
	f            *testfixtures.Loader
}

func (f *Fixture) MigrationTableSchema(tables ...interface{}) error {
	return f.orm.MigrationTableSchema(tables...)
}

func (f *Fixture) LoadMockData() error {
	if err := f.f.Load(); err != nil {
		return fmt.Errorf("[testfixtures] cannot load fixtures, err: %+v", err)
	}
	return nil
}

func (f *Fixture) RetDb() DB {
	return f.orm.RetDb()
}

type Option func(o *Fixture) error

func FixtureOptionOrm(orm IOrm) Option {
	return func(o *Fixture) (err error) {
		o.orm = orm
		return
	}
}

// FixtureOptionMockDataPath set path in order to load YAML files from a given directory.
func FixtureOptionMockDataPath(path string) Option {
	return func(o *Fixture) (err error) {
		_, err = ioutil.ReadDir(path)
		if err != nil {
			err = fmt.Errorf(`[testfixtures]: could not stat directory "%s": %w`, path, err)
			return
		}
		o.mockDataPath = path
		return
	}
}
