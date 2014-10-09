package plugins

import (
	"fmt"
	"time"
)

func Run(filename string) {

	gConfig, err := NewGlobalConfigFromFile(filename)
	if err != nil {
		fmt.Println(err)
		return
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
			fmt.Printf("Plugin '%v' start failed: %v\n", name, err)
			time.Sleep(time.Second)
			continue
		}
		break
	}
}
