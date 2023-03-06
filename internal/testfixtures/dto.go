package testfixtures

import (
	"database/sql"
	"gorm.io/gorm"
	"xorm.io/xorm"
)

type OrmType int

const (
	GORM OrmType = iota
	XORM
)

type DB struct {
	gOrmDB *gorm.DB
	xOrmDB *xorm.Engine
	__type OrmType
}

func (d DB) StandDB() (db *sql.DB) {
	switch d.__type {
	case GORM:
		db, _ = d.gOrmDB.DB()
	case XORM:
		db = d.xOrmDB.DB().DB
	}
	return
}

func (d DB) Gorm() *gorm.DB {
	return d.gOrmDB
}

func (d DB) Xorm() *xorm.Engine {
	return d.xOrmDB
}
