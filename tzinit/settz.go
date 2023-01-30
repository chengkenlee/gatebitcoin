package tzinit

import (
	"os"
)

func init() {
	_ = os.Setenv("TZ", "Asia/Shanghai")
}
