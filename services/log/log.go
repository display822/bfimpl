package log

import (
	"github.com/astaxie/beego/logs"
	"strings"
)

var GLogger *APILogger

func init() {
	// 在bootstrap未设置GLogger时候，防止在service等地方使用GLogger导致空指针panic
	// 是初始化顺序导致的
	GLogger = new(APILogger)
	GLogger.SetLevel(logs.LevelDebug)

}

type APILogger struct {
	logs.BeeLogger
}

// Print adapter for gorm logger
func (a *APILogger) Print(v ...interface{}) {
	fmstr := strings.Repeat("%v ", len(v))
	a.Info(fmstr, v)
}

func NewLogger() *APILogger {
	return new(APILogger)
}

func SetLogger(l *APILogger) {
	GLogger = l
}
