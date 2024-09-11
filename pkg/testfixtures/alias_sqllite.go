package testfixtures

import (
	"github.com/zput/toolkit/internal/testfixtures_wrap_sqllite"
)

// type

type (
	ThirdInitF    = testfixtures_wrap_sqllite.ThirdInitF
	Flags         = testfixtures_wrap_sqllite.Flags
	MockInterface = testfixtures_wrap_sqllite.MockInterface
	Mock          = testfixtures_wrap_sqllite.Mock
)

// variable

var (
	NewMock        = testfixtures_wrap_sqllite.NewMock
	GetDbFromCtx   = testfixtures_wrap_sqllite.GetDbFromCtx
	SetDbToCtxWrap = testfixtures_wrap_sqllite.SetDbToCtxWrap
)
