package plugins

import (
	"fmt"
	"reflect"
	"stalin/events"
	. "stalin/log"
	"time"
)

const retryStartSec = 1

type PluginPipeline struct {
	outputs map[string]PluginOutputInterface
	inputs  map[string]PluginInterface
}

var Pipeline = initPluginPipeline()

func initPluginPipeline() *PluginPipeline {
	return &PluginPipeline{
		outputs: make(map[string]PluginOutputInterface, 0),
		inputs:  make(map[string]PluginInterface, 0),
	}
}

func (p *PluginPipeline) AddPlugin(name string, i PluginInterface) error {
	if !isOutputPlugin(i) {
		p.inputs[name] = i
	} else {
		output, ok := i.(PluginOutputInterface)
		if !ok {
			return fmt.Errorf("Plugin %s: missing Inject method", name)
		}
		p.outputs[name] = output
	}
	return nil
}

func isOutputPlugin(i PluginInterface) bool {
	typ := reflect.TypeOf(i)
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		if method.Name == "Inject" {
			return true
		}
	}
	return false
}

func (p *PluginPipeline) loopStart(i PluginInterface, name string) {
	for {
		if err := i.Start(); err != nil {
			Log.Error("[%s] Start error (retry after %ds): %v", name, retryStartSec, err)
			time.Sleep(time.Duration(retryStartSec) * time.Second)
			continue
		}
		break
	}
}

func (p *PluginPipeline) inject(e *events.Events) {
	for _, plugin := range p.outputs {
		plugin.Inject(e)
	}
}

func (p *PluginPipeline) RunAll() {
	for name, plugin := range p.inputs {
		go func(i PluginInterface, n string) { p.loopStart(i, n) }(plugin, name)
	}
	for name, plugin := range p.outputs {
		go func(i PluginInterface, n string) { p.loopStart(i, n) }(plugin, name)
	}
}

func (p *PluginPipeline) Inject(e *events.Events) {
	go p.inject(e)
}
