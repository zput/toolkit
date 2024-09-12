package testfixtures_wrap_sqllite

import "context"

func (m *Mock) Base(ctx context.Context, ff ...ThirdInitF) context.Context {
	return m.baseByCtx(ctx, ff...)
}

func (m *Mock) BaseBySqlLite(ctx context.Context, ff ...ThirdInitF) context.Context {
	ff = append(ff, m.injectSqlLiteDB)
	return m.baseByCtx(ctx, ff...)
}

func (m *Mock) injectSqlLiteDB(ctx context.Context) context.Context {
	return m.InjectDB(ctx, "", "sqlite", m.mockI.RetDbPath(m.flags))(ctx)
}

// dsn := "username:password@tcp(127.0.0.1:3306)/your_database_name?charset=utf8mb4&parseTime=True&loc=Local"
