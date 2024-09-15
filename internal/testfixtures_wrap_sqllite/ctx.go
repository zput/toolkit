package testfixtures_wrap_sqllite

import (
	"context"
	"github.com/agiledragon/gomonkey/v2"
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

// ---

type WrapGoMonkey struct {
	*gomonkey.Patches
}

func SetGoMonkeyKeyToCtx(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, goMonkeyKey{}, &WrapGoMonkey{})
	return ctx
}

func SetGoMonkeyKeyToCtxWrap(ctx context.Context, dbPtr *WrapGoMonkey) context.Context {
	ctx = context.WithValue(ctx, goMonkeyKey{}, dbPtr)
	return ctx
}

func GetGoMonkeyKeyFromCtx(ctx context.Context) *WrapGoMonkey {
	wdb := ctx.Value(goMonkeyKey{}).(*WrapGoMonkey)
	return wdb
}

type goMonkeyKey struct{}
