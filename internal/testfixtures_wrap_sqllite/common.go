package testfixtures_wrap_sqllite

import (
	"context"
	"fmt"
	"github.com/zput/toolkit/internal/testfixtures"
	"github.com/zput/toolkit/internal/utils"
	"path/filepath"
	"strings"
)

type ThirdInitF = func(ctx context.Context) context.Context

type Flags struct {
	UseGoroutine     int    // 默认不使用协程, 直接串行执行
	MockDataSubPath  string // Mock数据库数据的yaml文件目录
	SelfDefineDBName string // DB名,主要用于sqlLite生成本地数据库时的名字(主要是防止并发导致sqlLite死锁)。当MockDataSubPath不为空，但是SelfDefineDBName为空，此时SelfDefineDBName被赋MockDataSubPath值
	IsOpenDbLog      bool   // 默认不打开db日志
}

type MockInterface interface {
	RetDbPath(flags Flags) string
	RetTables() (tables []interface{})
	// MockMutexLock()
}

// -----------------------------------------------------------

func NewMock(m MockInterface, flags Flags) *Mock {
	return &Mock{
		mockI: m,
		flags: flags,
	}
}

type Mock struct {
	mockI MockInterface
	flags Flags
}

func (m *Mock) baseByCtx(ctx context.Context, ff ...ThirdInitF) context.Context {
	var flags = m.flags

	if len(flags.MockDataSubPath) != 0 && len(flags.SelfDefineDBName) == 0 {
		flags.SelfDefineDBName = flags.MockDataSubPath
	}

	for _, f := range ff {
		ctx = f(ctx)
	}

	return ctx
}

// InjectDB 同步表结构
func (m *Mock) InjectDB(ctx context.Context, tablePrefix, driveName, dataSourceName string) ThirdInitF {
	return func(ctx context.Context) context.Context {
		var flags = m.flags
		if len(flags.MockDataSubPath) != 0 && len(flags.SelfDefineDBName) == 0 {
			flags.SelfDefineDBName = flags.MockDataSubPath
		}

		orm, errG := testfixtures.GenGorm(tablePrefix, driveName, dataSourceName, flags.IsOpenDbLog)
		if errG != nil {
			panic(errG)
		}

		testFixturesDb, err := testfixtures.SetUpFixture(
			filepath.Join(filepath.Dir(m.mockI.RetDbPath(flags)), flags.MockDataSubPath),
			orm,
			m.mockI.RetTables()...,
		)
		if err != nil {
			panic(err.Error() + "\n" + fmt.Sprintln(utils.ToString(flags)) + "\n" + m.mockI.RetDbPath(flags))
		}
		return SetDbToCtxWrap(ctx, &WrapDb{testFixturesDb.Gorm().WithContext(ctx)})
	}

}

func howManySlash(target string) string {
	idx := strings.LastIndex(target, "internal")
	if idx == -1 {
		return ""
	}

	nums := strings.Count(filepath.ToSlash(target[idx:]), "/")
	var ret string
	for i := nums; i >= 0; i-- {
		ret += "../"
	}
	return ret
}
