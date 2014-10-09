package plugins

import (
	"errors"
)

var PluginNotRegister = errors.New("Plugin is not registered (discribed)")
var PluginInitFatalError = errors.New("Fatal plugin error")
var PluginInitNoNeedToStartError = errors.New("Plugin no need to start")
var PluginNotImplemented = errors.New("MethodNotImpelemented")
