package testfixtures

import (
	"github.com/zput/toolkit/internal/testfixtures_wrap_sqllite"
)

// type

type (
	Append2Ctx        = testfixtures_wrap_sqllite.Append2Ctx
	Flags             = testfixtures_wrap_sqllite.Config
	CallerDbInterface = testfixtures_wrap_sqllite.CallerDbInterface
	ConfigInitF       = testfixtures_wrap_sqllite.ConfigInitF
	Mock              = testfixtures_wrap_sqllite.Mock

	WrapGoMonkey = testfixtures_wrap_sqllite.WrapGoMonkey
)

// variable

var (
	NewMock                     = testfixtures_wrap_sqllite.NewMock
	ConfigInitCallerDbInterface = testfixtures_wrap_sqllite.ConfigInitCallerDbInterface
	ConfigInitIsOpenDbLog       = testfixtures_wrap_sqllite.ConfigInitIsOpenDbLog

	GetDbFromCtx   = testfixtures_wrap_sqllite.GetDbFromCtx
	SetDbToCtxWrap = testfixtures_wrap_sqllite.SetDbToCtxWrap

	GetGoMonkeyKeyFromCtx   = testfixtures_wrap_sqllite.GetGoMonkeyKeyFromCtx
	SetGoMonkeyKeyToCtx     = testfixtures_wrap_sqllite.SetGoMonkeyKeyToCtx
	SetGoMonkeyKeyToCtxWrap = testfixtures_wrap_sqllite.SetGoMonkeyKeyToCtxWrap
)
