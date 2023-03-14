package tool

import (
	"database/sql"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestGenFixtureByExistDB(t *testing.T) {
	sqlDb, err := NewSqlDB()
	if err != nil {
		panic(err)
	}
	type args struct {
		db         *sql.DB
		dialect    string
		targetPath string
		tables     []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "success",
			args: args{
				db:         sqlDb,
				dialect:    "sqlite",
				targetPath: "./testdata",
				tables:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenFixtureByExistDB(tt.args.db, tt.args.dialect, tt.args.targetPath, tt.args.tables...); (err != nil) != tt.wantErr {
				t.Errorf("GenFixtureByExistDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func NewSqlDB() (db *sql.DB, err error) {
	var g *gorm.DB
	g, err = gorm.Open(sqlite.Open("gorm.db"))
	if err != nil {
		return
	}
	db, err = g.DB()
	return
}
