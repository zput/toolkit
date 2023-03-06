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
	var (
		path = "/internal/testfixtures/example/get_example_by_id"
		//path        = "/Users/edz/CODE/Self/zxcTool/ztTest/example/get_example_by_id"
		tablePrefix = "table_"
	)

	dir, _ := os.Getwd()
	fmt.Println("当前路径：", dir)

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

	if engine, err := example1(path, tablePrefix, "", ""); err != nil {
		panic(err)
	} else {
		var name string
		if _, err := engine.Xorm().SQL(`select name from table_example`).Get(&name); err != nil {
			panic(err)
		}
		fmt.Println(name)
		if name != "ztTest" {
			panic(fmt.Sprintf("expect ztTest, but get %s", name))
		} else {
			fmt.Println("(- v -), pass")
		}
	}

}

func example1(path string, tablePrefix, driveName, dataSourceName string) (db testfixtures.DB, err error) {

	xorm, err := testfixtures.NewXOrm(
		testfixtures.Dialect(driveName),
		testfixtures.DataSourceName(dataSourceName),
		testfixtures.TablePrefix(tablePrefix),
	)
	if err != nil {
		panic(err)
	}

	var f testfixtures.IFixture
	f, err = testfixtures.NewFixture(
		testfixtures.Orm(xorm),
		testfixtures.MockDataPath(path),
	)
	if err != nil {
		return
	}

	if err = f.MigrationTableSchema(
		new(Example),
	); err != nil {
		return
	}

	if err = f.LoadMockData(); err != nil {
		return
	}

	db = f.RetDb()
	return
}
