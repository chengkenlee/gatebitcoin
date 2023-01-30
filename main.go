// @program:     gatebitcoin
// @file:        main.go
// @author:      chengkenlee
// @create:      2023-01-05 21:36
// @description:
package main

import (
	"gatebitcoin/run"
	"gatebitcoin/util"
)

func main() {
	util.Parm()
	util.Loggers()
	run.Run()
}
