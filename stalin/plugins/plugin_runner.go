package plugins

import (
	"fmt"
	"os"
	"time"
)

func Run(filename string) {

	gConfig, err := NewGlobalConfigFromFile(filename)
	if err != nil {
		fmt.Printf("Error while init config: %v\n", err)
		os.Exit(1)
	}

	gConfig.Router.Run()

	for name, plugin := range gConfig.Inputs {
		go loopStart(plugin, name)
	}

	for name, plugin := range gConfig.Outputs {
		go loopStart(plugin, name)
	}

}

func loopStart(plugin Plugin, name string) {
	for {
		if err := plugin.Run(); err != nil {
			LogErr("Plugin '%v' start failed: %v", name, err)
			time.Sleep(time.Second)
			continue
		}
		break
	}
}
