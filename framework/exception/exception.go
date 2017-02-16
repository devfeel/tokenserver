package exception

import (
	"fmt"
	"os"
	"runtime"

	"TechService/framework/log"
)

//统一异常处理
func CatchError(title string, logtarget string, loglevel string, err interface{}) (errmsg string) {
	errmsg = fmt.Sprintln(err)
	os.Stdout.Write([]byte(title + " error! => " + errmsg + " => "))
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, true)
	logger.Log(title+" error! => "+errmsg+" => "+string(buf[:n]), logtarget, loglevel)
	return errmsg
}
