package testfixtures_wrap_sqllite

import (
	"context"
	"gorm.io/gorm"
)

type WrapDb struct {
	*gorm.DB
}

func SetDbToCtx(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, dbKey{}, &WrapDb{})
	return ctx
}

func SetDbToCtxWrap(ctx context.Context, dbPtr *WrapDb) context.Context {
	ctx = context.WithValue(ctx, dbKey{}, dbPtr)
	return ctx
}

func GetDbFromCtx(ctx context.Context) *WrapDb {
	wdb := ctx.Value(dbKey{}).(*WrapDb)
	return wdb
}

type dbKey struct{}
