package testfixtures

import (
	. "github.com/zput/toolkit/internal/testfixtures/common"
	. "github.com/zput/toolkit/internal/testfixtures/gorm_sqllite"
)

func GenGorm(tablePrefix, driveName, dataSourceName string, isOpenDebug bool) (db IOrm, err error) {
	gorm, err := NewGOrm(
		GOrmOptionDialect(driveName),
		GOrmOptionDataSourceName(dataSourceName),
		GOrmOptionTablePrefix(tablePrefix),
		GOrmOptionOpenDebug(isOpenDebug),
	)
	if err != nil {
		return
	}
	return gorm, nil
}

func SetUpFixture(mockDataPath string, orm IOrm, tables ...interface{}) (db DB, err error) {
	var f IFixture
	f, err = NewFixture(
		FixtureOptionOrm(orm),
		FixtureOptionMockDataPath(mockDataPath),
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
