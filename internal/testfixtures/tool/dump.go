package tool

import (
	"database/sql"
	"fmt"
	"github.com/go-testfixtures/testfixtures/v3"
	"xorm.io/xorm"
)

//const dbFilePath = "../../test.db"

func DumpData(dbFilePath, targetDir string) (err error) {
	engine, err := NewEngine(dbFilePath)
	if err != nil {
		return
	}

	dumper, err := testfixtures.NewDumper(
		testfixtures.DumpDatabase(engine.DB().DB),
		testfixtures.DumpDialect("sqlite"),
		testfixtures.DumpDirectory(targetDir),
	)
	if err != nil {
		return
	}
	if err = dumper.Dump(); err != nil {
		return
	}
	return
}

func NewEngine(sqliteFilePath string) (engine *xorm.Engine, err error) {
	engine, err = xorm.NewEngine("sqlite", sqliteFilePath+"?cache=shared&mode=memory")
	if err != nil {
		return nil, err
	}
	engine.ShowSQL(true)
	return
}

// GenFixtureByExistDB
// tables is optional, will dump all table if not given
func GenFixtureByExistDB(db *sql.DB, dialect, targetPath string, tables ...string) error {
	dumper, err := testfixtures.NewDumper(
		testfixtures.DumpDatabase(db),
		testfixtures.DumpDialect(dialect), // or your database of choice
		testfixtures.DumpDirectory(targetPath),
		testfixtures.DumpTables( // optional, will dump all table if not given
			tables...,
		),
	)
	if err != nil {
		return err
	}
	if err = dumper.Dump(); err != nil {
		return err
	}
	fmt.Println("success generate fixtures from exist database")
	return nil
}
