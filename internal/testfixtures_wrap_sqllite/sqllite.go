package testfixtures_wrap_sqllite

import (
	"context"
	"github.com/zput/toolkit/internal/testfixtures"
	"github.com/zput/toolkit/internal/utils"
	"gorm.io/gorm"
	"path/filepath"
)

type Config struct {
	CallerDbInterface
	IsOpenDbLog bool // 默认不打开db日志

	// UseGoroutine     int    // 默认不使用协程, 直接串行执行
	MockDataSubPath  string // Mock数据库数据的yaml文件目录
	SelfDefineDBName string // DB名,主要用于sqlLite生成本地数据库时的名字(主要是防止并发导致sqlLite死锁)。当MockDataSubPath不为空，但是SelfDefineDBName为空，此时SelfDefineDBName被赋MockDataSubPath值
}

func (c Config) String() string {
	return utils.ToString(c)
}

type CallerDbInterface interface {
	DefineDbTableModels() (tableModels []interface{})
	DefineDbLocationByWillGen() (dbLocation string)
}

type ConfigInitF = func(*Config)

func ConfigInitIsOpenDbLog(isOpenDbLog bool) ConfigInitF {
	return func(in *Config) {
		in.IsOpenDbLog = isOpenDbLog
	}
}

func ConfigInitCallerDbInterface(i CallerDbInterface) ConfigInitF {
	return func(in *Config) {
		in.CallerDbInterface = i
	}
}

// -----------------------------------------------------------
var _ Append2Ctx = NewMock().AppendSqlLiteToCtx()

func NewMock(ff ...ConfigInitF) *Mock {
	var tmp = new(Mock)
	for _, f := range ff {
		f(&tmp.Config)
	}
	return tmp
}

type Mock struct {
	Config
}

func (m *Mock) FixtureByGorm() *gorm.DB {
	return m.fixtureByGorm("", "sqlite", m.DefineDbLocationByWillGen())
}

func (m *Mock) fixtureByGorm(tablePrefix, driveName, dataSourceName string) *gorm.DB {
	// 1. gen orm
	orm, errG := testfixtures.GenGorm(tablePrefix, driveName, dataSourceName, m.IsOpenDbLog)
	if errG != nil {
		panic(errG)
	}
	// 2. fixture orm
	testFixturesDb, err := testfixtures.SetUpFixture(
		filepath.Join(filepath.Dir(m.DefineDbLocationByWillGen()), m.MockDataSubPath),
		orm,
		m.DefineDbTableModels()...,
	)
	if err != nil {
		panic(err.Error() + "\n" + m.String() + "\n" + m.DefineDbLocationByWillGen())
	}
	return testFixturesDb.Gorm()
}

func (m *Mock) AppendSqlLiteToCtx() Append2Ctx {
	return func(ctx context.Context) context.Context {
		return SetDbToCtxWrap(ctx, &WrapDb{m.FixtureByGorm().WithContext(ctx)})
	}
}

// dsn := "username:password@tcp(127.0.0.1:3306)/your_database_name?charset=utf8mb4&parseTime=True&loc=Local"
