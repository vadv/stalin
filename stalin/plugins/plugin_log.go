package plugins

import (
	"fmt"
)

func LogInfo(format string, a ...interface{}) {
	RuntimeStalinConfig.Logger.Infof("%v", fmt.Sprintf(format, a...))
}

func LogWarn(format string, a ...interface{}) {
	RuntimeStalinConfig.Logger.Warnf("%v", fmt.Sprintf(format, a...))
}

func LogDebug(format string, a ...interface{}) {
	RuntimeStalinConfig.Logger.Debugf("%v", fmt.Sprintf(format, a...))
}

func LogErr(format string, a ...interface{}) {
	RuntimeStalinConfig.Logger.Errorf("%v", fmt.Sprintf(format, a...))
}
