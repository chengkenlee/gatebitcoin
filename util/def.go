// @program:     gatebitcoin
// @file:        run.go
// @author:      chengkenlee
// @create:      2023-01-05 21:50
// @description:
package util

import (
	_ "gatebitcoin/tzinit"
	"github.com/blinkbean/dingtalk"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const ENCKEY = "**************"

var (
	TicksTime [3]string
	Config    *viper.Viper
	Logger    *zap.Logger
	Dingcli   *dingtalk.DingTalk
)
