package testfixtures

import (
	"github.com/zput/toolkit/internal/testfixtures"
	"github.com/zput/toolkit/internal/testfixtures/gorm_sqllite"
	"github.com/zput/toolkit/internal/testfixtures/tool"
	"github.com/zput/toolkit/internal/testfixtures/xorm_sqllite"
)

type (
	Fixture = testfixtures.Fixture
	GOrm    = gorm_sqllite.GOrm
	XOrm    = xorm_sqllite.XOrm
)

var (
	NewGOrm                  = gorm_sqllite.NewGOrm
	GOrmOptionDialect        = gorm_sqllite.GOrmOptionDialect
	GOrmOptionDataSourceName = gorm_sqllite.GOrmOptionDataSourceName
	GOrmOptionTablePrefix    = gorm_sqllite.GOrmOptionTablePrefix

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

var (
	GenGorm      = testfixtures.GenGorm
	SetUpFixture = testfixtures.SetUpFixture
)
