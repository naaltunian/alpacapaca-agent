package recovery

import (
	"fmt"

	"github.com/naaltunian/paca-agent/mailer"
)

func Recover() {
	if r := recover(); r != nil {
		errStr := fmt.Sprintf("%v", r)
		mailer.Notify("Agent panic: " + errStr)
	}
}
