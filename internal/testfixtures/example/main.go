package main

import (
	"fmt"
	"github.com/zput/toolkit/internal/testfixtures"
	"io/ioutil"
	"os"
)

type Example struct {
	Id   int64  `xorm:"pk autoincr bigint" json:"id"`
	Name string `xorm:"varchar(10) not null" json:"name"`
}

func main() {
	// testXorm()
	fmt.Println("=== ===")
	testGorm()
}

func testXorm() {
	var tablePrefix = "table_"

	xorm, errXorm := genXorm(tablePrefix, "", "")
	if errXorm != nil {
		panic(errXorm)
	}
	db, err := SetUpFixture(GetMockPath(), xorm, new(Example))
	if err != nil {
		panic(err)
	}

	verifyByXorm(db)
}

func testGorm() {
	var tablePrefix = "table_"

	xorm, errXorm := genGorm(tablePrefix, "sqlite", "gorm.db")
	if errXorm != nil {
		panic(errXorm)
	}
	db, err := SetUpFixture(GetMockPath(), xorm, new(Example))
	if err != nil {
		panic(err)
	}

	verifyByGorm(db)
}

func GetMockPath() string {
	var (
		path = "/internal/testfixtures/example/get_example_by_id"
	)

	dir, _ := os.Getwd()
	fmt.Println("当前路径：", dir)
	// TODO to make a blog
	//ex, err := os.Executable()
	//if err != nil {
	//	panic(err)
	//}
	//exPath := filepath.Dir(ex)
	//fmt.Println(exPath)

	path = dir + path
	fmt.Println("当前路径：", path)

	_, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	return path
}

func verifyByGorm(db testfixtures.DB) {
	var name string
	db.Gorm().Raw(`select name from table_example`).Scan(&name)

	fmt.Println(name)

	if name != "ztTest" {
		panic(fmt.Sprintf("expect ztTest, but get %s", name))
	} else {
		fmt.Println("(- v -), pass")
	}
}

func verifyByXorm(db testfixtures.DB) {
	var name string
	if _, err := db.Xorm().SQL(`select name from table_example`).Get(&name); err != nil {
		panic(err)
	}
	fmt.Println(name)
	if name != "ztTest" {
		panic(fmt.Sprintf("expect ztTest, but get %s", name))
	} else {
		fmt.Println("(- v -), pass")
	}
}

func genXorm(tablePrefix, driveName, dataSourceName string) (db testfixtures.IOrm, err error) {
	xorm, err := testfixtures.NewXOrm(
		testfixtures.DialectByXorm(driveName),
		testfixtures.DataSourceNameByXorm(dataSourceName),
		testfixtures.TablePrefixByXorm(tablePrefix),
	)
	if err != nil {
		return
	}
	return xorm, nil
}

func genGorm(tablePrefix, driveName, dataSourceName string) (db testfixtures.IOrm, err error) {
	gorm, err := testfixtures.NewGOrm(
		testfixtures.GOrmOptionDialect(driveName),
		testfixtures.GOrmOptionDataSourceName(dataSourceName),
		testfixtures.GOrmOptionTablePrefix(tablePrefix),
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
