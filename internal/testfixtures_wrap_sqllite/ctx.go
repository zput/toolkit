package testfixtures_wrap_sqllite

import (
	"context"
	"github.com/agiledragon/gomonkey/v2"
	"gorm.io/gorm"
)

type Append2Ctx = func(ctx context.Context) context.Context

func Append2CtxWrap(ctx context.Context, ff ...Append2Ctx) context.Context {
	for _, f := range ff {
		ctx = f(ctx)
	}
	return ctx
}

// -----------------------------------------------------------------------------------

type dbKey struct{}

type WrapDb struct {
	*gorm.DB
}

func (w *WrapDb) SetDbToCtxWrap() Append2Ctx {
	return func(ctx context.Context) context.Context {
		ctx = context.WithValue(ctx, dbKey{}, w)
		return ctx
	}
}
func SetDbToCtxWrap(ctx context.Context, w *WrapDb) context.Context {
	ctx = context.WithValue(ctx, dbKey{}, w)
	return ctx
}

func GetDbFromCtx(ctx context.Context) *WrapDb {
	wdb := ctx.Value(dbKey{}).(*WrapDb)
	return wdb
}

//func InitDbToCtx(ctx context.Context) context.Context {
//	ctx = context.WithValue(ctx, dbKey{}, &WrapDb{})
//	return ctx
//}

// ----------------------------------------------------

type goMonkeyKey struct{}

type WrapGoMonkey struct {
	*gomonkey.Patches
}

func (w *WrapGoMonkey) SetGoMonkeyKeyToCtxWrap() Append2Ctx {
	return func(ctx context.Context) context.Context {
		ctx = context.WithValue(ctx, goMonkeyKey{}, w)
		return ctx
	}
}

func SetGoMonkeyKeyToCtx(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, goMonkeyKey{}, &WrapGoMonkey{})
	return ctx
}

func SetGoMonkeyKeyToCtxWrap(ctx context.Context, w *WrapGoMonkey) context.Context {
	ctx = context.WithValue(ctx, goMonkeyKey{}, w)
	return ctx
}

func GetGoMonkeyKeyFromCtx(ctx context.Context) *WrapGoMonkey {
	wdb := ctx.Value(goMonkeyKey{}).(*WrapGoMonkey)
	return wdb
}

//func InitGoMonkeyKeyToCtx(ctx context.Context) context.Context {
//	ctx = context.WithValue(ctx, goMonkeyKey{}, &WrapGoMonkey{})
//	return ctx
//}

// ----------------------------------------------------
