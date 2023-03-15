package testfixtures

import (
	"github.com/zput/toolkit/internal/testfixtures"
	"github.com/zput/toolkit/internal/testfixtures/tool"
)

type (
	Fixture = testfixtures.Fixture
	GOrm    = testfixtures.GOrm
	XOrm    = testfixtures.XOrm
)

var (
	NewGOrm                  = testfixtures.NewGOrm
	GOrmOptionDialect        = testfixtures.GOrmOptionDialect
	GOrmOptionDataSourceName = testfixtures.GOrmOptionDataSourceName
	GOrmOptionTablePrefix    = testfixtures.GOrmOptionTablePrefix

	NewFixture                = testfixtures.NewFixture
	FixtureOptionMockDataPath = testfixtures.FixtureOptionMockDataPath
	FixtureOptionOrm          = testfixtures.FixtureOptionOrm
)

type (
	IFixture = testfixtures.IFixture
	IOrm     = testfixtures.IOrm
)

var (
	// dump yml 数据

	GenFixtureByExistDB = tool.GenFixtureByExistDB
	NewSqlDBByGorm      = tool.NewSqlDBByGorm
)

func Example() {
	var tablePrefix = "table_"

	gorm, errXorm := GenGorm(tablePrefix, "sqlite", "gorm.db", true)
	if errXorm != nil {
		panic(errXorm)
	}
	var mockDataPath = "/xx"
	db, err := SetUpFixture(mockDataPath, gorm)
	if err != nil {
		panic(err)
	}

	_ = db
}

func GenGorm(tablePrefix, driveName, dataSourceName string, isOpenDebug bool) (db testfixtures.IOrm, err error) {
	gorm, err := testfixtures.NewGOrm(
		testfixtures.GOrmOptionDialect(driveName),
		testfixtures.GOrmOptionDataSourceName(dataSourceName),
		testfixtures.GOrmOptionTablePrefix(tablePrefix),
		testfixtures.GOrmOptionOpenDebug(isOpenDebug),
	)
	if err != nil {
		return
	}
	return gorm, nil
}

func SetUpFixture(mockDataPath string, orm testfixtures.IOrm, tables ...interface{}) (db testfixtures.DB, err error) {
	var f testfixtures.IFixture
	f, err = testfixtures.NewFixture(
		testfixtures.FixtureOptionOrm(orm),
		testfixtures.FixtureOptionMockDataPath(mockDataPath),
	)
	if err != nil {
		return
	}

	if err = f.MigrationTableSchema(
		tables...,
	); err != nil {
		return
	}

	if err = f.LoadMockData(); err != nil {
		return
	}

	db = f.RetDb()
	return
}
