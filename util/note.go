// @program:     gatebitcoin
// @file:        run.go
// @author:      chengkenlee
// @create:      2023-01-05 21:50
// @description:
package util

import (
	"bufio"
	"fmt"
	_ "gatebitcoin/tzinit"
	"os"
)

func Writelog(filename, msg string) {

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Logger.Info(fmt.Sprintf("%s", err.Error()))
		return
	}
	defer file.Close()
	// NewWriter 默认缓冲区大小是 4096
	// 需要使用自定义缓冲区的writer 使用 NewWriterSize()方法
	buf := bufio.NewWriterSize(file, len(msg))

	buf.WriteString(msg)

	err = buf.Flush()
	if err != nil {
		Logger.Info(err.Error())
	}
}
